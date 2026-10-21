import 'dart:async';
import 'dart:convert';
import 'dart:io';

import 'package:flutter/services.dart';

class CoreReady {
  const CoreReady({required this.port, required this.token});

  final int port;
  final String token;
}

class CoreExit {
  const CoreExit({required this.code, required this.stderr});

  final int code;
  final String stderr;
}

CoreReady? parseCoreReadyLine(String line) {
  const prefix = 'ANTO_CORE_READY ';
  if (!line.startsWith(prefix)) return null;
  final value = jsonDecode(line.substring(prefix.length));
  if (value is! Map || value['port'] is! int || value['token'] is! String) {
    throw const FormatException('ANTO_CORE_READY 字段无效');
  }
  return CoreReady(port: value['port'] as int, token: value['token'] as String);
}

class CoreClient {
  static const _windowChannel = MethodChannel('anto/windows_window');
  Process? _process;
  CoreReady? _ready;
  Future<void>? _startOperation;
  final StreamController<CoreExit> _exitController =
      StreamController<CoreExit>.broadcast();
  String _stderrTail = '';

  Stream<CoreExit> get onExit => _exitController.stream;

  String get stderrTail => _stderrTail;

  bool get isRunning => _process != null && _ready != null;

  Future<void> start({String executable = 'anto-core.exe'}) {
    if (isRunning) return Future<void>.value();
    final pending = _startOperation;
    if (pending != null) return pending;
    late Future<void> operation;
    operation = _start(executable).whenComplete(() {
      if (identical(_startOperation, operation)) _startOperation = null;
    });
    _startOperation = operation;
    return operation;
  }

  Future<void> _start(String executable) async {
    if (_process != null) return;
    final executablePath = executable == 'anto-core.exe'
        ? File(Platform.resolvedExecutable).parent.path +
            Platform.pathSeparator +
            executable
        : executable;
    final process = await Process.start(
      executablePath,
      const [],
      runInShell: false,
      // Core loads cfg.yml from the current working directory. A packaged
      // client must use the directory beside anto-core.exe, regardless of
      // the directory from which anto-ui.exe was launched.
      workingDirectory: File(executablePath).parent.path,
    );
    try {
      await _windowChannel.invokeMethod<void>('attachCoreProcess', process.pid);
    } catch (error) {
      process.kill(ProcessSignal.sigterm);
      await process.exitCode.timeout(const Duration(seconds: 2), onTimeout: () {
        process.kill(ProcessSignal.sigkill);
        return -1;
      });
      throw StateError('无法绑定 Go Core 进程生命周期: $error');
    }
    _process = process;
    _ready = null;
    _stderrTail = '';
    final completer = Completer<CoreReady>();
    final stdoutSubscription = process.stdout
        // A failed Windows shell launch can emit non-UTF8 bytes before the
        // process exits. Keep the exit error actionable instead of allowing
        // the decoder stream to surface an unrelated FormatException.
        .transform(const Utf8Decoder(allowMalformed: true))
        .transform(const LineSplitter())
        .listen((line) {
      if (completer.isCompleted) return;
      try {
        final ready = parseCoreReadyLine(line);
        if (ready != null) completer.complete(ready);
      } catch (error) {
        completer.completeError(StateError('ANTO_CORE_READY 无效: $error'));
      }
    });
    final stderrSubscription = process.stderr
        .transform(const Utf8Decoder(allowMalformed: true))
        .listen(_appendStderr);
    final exitFuture = process.exitCode;
    unawaited(_watchExit(process, exitFuture, completer));
    try {
      _ready = await completer.future.timeout(const Duration(seconds: 10));
    } on TimeoutException {
      throw StateError('Go Core 启动超时${_stderrSuffix()}');
    } finally {
      await stdoutSubscription.cancel();
      await stderrSubscription.cancel();
      if (_ready == null && identical(_process, process)) {
        process.kill(ProcessSignal.sigterm);
        await exitFuture.timeout(const Duration(seconds: 2),
            onTimeout: () => -1);
        _process = null;
      }
    }
  }

  Future<void> _watchExit(Process process, Future<int> exitFuture,
      Completer<CoreReady> ready) async {
    final code = await exitFuture;
    if (!ready.isCompleted) {
      ready.completeError(
          StateError('Go Core 提前退出，退出码 $code${_stderrSuffix()}'));
    }
    if (identical(_process, process)) {
      _process = null;
      _ready = null;
      if (!_exitController.isClosed) {
        _exitController.add(CoreExit(code: code, stderr: _stderrTail));
      }
    }
  }

  void _appendStderr(String value) {
    _stderrTail += value;
    if (_stderrTail.length > 4000) {
      _stderrTail = _stderrTail.substring(_stderrTail.length - 4000);
    }
  }

  String _stderrSuffix() => _stderrTail.trim().isEmpty
      ? ''
      : ': ${_stderrTail.trim().split('\n').last}';

  Future<dynamic> get(String path) => _request('GET', path);

  Future<dynamic> put(String path, Map<String, dynamic> body) =>
      _request('PUT', path, body: body);

  Future<dynamic> post(String path, {Map<String, dynamic>? body}) =>
      _request('POST', path, body: body);

  /// Opens the Core event stream and decodes one JSON object per SSE frame.
  /// The returned stream owns its HttpClient and closes it when cancelled.
  Stream<Map<String, dynamic>> events() async* {
    final ready = _ready;
    if (ready == null) throw StateError('Go Core 尚未启动');
    final client = HttpClient();
    try {
      final request = await client
          .getUrl(Uri.parse('http://127.0.0.1:${ready.port}/api/v1/events'));
      request.headers.set('X-Anto-Token', ready.token);
      request.headers.set('Accept', 'text/event-stream');
      final response =
          await request.close().timeout(const Duration(seconds: 10));
      if (response.statusCode < 200 || response.statusCode >= 300) {
        final text = await utf8.decoder.bind(response).join();
        throw HttpException(text, uri: request.uri);
      }

      String? eventName;
      final data = StringBuffer();
      await for (final line
          in utf8.decoder.bind(response).transform(const LineSplitter())) {
        if (line.startsWith('event:')) {
          eventName = line.substring('event:'.length).trim();
        } else if (line.startsWith('data:')) {
          if (data.length > 0) data.write('\n');
          data.write(line.substring('data:'.length).trimLeft());
        } else if (line.isEmpty && data.length > 0) {
          final decoded = jsonDecode(data.toString());
          if (decoded is Map) {
            final event = decoded.cast<String, dynamic>();
            if (event['type'] == null && eventName != null) {
              event['type'] = eventName;
            }
            yield event;
          }
          eventName = null;
          data.clear();
        }
      }
      if (data.length > 0) {
        final decoded = jsonDecode(data.toString());
        if (decoded is Map) {
          final event = decoded.cast<String, dynamic>();
          if (event['type'] == null && eventName != null) {
            event['type'] = eventName;
          }
          yield event;
        }
      }
    } finally {
      client.close(force: true);
    }
  }

  Future<dynamic> _request(String method, String path,
      {Map<String, dynamic>? body}) async {
    final ready = _ready;
    if (ready == null) throw StateError('Go Core 尚未启动');
    final client = HttpClient();
    try {
      final request = await client.openUrl(
        method,
        Uri.parse('http://127.0.0.1:${ready.port}$path'),
      );
      request.headers.set('X-Anto-Token', ready.token);
      request.headers.contentType = ContentType.json;
      if (body != null) request.write(jsonEncode(body));
      final response =
          await request.close().timeout(const Duration(seconds: 10));
      final text = await utf8.decoder.bind(response).join();
      if (response.statusCode < 200 || response.statusCode >= 300) {
        throw HttpException(text, uri: request.uri);
      }
      return text.isEmpty ? null : jsonDecode(text);
    } finally {
      client.close(force: true);
    }
  }

  Future<void> stop() async {
    final process = _process;
    if (process == null) {
      _ready = null;
      return;
    }
    try {
      if (_ready != null) {
        await post('/api/v1/shutdown').timeout(const Duration(seconds: 2));
      }
    } catch (_) {
      // The process may already be shutting down. The local kill below is
      // still required to avoid leaving a stale Core process behind.
    } finally {
      process.kill(ProcessSignal.sigterm);
      await process.exitCode.timeout(const Duration(seconds: 2), onTimeout: () {
        process.kill(ProcessSignal.sigkill);
        return -1;
      });
      if (identical(_process, process)) {
        _process = null;
        _ready = null;
      }
    }
  }

  Future<void> dispose() async {
    await stop();
    await _exitController.close();
  }
}
