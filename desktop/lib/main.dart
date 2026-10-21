import 'dart:async';
import 'dart:convert';
import 'dart:io';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import 'about_page.dart';
import 'aliyun_regions.dart';
import 'core_client.dart';
import 'input_validation.dart';
import 'tencent_regions.dart';
import 'windows_picker.dart';

const _ink = Color(0xff202638);
const _muted = Color(0xff7c8498);
const _brand = Color(0xff5b5ee8);
const _sidebar = Color(0xff20243b);

const _fallbackLanguages = <Map<String, String>>[
  {'id': 'en', 'name': '英语'},
  {'id': 'zh', 'name': '中文'},
  {'id': 'zh-cn', 'name': '中文（简体）'},
  {'id': 'ja', 'name': '日语'},
  {'id': 'ko', 'name': '韩语'},
  {'id': 'fr', 'name': '法语'},
  {'id': 'de', 'name': '德语'},
  {'id': 'ru', 'name': '俄语'},
];

void main() => runApp(const AntoApp());

class AntoApp extends StatefulWidget {
  const AntoApp({super.key, this.core, this.picker, this.autoConnect = true});

  final CoreClient? core;
  final WindowsPicker? picker;
  final bool autoConnect;

  @override
  State<AntoApp> createState() => _AntoAppState();
}

class _AntoAppState extends State<AntoApp> {
  late final CoreClient _core;
  late final WindowsPicker _picker;
  final _navigatorKey = GlobalKey<NavigatorState>();
  static const _dropChannel = MethodChannel('anto/windows_drop');
  static const _windowChannel = MethodChannel('anto/windows_window');
  final _filesController = TextEditingController();
  final _inputDirectoryController = TextEditingController();
  final _directoryController = TextEditingController();
  final _fromController = TextEditingController(text: 'en');
  final _toController = TextEditingController(text: 'zh');
  final _concurrencyController = TextEditingController(text: '5');
  Timer? _refreshTimer;
  Timer? _reconnectTimer;
  StreamSubscription<CoreExit>? _coreExitSubscription;
  StreamSubscription<Map<String, dynamic>>? _eventsSubscription;
  String _status = '正在启动 Go Core...';
  List<dynamic> _translators = const [];
  List<dynamic> _tasks = const [];
  Map<String, dynamic> _appConfig = <String, dynamic>{};
  Map<String, Map<String, dynamic>> _providerConfigs =
      <String, Map<String, dynamic>>{};
  Map<String, bool> _providerConfigured = <String, bool>{};
  final Map<String, Map<String, TextEditingController>> _providerControllers =
      <String, Map<String, TextEditingController>>{};
  String? _selectedTranslator;
  String _mode = '全量翻译';
  int _trackExport = 2;
  int _page = 0;
  bool _keepOriginal = true;
  bool _notifyOnComplete = true;
  bool _autoConnect = true;
  bool _minimizeToTray = false;
  bool _testingConnection = false;
  bool _translationReloadPending = false;
  bool _reloadingCore = false;
  Future<bool>? _translationActivation;
  String _logLevel = '信息';
  String _appVersion = '3.7.0+1';

  void _applyConfig(Map<String, dynamic> config) {
    final app =
        (config['app'] as Map?)?.cast<String, dynamic>() ?? <String, dynamic>{};
    _appConfig = app;
    _fromController.text =
        (app['default_from_lang'] as String?)?.isNotEmpty == true
            ? app['default_from_lang'] as String
            : _fromController.text;
    _toController.text = (app['default_to_lang'] as String?)?.isNotEmpty == true
        ? app['default_to_lang'] as String
        : _toController.text;
    _directoryController.text = app['output_directory'] as String? ?? '';
    _trackExport = (app['track_export'] as num?)?.toInt() ?? _trackExport;
    _keepOriginal = app['keep_original'] as bool? ?? _keepOriginal;
    _notifyOnComplete = app['notify_on_complete'] as bool? ?? _notifyOnComplete;
    _autoConnect = app['auto_connect'] as bool? ?? _autoConnect;
    _minimizeToTray = app['minimize_to_tray'] as bool? ?? _minimizeToTray;
    _logLevel = app['log_level'] as String? ?? _logLevel;
    _concurrencyController.text =
        ((app['concurrency'] as num?)?.toInt() ?? 5).toString();

    _applyProviderConfig(config);
  }

  void _applyProviderConfig(Map<String, dynamic> config) {
    final providers = (config['providers'] as Map?)?.cast<String, dynamic>() ??
        <String, dynamic>{};
    final next = <String, Map<String, dynamic>>{};
    final configured = <String, bool>{};
    for (final entry in providers.entries) {
      final provider =
          (entry.value as Map?)?.cast<String, dynamic>() ?? <String, dynamic>{};
      final values = (provider['values'] as Map?)?.cast<String, dynamic>() ??
          <String, dynamic>{};
      next[entry.key] = values;
      configured[entry.key] = provider['configured'] as bool? ?? false;
      final controllers = _providerControllers.putIfAbsent(
          entry.key, () => <String, TextEditingController>{});
      for (final field in values.entries) {
        if (_isInternalProviderField(field.key)) continue;
        final controller =
            controllers.putIfAbsent(field.key, () => TextEditingController());
        controller.text = _displayValue(field.value);
      }
    }
    for (final id in _providerControllers.keys.toList()) {
      if (next.containsKey(id)) continue;
      for (final controller in _providerControllers.remove(id)!.values) {
        controller.dispose();
      }
    }
    _providerConfigs = next;
    _providerConfigured = configured;
  }

  String _displayValue(dynamic value) => value == null ? '' : '$value';

  bool _isInternalProviderField(String key) =>
      key == 'configured' || key == 'values';

  bool _isSecretField(String key) {
    final normalized = key.toLowerCase();
    return normalized.contains('key') ||
        normalized.contains('secret') ||
        normalized.contains('token') ||
        normalized.contains('password');
  }

  bool _isProviderNumberField(String key) =>
      key == 'qps' ||
      key == 'max_single_text_length' ||
      key == 'max_coroutine_num';

  String _providerFieldLabel(String key) {
    const labels = <String, String>{
      'ak_id': 'Access Key',
      'sk_key': 'Secret Key',
      'ak_secret': 'Access Secret',
      'secret_id': 'Secret ID',
      'secret_key': 'Secret Key',
      'app_id': '应用 ID',
      'service_tier': '服务版本',
      'app_key': '应用 ID',
      'app_secret': '应用密钥',
      'api_key': 'API Key',
      'base_url': '服务地址',
      'model': '模型名',
      'prompt': '翻译提示词',
      'api_secret': 'API Secret',
      'token': '访问令牌（Token）',
      'access_key': 'Access Key',
      'project_id': '项目 ID',
      'project_key': '项目 / 模型',
      'region': '区域',
      'data_id': '数据 ID',
      'max_single_text_length': '单次最大长度',
      'max_coroutine_num': '协程数量',
      'qps': 'QPS',
    };
    return labels[key] ?? key.replaceAll('_', ' ');
  }

  String _translatorName(String id) {
    for (final item in _translators.whereType<Map<String, dynamic>>()) {
      if (item['id'] == id) return item['name'] as String? ?? id;
    }
    return id;
  }

  Future<void> _editProviderParameters(String id) async {
    final dialogContext = _navigatorKey.currentContext;
    if (dialogContext == null) return;
    final controllers =
        _providerControllers[id] ?? <String, TextEditingController>{};
    final original = <String, String>{
      for (final entry in controllers.entries) entry.key: entry.value.text,
    };
    final saved = await showDialog<bool>(
      context: dialogContext,
      builder: (dialogContext) => AlertDialog(
        title: Text('${_translatorName(id)} 参数'),
        content: ConstrainedBox(
          constraints: BoxConstraints(
            maxWidth: 560,
            maxHeight: MediaQuery.sizeOf(dialogContext).height * 0.65,
          ),
          child: SingleChildScrollView(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisSize: MainAxisSize.min,
              children: [
                const Text('填写后点击“完成”，再点击翻译服务卡片右上角的“保存”。',
                    style: TextStyle(color: _muted, fontSize: 12)),
                const SizedBox(height: 14),
                _providerForm(id),
              ],
            ),
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(dialogContext).pop(false),
            child: const Text('取消'),
          ),
          FilledButton(
            onPressed: () => Navigator.of(dialogContext).pop(true),
            style: FilledButton.styleFrom(backgroundColor: _brand),
            child: const Text('完成'),
          ),
        ],
      ),
    );
    if (saved != true) {
      for (final entry in original.entries) {
        controllers[entry.key]?.text = entry.value;
      }
    }
    if (mounted) setState(() {});
  }

  @override
  void initState() {
    super.initState();
    _core = widget.core ?? CoreClient();
    _picker = widget.picker ?? const WindowsPicker();
    _coreExitSubscription = _core.onExit.listen((exit) {
      if (!mounted || _reloadingCore) return;
      setState(() => _status = exit.stderr.trim().isEmpty
          ? 'Go Core 已退出（${exit.code}）'
          : 'Go Core 已退出：${exit.stderr.trim().split('\n').last}');
      _scheduleReconnect();
    });
    _dropChannel.setMethodCallHandler(_handleDropCall);
    unawaited(_loadAppVersion());
    if (widget.autoConnect) _connect();
  }

  Future<void> _loadAppVersion() async {
    try {
      final version =
          await _windowChannel.invokeMethod<String>('getAppVersion');
      if (mounted && version != null && version.isNotEmpty) {
        setState(() => _appVersion = version);
      }
    } catch (_) {
      // The native channel is unavailable in widget tests.
    }
  }

  Future<dynamic> _handleDropCall(MethodCall call) async {
    if (call.method != 'filesDropped') return null;
    final paths = (call.arguments as List?)
            ?.whereType<String>()
            .map((path) => path.trim())
            .where((path) => path.isNotEmpty)
            .toList() ??
        <String>[];
    if (paths.isEmpty || !mounted) return null;
    final files = _parseFiles(_filesController.text).toSet();
    var inputDirectory = _inputDirectoryController.text.trim();
    for (final path in paths) {
      if (path.toLowerCase().endsWith('.srt')) {
        files.add(path);
      } else if (inputDirectory.isEmpty && Directory(path).existsSync()) {
        inputDirectory = path;
      }
    }
    setState(() {
      _filesController.text = files.join('\n');
      _inputDirectoryController.text = inputDirectory;
    });
    return null;
  }

  Future<bool> _connect({bool preserveDrafts = false}) async {
    try {
      await _core.start();
      final result = await _core.get('/api/v1/translators') as List<dynamic>;
      final config = await _core.get('/api/v1/config') as Map<String, dynamic>;
      if (!mounted) return false;
      setState(() {
        _translators = result;
        if (preserveDrafts) {
          _appConfig = (config['app'] as Map?)?.cast<String, dynamic>() ??
              <String, dynamic>{};
          _applyProviderConfig(config);
        } else {
          _applyConfig(config);
          final configuredTranslator =
              _appConfig['default_translator'] as String?;
          _selectedTranslator = _knownTranslatorOrFirst(configuredTranslator);
        }
        _ensureTranslatorDirection();
        _status = 'Go Core 已连接';
      });
      unawaited(_syncWindowBehavior());
      await _refreshTasks();
      await _eventsSubscription?.cancel();
      _eventsSubscription = _core.events().listen(_applyEvent,
          onError: (_) => _scheduleReconnect(), onDone: _scheduleReconnect);
      // Polling remains a slow recovery path if a proxy or a transient Core
      // restart interrupts the event stream.
      _refreshTimer?.cancel();
      _refreshTimer =
          Timer.periodic(const Duration(seconds: 10), (_) => _refreshTasks());
      return true;
    } catch (error) {
      if (mounted) setState(() => _status = '连接 Go Core 失败: $error');
      return false;
    }
  }

  bool _translatorIsUsable(String? id) {
    return isUsableTranslator(_translators, id);
  }

  String? _knownTranslatorOrFirst(String? id) {
    for (final translator in _translators.whereType<Map<String, dynamic>>()) {
      if (translator['id'] == id) return id;
    }
    for (final translator in _translators.whereType<Map<String, dynamic>>()) {
      if (translator['id'] is String) return translator['id'] as String;
    }
    return null;
  }

  Future<void> _syncWindowBehavior() async {
    try {
      await _windowChannel.invokeMethod<void>(
          'setMinimizeToTray', _minimizeToTray);
    } catch (_) {
      // The channel is only available on the Windows runner.
    }
  }

  List<Map<String, dynamic>> _availableTranslators() => _translators
      .whereType<Map>()
      .map((item) => item.cast<String, dynamic>())
      .where((item) => item['valid'] == true)
      .toList();

  void _applyEvent(Map<String, dynamic> event) {
    if (!mounted) return;
    final type = event['type'] as String? ?? '';
    if (type == 'tasks.snapshot') {
      final tasks = (event['tasks'] as List?)
              ?.whereType<Map>()
              .map((task) => task.cast<String, dynamic>())
              .toList() ??
          <Map<String, dynamic>>[];
      tasks.sort((a, b) => _taskDate(b).compareTo(_taskDate(a)));
      setState(() => _tasks = tasks);
      _checkPendingTranslationReload();
      return;
    }
    final task = (event['task'] as Map?)?.cast<String, dynamic>();
    final id = task?['id'] as String? ?? event['task_id'] as String?;
    if (id == null || task == null) return;
    final next = _tasks.whereType<Map<String, dynamic>>().toList();
    final index = next.indexWhere((item) => item['id'] == id);
    if (index < 0) {
      next.add(task);
    } else {
      next[index] = task;
    }
    next.sort((a, b) => _taskDate(b).compareTo(_taskDate(a)));
    setState(() {
      _tasks = next;
      if (type == 'task.failed') _status = '任务失败';
    });
    _checkPendingTranslationReload();
    if (type == 'task.failed') {
      final message = (task['error'] as String?)?.trim().isNotEmpty == true
          ? (task['error'] as String).trim()
          : (event['message'] as String?)?.trim() ?? 'Core 未提供失败原因';
      unawaited(_showMessage('翻译任务失败', message, error: true));
    } else if (type == 'task.completed' && _notifyOnComplete) {
      final outputDirectory = (task['output_directory'] as String?)?.trim();
      final message = outputDirectory == null || outputDirectory.isEmpty
          ? '字幕翻译已完成。'
          : '字幕翻译已完成，输出目录：$outputDirectory';
      unawaited(_showMessage('翻译完成', message));
    } else if (type == 'task.cancelled') {
      unawaited(_showMessage('任务已取消', '翻译任务已停止。'));
    }
  }

  void _scheduleReconnect() {
    if (!mounted ||
        _reloadingCore ||
        !_autoConnect ||
        _reconnectTimer != null) {
      return;
    }
    _reconnectTimer = Timer(const Duration(seconds: 2), () {
      _reconnectTimer = null;
      if (mounted && !_core.isRunning) _connect();
    });
  }

  Future<void> _refreshTasks() async {
    try {
      final result = await _core.get('/api/v1/tasks') as List<dynamic>;
      final tasks = result.whereType<Map<String, dynamic>>().toList()
        ..sort((a, b) => _taskDate(b).compareTo(_taskDate(a)));
      if (mounted) {
        setState(() => _tasks = tasks);
        _checkPendingTranslationReload();
      }
    } catch (_) {}
  }

  bool _hasActiveTasks(List<dynamic> tasks) =>
      tasks.whereType<Map>().any((task) => const {
            'queued',
            'reading',
            'translating',
            'writing'
          }.contains(task['status']));

  void _checkPendingTranslationReload() {
    if (!_translationReloadPending ||
        _reloadingCore ||
        _translationActivation != null ||
        _hasActiveTasks(_tasks)) {
      return;
    }
    unawaited(_activateTranslationWhenIdle().catchError((Object error) async {
      if (!mounted) return false;
      setState(() => _status = '翻译服务重新加载失败');
      await _showMessage('翻译服务重新加载失败', _friendlyError(error), error: true);
      return false;
    }));
  }

  Future<bool> _activateTranslationWhenIdle() {
    if (!_translationReloadPending) return Future.value(false);
    return _translationActivation ??= _activateTranslationWhenIdleOnce()
        .whenComplete(() => _translationActivation = null);
  }

  Future<bool> _activateTranslationWhenIdleOnce() async {
    final latest = await _core.get('/api/v1/tasks') as List<dynamic>;
    if (_hasActiveTasks(latest)) return false;
    if (!mounted) return false;
    setState(() {
      _reloadingCore = true;
      _status = '正在重新加载翻译服务';
    });
    _reconnectTimer?.cancel();
    _reconnectTimer = null;
    _refreshTimer?.cancel();
    _refreshTimer = null;
    try {
      // The event stream stays open until Core exits. Stop Core first so
      // cancelling the subscription cannot wait forever for the next event.
      await _core.stop();
      await _eventsSubscription?.cancel();
      _eventsSubscription = null;
      if (!await _connect(preserveDrafts: true)) {
        throw StateError('Core 重新启动失败，请重新打开软件后检查翻译服务设置。');
      }
      if (mounted) {
        setState(() {
          _translationReloadPending = false;
          _status = '翻译服务已生效';
        });
      }
      return true;
    } catch (_) {
      _translationReloadPending = false;
      rethrow;
    } finally {
      if (mounted) {
        setState(() => _reloadingCore = false);
      } else {
        _reloadingCore = false;
      }
    }
  }

  DateTime _taskDate(Map<String, dynamic> task) {
    return DateTime.tryParse(task['updated_at'] as String? ??
            task['created_at'] as String? ??
            '') ??
        DateTime.fromMillisecondsSinceEpoch(0);
  }

  void _resetTaskForm() {
    _filesController.clear();
    _inputDirectoryController.clear();
    _fromController.text =
        (_appConfig['default_from_lang'] as String?)?.trim().isNotEmpty == true
            ? _appConfig['default_from_lang'] as String
            : 'en';
    _toController.text =
        (_appConfig['default_to_lang'] as String?)?.trim().isNotEmpty == true
            ? _appConfig['default_to_lang'] as String
            : 'zh';
    final defaultTranslator = _appConfig['default_translator'] as String?;
    _selectedTranslator = _knownTranslatorOrFirst(defaultTranslator);
    _mode = '全量翻译';
    _ensureTranslatorDirection();
  }

  Future<void> _createTask() async {
    if (_reloadingCore) return;
    final files = _filesController.text
        .split(RegExp(r'[;\r\n]+'))
        .map((value) => value.trim())
        .where(
            (value) => value.isNotEmpty && value.toLowerCase().endsWith('.srt'))
        .toList();
    final inputDirectory = _inputDirectoryController.text.trim();
    if (files.isEmpty && inputDirectory.isEmpty) {
      await _showMessage('缺少输入文件', '请填写 SRT 文件路径，或选择包含 SRT 文件的目录。',
          error: true);
      return;
    }
    if (_selectedTranslator == null) {
      await _showMessage('未选择翻译引擎', '请先在翻译引擎下拉框中选择一个可用引擎。', error: true);
      return;
    }
    if (!_translatorIsUsable(_selectedTranslator)) {
      await _showMessage('翻译引擎不可用', '当前引擎尚未配置或已被禁用，请先在设置页完成配置。', error: true);
      return;
    }
    if (!isLanguagePairValid(_fromController.text, _toController.text)) {
      await _showMessage('语言设置无效', '源语言和目标语言不能相同。', error: true);
      return;
    }
    if (_selectedTranslator == 'caiyunai' &&
        !_caiyunDirectionSupported(_fromController.text, _toController.text)) {
      await _showMessage('语言设置无效', '彩云小译暂不支持所选翻译方向。', error: true);
      return;
    }
    try {
      await _core.post('/api/v1/tasks', body: {
        'files': files,
        'directory': inputDirectory,
        'translator_id': _selectedTranslator,
        'from_lang': _fromController.text.trim(),
        'to_lang': _toController.text.trim(),
        'translate_mode': _mode,
        'main_track_report': '来源语种',
        'flag_track_export': _trackExport,
        'output_directory': _directoryController.text.trim(),
      });
      if (mounted) {
        setState(() {
          _resetTaskForm();
          _status = '任务已提交';
        });
      }
      await _refreshTasks();
      await _showMessage('任务已提交', '翻译任务已加入队列，任务列表会自动刷新。');
    } catch (error) {
      if (mounted) setState(() => _status = '提交任务失败');
      await _showMessage('提交任务失败', _friendlyError(error), error: true);
    }
  }

  Future<void> _chooseFiles() async {
    try {
      final selected = await _picker.pickFiles();
      if (selected.isEmpty || !mounted) return;
      final current = _parseFiles(_filesController.text);
      final merged = <String>{...current, ..._parseFiles(selected.join('\n'))};
      _filesController.text = merged.join('\n');
      setState(() {});
    } catch (error) {
      await _showMessage('选择文件失败', _friendlyError(error), error: true);
    }
  }

  Future<void> _chooseInputDirectory() async {
    try {
      final selected = await _picker.pickDirectory();
      if (selected == null || selected.isEmpty || !mounted) return;
      setState(() => _inputDirectoryController.text = selected);
    } catch (error) {
      await _showMessage('选择目录失败', _friendlyError(error), error: true);
    }
  }

  Future<void> _chooseOutputDirectory() async {
    try {
      final selected = await _picker.pickDirectory();
      if (selected == null || selected.isEmpty || !mounted) return;
      setState(() => _directoryController.text = selected);
    } catch (error) {
      await _showMessage('选择输出目录失败', _friendlyError(error), error: true);
    }
  }

  Directory get _logsDirectory => Directory(
      '${File(Platform.resolvedExecutable).parent.path}${Platform.pathSeparator}logs');

  Future<void> _openLogs() async {
    final directory = _logsDirectory;
    if (!directory.existsSync()) {
      await _showMessage('暂无日志', '日志目录尚未创建。');
      return;
    }
    await Process.start('explorer.exe', <String>[directory.path]);
  }

  Future<void> _clearOldLogs() async {
    final dialogContext = _navigatorKey.currentContext;
    if (dialogContext == null) return;
    final confirmed = await showDialog<bool>(
            context: dialogContext,
            builder: (dialogContext) => AlertDialog(
                  title: const Text('清理历史日志'),
                  content: const Text('将删除当天之前的日志文件，是否继续？'),
                  actions: [
                    TextButton(
                        onPressed: () => Navigator.of(dialogContext).pop(false),
                        child: const Text('取消')),
                    FilledButton(
                        onPressed: () => Navigator.of(dialogContext).pop(true),
                        style: FilledButton.styleFrom(backgroundColor: _brand),
                        child: const Text('清理'))
                  ],
                )) ??
        false;
    if (!confirmed) return;
    final directory = _logsDirectory;
    if (!directory.existsSync()) {
      await _showMessage('日志清理完成', '暂无历史日志文件。');
      return;
    }
    final today = DateTime.now();
    var removed = 0;
    await for (final entity in directory.list()) {
      if (entity is! File || !entity.path.toLowerCase().endsWith('.log')) {
        continue;
      }
      final modified = await entity.lastModified();
      if (modified.year == today.year &&
          modified.month == today.month &&
          modified.day == today.day) {
        continue;
      }
      try {
        await entity.delete();
        removed++;
      } catch (_) {
        // A file may be held by the running Core; leave it for the next pass.
      }
    }
    await _showMessage('日志清理完成', '已删除 $removed 个历史日志文件。');
  }

  List<String> _parseFiles(String text) => parseSrtPaths(text);

  List<Map<String, dynamic>> _languagesForTranslator() {
    for (final item in _translators.whereType<Map<String, dynamic>>()) {
      if (item['id'] != _selectedTranslator) continue;
      return (item['languages'] as List?)
              ?.whereType<Map>()
              .map((language) => language.cast<String, dynamic>())
              .where((language) => language['id'] is String)
              .toList() ??
          <Map<String, dynamic>>[];
    }
    return _fallbackLanguages
        .map((language) => Map<String, dynamic>.from(language))
        .toList();
  }

  bool _caiyunDirectionSupported(String from, String to) {
    if (from == to) return false;
    if (from == 'zh' || from == 'zh-Hant') {
      return const {'zh', 'zh-Hant', 'en', 'ja', 'ko'}.contains(to);
    }
    return const {
          'en',
          'ja',
          'ko',
          'de',
          'es',
          'fr',
          'it',
          'pt',
          'ru',
          'tr',
          'vi'
        }.contains(from) &&
        const {'zh', 'zh-Hant'}.contains(to);
  }

  void _ensureTranslatorDirection() {
    if (_selectedTranslator == 'caiyunai' &&
        !_caiyunDirectionSupported(_fromController.text, _toController.text)) {
      _fromController.text = 'en';
      _toController.text = 'zh';
    }
    if (_selectedTranslator == 'deepl_pro' ||
        _selectedTranslator == 'youdao_batch') {
      final languages = _languagesForTranslator()
          .map((language) => language['id'])
          .whereType<String>()
          .toSet();
      if (!languages.contains(_fromController.text) ||
          !languages.contains(_toController.text) ||
          _fromController.text == _toController.text) {
        _fromController.text = 'en';
        _toController.text = 'zh';
      }
    }
  }

  Widget _languageField(String label, TextEditingController controller) {
    final languages = _languagesForTranslator()
        .where((language) =>
            controller != _toController ||
            _selectedTranslator != 'caiyunai' ||
            _caiyunDirectionSupported(
                _fromController.text, language['id'] as String))
        .toList();
    final selected =
        languages.any((language) => language['id'] == controller.text)
            ? controller.text
            : null;
    return _field(
        label,
        DropdownButtonFormField<String>(
            value: selected,
            isExpanded: true,
            hint: Text(controller.text),
            items: languages
                .map((language) => DropdownMenuItem<String>(
                    value: language['id'] as String,
                    child: Text(language['name'] as String? ??
                        language['id'] as String)))
                .toList(),
            onChanged: (value) {
              if (value != null) {
                setState(() {
                  controller.text = value;
                  if (controller == _fromController &&
                      _selectedTranslator == 'caiyunai' &&
                      !_caiyunDirectionSupported(value, _toController.text)) {
                    _toController.text = value == 'zh' ? 'en' : 'zh';
                  }
                });
              }
            }));
  }

  Future<void> _cancelTask(String id) async {
    try {
      await _core.post('/api/v1/tasks/$id/cancel');
      await _refreshTasks();
      await _showMessage('任务取消中', '已发送取消请求，任务状态会在列表中更新。');
    } catch (error) {
      if (mounted) setState(() => _status = '取消任务失败');
      await _showMessage('取消任务失败', _friendlyError(error), error: true);
    }
  }

  Future<void> _showMessage(String title, String message,
      {bool error = false}) async {
    final dialogContext = _navigatorKey.currentContext;
    if (!mounted || dialogContext == null) return;
    await showDialog<void>(
      context: dialogContext,
      builder: (dialogContext) => AlertDialog(
        title: Row(children: [
          Icon(error ? Icons.error_outline : Icons.check_circle_outline,
              color: error ? const Color(0xffd9534f) : _brand),
          const SizedBox(width: 9),
          Expanded(child: Text(title)),
        ]),
        content: Text(message),
        actions: [
          FilledButton(
              onPressed: () => Navigator.of(dialogContext).pop(),
              style: FilledButton.styleFrom(backgroundColor: _brand),
              child: const Text('知道了'))
        ],
      ),
    );
  }

  String _friendlyError(Object error) {
    final candidates = <String>[
      if (error is HttpException) error.message,
      error.toString(),
    ];
    for (final candidate in candidates) {
      final text = candidate.trim();
      final start = text.indexOf('{');
      final end = text.lastIndexOf('}');
      if (start >= 0 && end > start) {
        try {
          final decoded = jsonDecode(text.substring(start, end + 1));
          final apiError = decoded is Map ? decoded['error'] : null;
          final message = apiError is Map ? apiError['message'] : null;
          if (message is String && message.isNotEmpty) return message;
        } catch (_) {
          // Try the next representation when the body is not JSON.
        }
      }
      if (text.isNotEmpty && !text.startsWith('HttpException:')) {
        return text.startsWith('Exception: ') ? text.substring(11) : text;
      }
    }
    return error.toString();
  }

  Future<void> _saveSettings(String section) async {
    var written = false;
    try {
      final app = <String, dynamic>{};
      final providers = <String, dynamic>{};
      if (section == 'translation') {
        app['default_translator'] = _selectedTranslator ?? '';
      } else if (section == 'output') {
        app.addAll({
          'output_directory': _directoryController.text.trim(),
          'track_export': _trackExport,
          'keep_original': _keepOriginal,
          'notify_on_complete': _notifyOnComplete,
        });
      } else if (section == 'app') {
        final concurrency = int.tryParse(_concurrencyController.text.trim());
        if (concurrency == null || concurrency < 1 || concurrency > 256) {
          await _showMessage('设置无效', '并发数量必须是 1 到 256 之间的整数。', error: true);
          return;
        }
        app.addAll({
          'auto_connect': _autoConnect,
          'minimize_to_tray': _minimizeToTray,
          'log_level': _logLevel,
          'concurrency': concurrency,
        });
      }
      for (final entry in section == 'translation'
          ? _providerControllers.entries
          : <MapEntry<String, Map<String, TextEditingController>>>[]) {
        final values = <String, dynamic>{};
        final original = _providerConfigs[entry.key] ?? <String, dynamic>{};
        for (final field in entry.value.entries) {
          final key = field.key;
          if (_isInternalProviderField(key)) continue;
          final text = field.value.text.trim();
          if (text == _displayValue(original[key]).trim()) continue;
          if (_isProviderNumberField(key)) {
            final value = int.tryParse(text);
            final valid = value != null &&
                value > 0 &&
                (key != 'qps' || value <= 100000) &&
                (key != 'max_single_text_length' || value <= 1000000) &&
                (entry.key != 'tencent_cloud_mt' ||
                    key != 'max_single_text_length' ||
                    value <= 2000) &&
                (entry.key != 'deepseek' ||
                    key != 'max_single_text_length' ||
                    value <= 6000) &&
                (key != 'max_coroutine_num' || value <= 1024);
            if (!valid) {
              await _showMessage(
                  '设置无效', '${_providerFieldLabel(key)} 必须是有效的正整数，且不能超过允许范围。',
                  error: true);
              return;
            }
            values[key] = value;
          } else {
            values[key] = text;
          }
        }
        if (values.isNotEmpty) providers[entry.key] = values;
      }
      await _core.put('/api/v1/config', {'app': app, 'providers': providers});
      written = true;
      final saved = await _core.get('/api/v1/config') as Map<String, dynamic>;
      final translators =
          await _core.get('/api/v1/translators') as List<dynamic>;
      if (mounted) {
        setState(() {
          _appConfig =
              (saved['app'] as Map?)?.cast<String, dynamic>() ?? _appConfig;
          if (section == 'translation') _applyProviderConfig(saved);
          _translators = translators;
        });
      }
      if (section == 'app') unawaited(_syncWindowBehavior());
      var activated = false;
      if (section == 'translation' && mounted) {
        setState(() {
          _translationReloadPending = true;
          _status = '翻译服务已保存，等待生效';
        });
        activated = await _activateTranslationWhenIdle();
      } else if (mounted) {
        setState(() => _status = '设置已保存');
      }
      final label = switch (section) {
        'translation' => '翻译服务',
        'output' => '输出设置',
        _ => '应用设置',
      };
      await _showMessage(
          '设置已保存',
          section == 'translation'
              ? activated
                  ? '$label 已保存并生效。'
                  : '$label 已保存，当前任务完成后会自动生效。'
              : '$label 已写入 cfg.yml，并已重新读取验证。');
    } catch (error) {
      final title = written ? '设置已写入，但重新加载失败' : '设置保存失败';
      if (mounted) setState(() => _status = title);
      await _showMessage(title, _friendlyError(error), error: true);
    }
  }

  Future<void> _testConnection() async {
    final translatorId = _selectedTranslator;
    if (translatorId == null || _testingConnection) return;
    final configuredRegion =
        _providerControllers[translatorId]?['region']?.text.trim() ?? '';
    final region = translatorId == 'baidu_cloud_mt' ||
            translatorId == 'caiyunai' ||
            translatorId == 'deepl_pro' ||
            translatorId == 'deepseek' ||
            translatorId == 'youdao_batch'
        ? ''
        : configuredRegion.isEmpty
            ? (translatorId == 'tencent_cloud_mt'
                ? 'ap-guangzhou'
                : 'cn-hangzhou')
            : configuredRegion;
    setState(() => _testingConnection = true);
    try {
      final result = await _core.post(
              '/api/v1/translators/$translatorId/test-connection',
              body: region.isEmpty ? <String, dynamic>{} : {'region': region})
          as Map<String, dynamic>;
      final endpoint = result['endpoint'] as String? ?? '';
      final latency = (result['latency_ms'] as num?)?.toInt();
      final displayEndpoint =
          translatorId == 'deepseek' ? endpoint : '$endpoint:443';
      await _showMessage('引擎网络可达',
          '已连接 $displayEndpoint${latency == null ? '' : '，耗时 $latency 毫秒'}。此检查不验证访问密钥或翻译额度。');
    } catch (error) {
      await _showMessage('引擎连接失败', _friendlyError(error), error: true);
    } finally {
      if (mounted) setState(() => _testingConnection = false);
    }
  }

  Future<void> _openTaskDirectory(Map<String, dynamic> task) async {
    final currentFile = task['current_file'] as String?;
    final outputDirectory = task['output_directory'] as String?;
    final directory = outputDirectory ??
        (currentFile == null ? null : Directory(currentFile).parent.path);
    if (directory == null || directory.isEmpty) {
      await _showMessage('目录不可用', '任务没有记录可打开的生成目录。', error: true);
      return;
    }
    try {
      await Process.run('explorer.exe', [directory]);
    } catch (error) {
      if (mounted) setState(() => _status = '打开目录失败');
      await _showMessage('打开目录失败', _friendlyError(error), error: true);
    }
  }

  @override
  void dispose() {
    _refreshTimer?.cancel();
    _reconnectTimer?.cancel();
    final eventsSubscription = _eventsSubscription;
    if (eventsSubscription != null) unawaited(eventsSubscription.cancel());
    final exitSubscription = _coreExitSubscription;
    if (exitSubscription != null) unawaited(exitSubscription.cancel());
    _filesController.dispose();
    _inputDirectoryController.dispose();
    _directoryController.dispose();
    _fromController.dispose();
    _toController.dispose();
    _concurrencyController.dispose();
    for (final provider in _providerControllers.values) {
      for (final controller in provider.values) {
        controller.dispose();
      }
    }
    unawaited(_core.dispose());
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      navigatorKey: _navigatorKey,
      debugShowCheckedModeBanner: false,
      title: 'Anto-你的字幕翻译小助手',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: _brand),
        scaffoldBackgroundColor: const Color(0xfff5f7fb),
        fontFamily: 'Segoe UI',
        filledButtonTheme: FilledButtonThemeData(
          style: FilledButton.styleFrom(
            backgroundColor: _brand,
            foregroundColor: Colors.white,
            disabledBackgroundColor: const Color(0xffe8eaf1),
            disabledForegroundColor: const Color(0xffa0a7b6),
            minimumSize: const Size(0, 36),
            padding: const EdgeInsets.symmetric(horizontal: 18, vertical: 9),
            shape:
                RoundedRectangleBorder(borderRadius: BorderRadius.circular(10)),
          ),
        ),
        outlinedButtonTheme: OutlinedButtonThemeData(
          style: ButtonStyle(
            backgroundColor: MaterialStateProperty.resolveWith((states) =>
                states.contains(MaterialState.hovered)
                    ? const Color(0xffe5e7ff)
                    : const Color(0xfff0f1ff)),
            foregroundColor: const MaterialStatePropertyAll(_brand),
            side: const MaterialStatePropertyAll(BorderSide.none),
            minimumSize: const MaterialStatePropertyAll(Size(0, 36)),
            padding: const MaterialStatePropertyAll(
                EdgeInsets.symmetric(horizontal: 16, vertical: 8)),
            shape: MaterialStatePropertyAll(RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(10))),
          ),
        ),
        textButtonTheme: TextButtonThemeData(
          style: TextButton.styleFrom(
            foregroundColor: _brand,
            minimumSize: const Size(0, 34),
            shape:
                RoundedRectangleBorder(borderRadius: BorderRadius.circular(9)),
          ),
        ),
        iconButtonTheme: IconButtonThemeData(
          style: IconButton.styleFrom(
            foregroundColor: _brand,
            hoverColor: const Color(0xffeceeff),
            shape:
                RoundedRectangleBorder(borderRadius: BorderRadius.circular(9)),
          ),
        ),
        inputDecorationTheme: InputDecorationTheme(
          isDense: true,
          contentPadding:
              const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
          filled: true,
          fillColor: Colors.white,
          border: OutlineInputBorder(
              borderRadius: BorderRadius.circular(9),
              borderSide: const BorderSide(color: Color(0xffdfe3ee))),
          enabledBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(9),
              borderSide: const BorderSide(color: Color(0xffdfe3ee))),
          focusedBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(9),
              borderSide:
                  const BorderSide(color: Color(0xff999bf2), width: 1.4)),
        ),
      ),
      home: _buildShell(),
    );
  }

  Widget _buildShell() {
    return Scaffold(
      body: LayoutBuilder(builder: (context, constraints) {
        final compact = constraints.maxWidth < 700;
        final padding =
            EdgeInsets.fromLTRB(compact ? 14 : 40, 30, compact ? 14 : 40, 44);
        return Row(crossAxisAlignment: CrossAxisAlignment.start, children: [
          _buildSidebar(compact),
          Expanded(
              child: Padding(
                  padding: padding,
                  child: LayoutBuilder(builder: (context, content) {
                    if (_page == 3) return AboutPage(version: _appVersion);
                    if (_page == 2) {
                      return _buildSettingsPage(
                          content.maxWidth, content.maxHeight);
                    }
                    if (_page == 1) return _buildQueuePage();
                    return _buildWorkspacePage(
                        content.maxWidth, content.maxHeight);
                  }))),
        ]);
      }),
    );
  }

  Widget _buildSidebar(bool compact) {
    final labels = ['字幕翻译', '任务队列', '我的设置', '关于应用'];
    final icons = [
      Icons.subtitles_outlined,
      Icons.list_alt_outlined,
      Icons.settings_outlined,
      Icons.info_outline
    ];
    return Container(
      width: compact ? 76 : 224,
      color: _sidebar,
      padding:
          EdgeInsets.fromLTRB(compact ? 10 : 16, 24, compact ? 10 : 16, 14),
      child: Column(children: [
        compact ? const LogoMark(size: 38) : const LogoLockup(),
        const SizedBox(height: 28),
        if (!compact)
          const Align(
              alignment: Alignment.centerLeft,
              child: Padding(
                  padding: EdgeInsets.only(left: 12, bottom: 9),
                  child: Text('工作台',
                      style: TextStyle(
                          color: Color(0xff777d9c),
                          fontSize: 11,
                          letterSpacing: 1.1)))),
        Expanded(
            child: SingleChildScrollView(
                child: Column(
                    children: List.generate(
                        labels.length,
                        (index) => Padding(
                            padding: const EdgeInsets.only(bottom: 6),
                            child: _navButton(labels[index], icons[index],
                                index, compact)))))),
        if (!compact)
          Container(
              width: double.infinity,
              padding: const EdgeInsets.fromLTRB(12, 13, 12, 2),
              decoration: const BoxDecoration(
                  border: Border(top: BorderSide(color: Color(0xff343951)))),
              child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(_status,
                        style: const TextStyle(
                            color: Color(0xff858ba7), fontSize: 12)),
                    const SizedBox(height: 4),
                    const Text('● 本地服务运行中',
                        style:
                            TextStyle(color: Color(0xff5fc79b), fontSize: 12))
                  ]))
        else
          const Icon(Icons.circle, size: 7, color: Color(0xff5fc79b)),
      ]),
    );
  }

  Widget _navButton(String label, IconData icon, int index, bool compact) {
    final active = _page == index;
    return Material(
        color: Colors.transparent,
        child: InkWell(
          borderRadius: BorderRadius.circular(10),
          onTap: () => setState(() => _page = index),
          child: Container(
            width: double.infinity,
            padding: const EdgeInsets.symmetric(vertical: 11, horizontal: 12),
            decoration: BoxDecoration(
                borderRadius: BorderRadius.circular(10),
                gradient: active
                    ? const LinearGradient(
                        colors: [Color(0xff5e61e9), Color(0xff6669ed)])
                    : null,
                boxShadow: active
                    ? const [
                        BoxShadow(
                            color: Color(0x405356e0),
                            blurRadius: 16,
                            offset: Offset(0, 7))
                      ]
                    : null),
            child: Row(
                mainAxisAlignment: compact
                    ? MainAxisAlignment.center
                    : MainAxisAlignment.start,
                children: [
                  Icon(icon,
                      size: 18,
                      color: active ? Colors.white : const Color(0xffb9bed3)),
                  if (!compact) ...[
                    const SizedBox(width: 11),
                    Text(label,
                        style: TextStyle(
                            color:
                                active ? Colors.white : const Color(0xffb9bed3),
                            fontWeight:
                                active ? FontWeight.w600 : FontWeight.w400))
                  ]
                ]),
          ),
        ));
  }

  Widget _buildWorkspacePage(double width, double height) {
    final wide = width >= 700;
    final composer = _card(
        Column(children: [
          Expanded(child: SingleChildScrollView(child: _buildComposer())),
          _buildComposerActions(),
        ]),
        key: const Key('workspace-composer-card'));
    final recent = _card(SingleChildScrollView(child: _buildRecentTasks()),
        key: const Key('workspace-recent-card'));
    if (wide) {
      return Row(crossAxisAlignment: CrossAxisAlignment.stretch, children: [
        Expanded(child: composer),
        const SizedBox(width: 20),
        SizedBox(width: 360, child: recent),
      ]);
    }
    return SingleChildScrollView(
        child: Column(children: [
      SizedBox(height: height, child: composer),
      const SizedBox(height: 20),
      SizedBox(height: height, child: recent),
    ]));
  }

  Widget _buildComposer() {
    final selectedFiles = _parseFiles(_filesController.text);
    return Padding(
        padding: const EdgeInsets.fromLTRB(22, 22, 22, 0),
        child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
          const Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text('新建翻译任务',
                    style: TextStyle(
                        color: _ink,
                        fontWeight: FontWeight.w700,
                        fontSize: 15)),
                Text('支持 .srt', style: TextStyle(color: _muted, fontSize: 12))
              ]),
          const SizedBox(height: 17),
          Container(
              constraints: const BoxConstraints(minHeight: 185),
              width: double.infinity,
              decoration: BoxDecoration(
                  border: Border.all(
                      color: const Color(0xffbdc2f4), style: BorderStyle.solid),
                  borderRadius: BorderRadius.circular(13),
                  color: const Color(0xfff5f5ff)),
              child: Padding(
                  padding: const EdgeInsets.all(16),
                  child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Container(
                            width: 48,
                            height: 48,
                            decoration: BoxDecoration(
                                color: const Color(0xffe1e2ff),
                                borderRadius: BorderRadius.circular(14)),
                            child: const Icon(Icons.file_upload_outlined,
                                color: _brand, size: 24)),
                        const SizedBox(height: 10),
                        Text(
                            selectedFiles.isEmpty
                                ? '选择 SRT 文件或目录'
                                : '已选择 ${selectedFiles.length} 个 SRT 文件',
                            style: const TextStyle(
                                fontWeight: FontWeight.w700, fontSize: 15)),
                        const SizedBox(height: 4),
                        const Text('支持拖放 SRT 文件或目录，也可以直接粘贴路径',
                            style: TextStyle(color: _muted, fontSize: 12)),
                        const SizedBox(height: 12),
                        Row(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              OutlinedButton(
                                  onPressed: _chooseFiles,
                                  child: const Text('选择文件')),
                              const SizedBox(width: 9),
                              OutlinedButton(
                                  onPressed: _chooseInputDirectory,
                                  child: const Text('选择目录'))
                            ])
                      ]))),
          if (selectedFiles.isNotEmpty) ...[
            const SizedBox(height: 12),
            _selectedFilesList(selectedFiles),
          ],
          const SizedBox(height: 17),
          _responsiveFields(),
          if (_selectedTranslator == 'caiyunai') ...[
            const SizedBox(height: 8),
            const Text('数据来自彩云小译',
                style: TextStyle(color: _muted, fontSize: 12)),
          ],
        ]));
  }

  Widget _buildComposerActions() => Padding(
      padding: const EdgeInsets.fromLTRB(22, 14, 22, 22),
      child: Row(mainAxisAlignment: MainAxisAlignment.end, children: [
        OutlinedButton(
            onPressed: () => setState(_resetTaskForm), child: const Text('重置')),
        const SizedBox(width: 9),
        FilledButton(
            onPressed:
                !_reloadingCore && _translatorIsUsable(_selectedTranslator)
                    ? _createTask
                    : null,
            style: FilledButton.styleFrom(
                backgroundColor: _brand,
                padding:
                    const EdgeInsets.symmetric(horizontal: 16, vertical: 11)),
            child: const Text('开始翻译'))
      ]));

  Widget _selectedFilesList(List<String> files) {
    return Container(
        width: double.infinity,
        padding: const EdgeInsets.fromLTRB(12, 10, 12, 8),
        decoration: BoxDecoration(
            color: const Color(0xfff8f9ff),
            border: Border.all(color: const Color(0xffe5e7f4)),
            borderRadius: BorderRadius.circular(10)),
        child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
          Text('已选文件（${files.length}）',
              style: const TextStyle(
                  color: _ink, fontSize: 12, fontWeight: FontWeight.w600)),
          const SizedBox(height: 5),
          SizedBox(
              height: (files.length * 48.0).clamp(48.0, 192.0),
              child: ListView.builder(
                  itemCount: files.length,
                  itemBuilder: (context, index) {
                    final path = files[index];
                    final name = path.split(RegExp(r'[\\/]')).last;
                    return SizedBox(
                        height: 48,
                        child: Row(children: [
                          const Icon(Icons.description_outlined,
                              color: _brand, size: 18),
                          const SizedBox(width: 8),
                          Expanded(
                              child: Tooltip(
                                  message: path,
                                  child: Column(
                                      crossAxisAlignment:
                                          CrossAxisAlignment.start,
                                      mainAxisAlignment:
                                          MainAxisAlignment.center,
                                      children: [
                                        Text(name,
                                            maxLines: 1,
                                            overflow: TextOverflow.ellipsis,
                                            style: const TextStyle(
                                                fontWeight: FontWeight.w600,
                                                fontSize: 12)),
                                        Text(path,
                                            maxLines: 1,
                                            overflow: TextOverflow.ellipsis,
                                            style: const TextStyle(
                                                color: _muted, fontSize: 11)),
                                      ]))),
                          IconButton(
                              tooltip: '移除 $name',
                              onPressed: () {
                                final remaining = List<String>.from(files)
                                  ..removeAt(index);
                                setState(() => _filesController.text =
                                    remaining.join('\n'));
                              },
                              icon: const Icon(Icons.close, size: 17))
                        ]));
                  }))
        ]));
  }

  Widget _responsiveFields() {
    return LayoutBuilder(builder: (context, constraints) {
      final oneColumn = constraints.maxWidth < 400;
      final children = [
        _languageField('源语言', _fromController),
        _languageField('目标语言', _toController),
        _field(
            '翻译引擎',
            DropdownButtonFormField<String>(
                isExpanded: true,
                value: _translatorIsUsable(_selectedTranslator)
                    ? _selectedTranslator
                    : null,
                hint: const Text('暂无可用引擎'),
                items: _availableTranslators().map((value) {
                  return DropdownMenuItem(
                      value: value['id'] as String,
                      child: Text(value['name'] as String,
                          overflow: TextOverflow.ellipsis));
                }).toList(),
                onChanged: (value) => setState(() {
                      _selectedTranslator = value;
                      _ensureTranslatorDirection();
                    }))),
        _field(
            '翻译模式',
            DropdownButtonFormField<String>(
                value: _mode,
                items: const [
                  DropdownMenuItem(value: '全量翻译', child: Text('全量翻译')),
                  DropdownMenuItem(value: '增量翻译', child: Text('增量翻译'))
                ],
                onChanged: (value) => setState(() => _mode = value ?? _mode))),
        _field(
            '输入目录',
            TextField(
                controller: _inputDirectoryController,
                decoration: InputDecoration(
                    hintText: '可选，扫描目录中的 .srt',
                    suffixIcon: IconButton(
                        tooltip: '选择目录',
                        onPressed: _chooseInputDirectory,
                        icon: const Icon(Icons.folder_open_outlined))))),
      ];
      final fieldWidth =
          oneColumn ? constraints.maxWidth : (constraints.maxWidth - 13) / 2;
      return Wrap(
          spacing: 13,
          runSpacing: 8,
          children: children
              .map((child) => SizedBox(width: fieldWidth, child: child))
              .toList());
    });
  }

  Widget _field(String label, Widget child) =>
      Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
        Text(label, style: const TextStyle(color: _muted, fontSize: 12)),
        const SizedBox(height: 6),
        SizedBox(height: 48, child: child)
      ]);

  Widget _buildRecentTasks() {
    return Padding(
        padding: const EdgeInsets.all(22),
        child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
          Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
            const Text('最近任务',
                style: TextStyle(
                    color: _ink, fontWeight: FontWeight.w700, fontSize: 15)),
            TextButton(
                onPressed: () => setState(() => _page = 1),
                style: TextButton.styleFrom(
                    foregroundColor: _muted,
                    padding: EdgeInsets.zero,
                    minimumSize: Size.zero,
                    tapTargetSize: MaterialTapTargetSize.shrinkWrap),
                child: const Text('查看全部'))
          ]),
          const SizedBox(height: 10),
          _taskList(limit: 3)
        ]));
  }

  Widget _buildQueuePage() {
    return _card(
        Padding(
            padding: const EdgeInsets.all(20),
            child:
                Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
              Row(children: [
                const Text('任务记录',
                    style:
                        TextStyle(fontWeight: FontWeight.w700, fontSize: 15)),
                const SizedBox(width: 8),
                const Text('自动刷新',
                    style: TextStyle(color: _muted, fontSize: 12)),
                const Spacer(),
                FilledButton(
                    onPressed: () => setState(() => _page = 0),
                    child: const Text('新建任务')),
              ]),
              const SizedBox(height: 17),
              Expanded(child: SingleChildScrollView(child: _taskList())),
            ])),
        key: const Key('queue-tasks-card'));
  }

  Widget _taskList({int? limit}) {
    final tasks = limit == null ? _tasks : _tasks.take(limit).toList();
    if (tasks.isEmpty) {
      return const Padding(
          padding: EdgeInsets.symmetric(vertical: 24),
          child: Center(child: Text('暂无任务', style: TextStyle(color: _muted))));
    }
    return Column(
        children: tasks
            .map((item) => _taskRow(item as Map<String, dynamic>))
            .toList());
  }

  String _taskStatusLabel(String status) {
    switch (status) {
      case 'queued':
        return '等待中';
      case 'reading':
        return '读取中';
      case 'translating':
        return '翻译中';
      case 'writing':
        return '写入中';
      case 'completed':
        return '已完成';
      case 'failed':
        return '失败';
      case 'cancelled':
        return '已取消';
      default:
        return status;
    }
  }

  String _taskFileLabel(Map<String, dynamic> task) {
    final currentFile = (task['current_file'] as String?)?.trim();
    if (currentFile != null && currentFile.isNotEmpty) {
      return currentFile.split(RegExp(r'[\\/]')).last;
    }
    final outputFiles = task['output_files'] as List?;
    if (outputFiles != null && outputFiles.isNotEmpty) {
      final outputFile = outputFiles.first;
      if (outputFile is String && outputFile.trim().isNotEmpty) {
        return outputFile.split(RegExp(r'[\\/]')).last;
      }
    }
    return task['status'] == 'queued' ? '等待处理' : '翻译任务';
  }

  Widget _taskRow(Map<String, dynamic> task) {
    final status = task['status'] as String? ?? 'queued';
    final finished =
        status == 'completed' || status == 'failed' || status == 'cancelled';
    final progress = ((task['completed'] as num? ?? 0) /
            ((task['total'] as num? ?? 1) == 0 ? 1 : (task['total'] as num)))
        .clamp(0, 1)
        .toDouble();
    final errorMessage = (task['error'] as String?)?.trim() ?? '';
    return Container(
        padding: const EdgeInsets.symmetric(vertical: 12),
        decoration: const BoxDecoration(
            border: Border(bottom: BorderSide(color: Color(0xffeff0f5)))),
        child: Row(crossAxisAlignment: CrossAxisAlignment.start, children: [
          Container(
              width: 34,
              height: 34,
              decoration: BoxDecoration(
                  color: const Color(0xfff0f1ff),
                  borderRadius: BorderRadius.circular(9)),
              child: const Icon(Icons.description_outlined,
                  color: Color(0xff5d60dc), size: 18)),
          const SizedBox(width: 11),
          Expanded(
              child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                Text(_taskFileLabel(task),
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                    style: const TextStyle(
                        fontWeight: FontWeight.w600, fontSize: 13)),
                const SizedBox(height: 3),
                Text(
                    '${task['completed'] ?? 0}/${task['total'] ?? 0} · ${_taskStatusLabel(status)}',
                    style: const TextStyle(color: _muted, fontSize: 11)),
                if (status == 'failed' && errorMessage.isNotEmpty) ...[
                  const SizedBox(height: 4),
                  Text(errorMessage,
                      maxLines: 2,
                      overflow: TextOverflow.ellipsis,
                      style: const TextStyle(
                          color: Color(0xffc54f59), fontSize: 11)),
                ],
                if (!finished) ...[
                  const SizedBox(height: 8),
                  ClipRRect(
                      borderRadius: BorderRadius.circular(10),
                      child: LinearProgressIndicator(
                          value: progress,
                          minHeight: 5,
                          color: _brand,
                          backgroundColor: const Color(0xffeceef6)))
                ]
              ])),
          const SizedBox(width: 7),
          if (finished && status == 'completed')
            IconButton(
                tooltip: '打开生成目录',
                constraints: const BoxConstraints(minWidth: 34, minHeight: 34),
                padding: EdgeInsets.zero,
                visualDensity: VisualDensity.compact,
                onPressed: () => _openTaskDirectory(task),
                icon: const Icon(Icons.folder_open_outlined,
                    color: _brand, size: 19))
          else if (!finished)
            IconButton(
                tooltip: '取消任务',
                constraints: const BoxConstraints(minWidth: 34, minHeight: 34),
                padding: EdgeInsets.zero,
                visualDensity: VisualDensity.compact,
                onPressed: () => _cancelTask(task['id'] as String),
                icon: const Icon(Icons.stop_circle_outlined,
                    color: Color(0xff9aa1b4), size: 19))
        ]));
  }

  Widget _buildSettingsPage(double width, double height) {
    final oneColumn = width < 860;
    final translationCard = _settingsCard(
        '翻译服务',
        '配置默认引擎和引擎参数',
        [
          _field(
              '默认翻译引擎',
              DropdownButtonFormField<String>(
                  value: _selectedTranslator,
                  items: _translators.map((item) {
                    final value = item as Map<String, dynamic>;
                    return DropdownMenuItem(
                        value: value['id'] as String,
                        child: Text(value['name'] as String));
                  }).toList(),
                  onChanged: (value) =>
                      setState(() => _selectedTranslator = value))),
          if (_selectedTranslator != null)
            _providerParametersEntry(_selectedTranslator!),
        ],
        section: 'translation',
        footerActions: _selectedTranslator == null
            ? const []
            : _providerParameterActions(_selectedTranslator!),
        key: const Key('settings-translation-card'));
    final outputCard = _settingsCard(
        '输出设置',
        '控制翻译后的字幕文件格式',
        [
          _settingsFieldPair(
              _field(
                  '输出目录',
                  TextField(
                      controller: _directoryController,
                      readOnly: true,
                      onTap: _chooseOutputDirectory,
                      decoration: InputDecoration(
                          hintText: '与源文件同目录 / anto',
                          suffixIcon: IconButton(
                              tooltip: '选择输出目录',
                              onPressed: _chooseOutputDirectory,
                              icon: const Icon(Icons.folder_open_outlined))))),
              _field(
                  '导出轨道',
                  DropdownButtonFormField<int>(
                      value: _trackExport,
                      items: const [
                        DropdownMenuItem(value: 0, child: Text('双轨')),
                        DropdownMenuItem(value: 1, child: Text('主轨')),
                        DropdownMenuItem(value: 2, child: Text('副轨'))
                      ],
                      onChanged: (value) => setState(
                          () => _trackExport = value ?? _trackExport))),
              firstFlex: 3,
              secondFlex: 2),
          _switchRow('保留原始字幕文件', _keepOriginal,
              (value) => setState(() => _keepOriginal = value)),
          _switchRow('翻译完成后提示', _notifyOnComplete,
              (value) => setState(() => _notifyOnComplete = value))
        ],
        section: 'output',
        key: const Key('settings-output-card'));
    final appCard = _settingsCard(
        '应用设置',
        '调整客户端运行方式',
        [
          _settingsFieldPair(
              _field(
                  '并发数量',
                  TextField(
                      controller: _concurrencyController,
                      keyboardType: TextInputType.number,
                      decoration: const InputDecoration(hintText: '5'))),
              _field(
                  '日志级别',
                  DropdownButtonFormField<String>(
                      value: _logLevel,
                      items: const [
                        DropdownMenuItem(value: '信息', child: Text('信息')),
                        DropdownMenuItem(value: '调试', child: Text('调试')),
                        DropdownMenuItem(value: '错误', child: Text('错误'))
                      ],
                      onChanged: (value) =>
                          setState(() => _logLevel = value ?? _logLevel)))),
          _switchRow('启动时自动连接 Go Core', _autoConnect,
              (value) => setState(() => _autoConnect = value)),
          _switchRow('关闭窗口时最小化到托盘', _minimizeToTray,
              (value) => setState(() => _minimizeToTray = value)),
        ],
        section: 'app',
        footerActions: [
          OutlinedButton(onPressed: _openLogs, child: const Text('打开日志目录')),
          OutlinedButton(onPressed: _clearOldLogs, child: const Text('清理历史日志'))
        ],
        key: const Key('settings-app-card'));
    if (oneColumn) {
      return SingleChildScrollView(
          child: Column(children: [
        SizedBox(height: height, child: translationCard),
        const SizedBox(height: 16),
        SizedBox(height: height, child: outputCard),
        const SizedBox(height: 16),
        SizedBox(height: height, child: appCard),
      ]));
    }
    return Row(crossAxisAlignment: CrossAxisAlignment.stretch, children: [
      Expanded(
          child: Column(children: [
        Expanded(child: translationCard),
        const SizedBox(height: 16),
        Expanded(child: outputCard),
      ])),
      const SizedBox(width: 18),
      Expanded(child: appCard),
    ]);
  }

  Widget _settingsFieldPair(Widget first, Widget second,
      {int firstFlex = 1, int secondFlex = 1}) {
    return LayoutBuilder(builder: (context, constraints) {
      if (constraints.maxWidth < 360) {
        return Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [first, const SizedBox(height: 13), second]);
      }
      return Row(crossAxisAlignment: CrossAxisAlignment.start, children: [
        Expanded(flex: firstFlex, child: first),
        const SizedBox(width: 12),
        Expanded(flex: secondFlex, child: second),
      ]);
    });
  }

  Widget _providerForm(String id) {
    final values = _providerConfigs[id] ?? <String, dynamic>{};
    final controllers =
        _providerControllers[id] ?? <String, TextEditingController>{};
    final fields =
        values.keys.where((key) => !_isInternalProviderField(key)).toList();
    const fieldOrders = <String, List<String>>{
      'deepseek': [
        'api_key',
        'base_url',
        'model',
        'prompt',
        'max_single_text_length',
        'qps',
        'max_coroutine_num'
      ],
      'ali_cloud_mt': [
        'ak_id',
        'ak_secret',
        'region',
        'max_coroutine_num',
        'max_single_text_length',
        'qps'
      ],
      'tencent_cloud_mt': [
        'secret_id',
        'secret_key',
        'region',
        'max_coroutine_num',
        'max_single_text_length',
        'qps'
      ],
      'baidu_cloud_mt': [
        'app_id',
        'app_secret',
        'service_tier',
        'max_coroutine_num'
      ],
      'caiyunai': [
        'token',
        'max_coroutine_num',
        'max_single_text_length',
        'qps'
      ],
      'youdao_batch': [
        'app_key',
        'app_secret',
        'max_coroutine_num',
        'max_single_text_length',
        'qps'
      ],
      'deepl_pro': [
        'api_key',
        'max_coroutine_num',
        'max_single_text_length',
        'qps'
      ],
    };
    final order = fieldOrders[id];
    if (order != null) {
      fields.sort((first, second) {
        final firstIndex = order.indexOf(first);
        final secondIndex = order.indexOf(second);
        if (firstIndex >= 0 && secondIndex >= 0) {
          return firstIndex.compareTo(secondIndex);
        }
        if (firstIndex >= 0) return -1;
        if (secondIndex >= 0) return 1;
        return first.compareTo(second);
      });
    } else {
      fields.sort();
    }
    if (fields.isEmpty) {
      return const Text('该翻译引擎暂无可编辑参数。',
          style: TextStyle(color: _muted, fontSize: 12));
    }
    return Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
      if (!(_providerConfigured[id] ?? false))
        Padding(
            padding: const EdgeInsets.only(bottom: 10),
            child: Text(_providerConfigurationStatus(id),
                style:
                    const TextStyle(color: Color(0xffc28725), fontSize: 12))),
      ...fields.map((key) {
        if (id == 'ali_cloud_mt' && key == 'region') {
          return Padding(
              padding: const EdgeInsets.only(bottom: 10),
              child: _aliRegionField(controllers[key]));
        }
        if (id == 'tencent_cloud_mt' && key == 'region') {
          return Padding(
              padding: const EdgeInsets.only(bottom: 10),
              child: _tencentRegionField(controllers[key]));
        }
        if (id == 'baidu_cloud_mt' && key == 'service_tier') {
          return Padding(
              padding: const EdgeInsets.only(bottom: 10),
              child: _baiduTierField(controllers[key]));
        }
        final numeric = _isProviderNumberField(key);
        return Padding(
            padding: const EdgeInsets.only(bottom: 10),
            child: _field(
                _providerFieldLabel(key),
                TextField(
                    controller: controllers[key],
                    obscureText: _isSecretField(key),
                    minLines: key == 'prompt' ? 4 : 1,
                    maxLines: key == 'prompt' ? 8 : 1,
                    keyboardType:
                        numeric ? TextInputType.number : TextInputType.text,
                    decoration: InputDecoration(
                        hintText:
                            _isSecretField(key) && values[key] == '********'
                                ? '已保存（留空或保持遮罩以保留）'
                                : null))));
      })
    ]);
  }

  Widget _aliRegionField(TextEditingController? controller) {
    if (controller == null) return const SizedBox.shrink();
    return ValueListenableBuilder<TextEditingValue>(
        valueListenable: controller,
        builder: (context, value, _) {
          final configured = value.text.trim();
          final selected = configured.isEmpty ? 'cn-hangzhou' : configured;
          final region = aliyunRegionById(selected);
          return Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                _field(
                    '区域',
                    DropdownButtonFormField<String>(
                        value: region?.id,
                        isExpanded: true,
                        hint: Text(
                            configured.isEmpty ? '请选择地区' : '未知地区：$configured'),
                        items: aliyunRegions
                            .map((item) => DropdownMenuItem<String>(
                                value: item.id,
                                child: Text('${item.name} (${item.id})',
                                    overflow: TextOverflow.ellipsis)))
                            .toList(),
                        onChanged: (id) {
                          if (id != null) controller.text = id;
                        })),
                const SizedBox(height: 5),
                Text(
                    region == null
                        ? '当前地区不在阿里云机器翻译公网接入点列表中，请重新选择。'
                        : '公网接入地址：${region.publicEndpoint}',
                    style: const TextStyle(color: _muted, fontSize: 11)),
              ]);
        });
  }

  Widget _tencentRegionField(TextEditingController? controller) {
    if (controller == null) return const SizedBox.shrink();
    return ValueListenableBuilder<TextEditingValue>(
        valueListenable: controller,
        builder: (context, value, _) {
          final configured = value.text.trim();
          final selected = configured.isEmpty ? 'ap-guangzhou' : configured;
          return Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                _field(
                    '区域',
                    DropdownButtonFormField<String>(
                        value: tencentRegionById(selected)?.id,
                        isExpanded: true,
                        hint: Text(
                            configured.isEmpty ? '请选择地区' : '未知地区：$configured'),
                        items: tencentRegions
                            .map((item) => DropdownMenuItem<String>(
                                value: item.id,
                                child: Text('${item.name} (${item.id})',
                                    overflow: TextOverflow.ellipsis)))
                            .toList(),
                        onChanged: (id) {
                          if (id != null) controller.text = id;
                        })),
                const SizedBox(height: 5),
                const Text('公网接入地址：tmt.tencentcloudapi.com',
                    style: TextStyle(color: _muted, fontSize: 11)),
              ]);
        });
  }

  Widget _baiduTierField(TextEditingController? controller) {
    if (controller == null) return const SizedBox.shrink();
    return ValueListenableBuilder<TextEditingValue>(
        valueListenable: controller,
        builder: (context, value, _) {
          const limits = <String, (String, String)>{
            'standard': ('标准版', '单次 1000 字符 · 1 QPS'),
            'advanced': ('高级版', '单次 6000 字符 · 10 QPS'),
            'premium': ('尊享版', '单次 6000 字符 · 100 QPS'),
          };
          final tier =
              value.text.trim().isEmpty ? 'standard' : value.text.trim();
          return Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                _field(
                    '服务版本',
                    DropdownButtonFormField<String>(
                        value: limits.containsKey(tier) ? tier : null,
                        isExpanded: true,
                        hint: Text('未知版本：$tier'),
                        items: limits.entries
                            .map((entry) => DropdownMenuItem<String>(
                                value: entry.key, child: Text(entry.value.$1)))
                            .toList(),
                        onChanged: (selected) {
                          if (selected != null) controller.text = selected;
                        })),
                const SizedBox(height: 5),
                Text('${limits[tier]?.$2 ?? '请选择有效版本'}；单次请求还受 6000 字节限制',
                    style: const TextStyle(color: _muted, fontSize: 11)),
              ]);
        });
  }

  Widget _providerParametersEntry(String id) {
    final fields = (_providerConfigs[id] ?? <String, dynamic>{})
        .keys
        .where((key) => !_isInternalProviderField(key))
        .toList();
    final configured = _providerConfigured[id] ?? false;
    return Text(
      fields.isEmpty ? '该引擎暂无可编辑参数' : _providerConfigurationStatus(id),
      style: TextStyle(
        color: configured ? _muted : const Color(0xffc28725),
        fontSize: 12,
      ),
    );
  }

  String _providerConfigurationStatus(String id) {
    if (id == 'deepseek') {
      final values = _providerConfigs[id] ?? <String, dynamic>{};
      final baseUrl = (values['base_url'] as String? ?? '').trim();
      final missing = <String>[
        if (baseUrl.isEmpty) '服务地址',
        if ((values['model'] as String? ?? '').trim().isEmpty) '模型名',
        if (baseUrl == 'https://api.deepseek.com' &&
            (values['api_key'] as String? ?? '').trim().isEmpty)
          'API Key',
      ];
      if (missing.isNotEmpty) return '${missing.join('、')} 尚未保存';
    }
    if (id == 'ali_cloud_mt' ||
        id == 'tencent_cloud_mt' ||
        id == 'baidu_cloud_mt' ||
        id == 'caiyunai' ||
        id == 'deepl_pro' ||
        id == 'youdao_batch') {
      final values = _providerConfigs[id] ?? <String, dynamic>{};
      final missing = id == 'ali_cloud_mt'
          ? <String>[
              if ((values['ak_id'] as String? ?? '').trim().isEmpty)
                'Access Key',
              if ((values['ak_secret'] as String? ?? '').trim().isEmpty)
                'Access Secret',
            ]
          : id == 'tencent_cloud_mt'
              ? <String>[
                  if ((values['secret_id'] as String? ?? '').trim().isEmpty)
                    'Secret ID',
                  if ((values['secret_key'] as String? ?? '').trim().isEmpty)
                    'Secret Key',
                ]
              : id == 'baidu_cloud_mt'
                  ? <String>[
                      if ((values['app_id'] as String? ?? '').trim().isEmpty)
                        '应用 ID',
                      if ((values['app_secret'] as String? ?? '')
                          .trim()
                          .isEmpty)
                        '应用密钥',
                    ]
                  : id == 'caiyunai'
                      ? <String>[
                          if ((values['token'] as String? ?? '').trim().isEmpty)
                            'API Token',
                        ]
                      : id == 'youdao_batch'
                          ? <String>[
                              if ((values['app_key'] as String? ?? '')
                                  .trim()
                                  .isEmpty)
                                '应用 ID',
                              if ((values['app_secret'] as String? ?? '')
                                  .trim()
                                  .isEmpty)
                                '应用密钥',
                            ]
                          : <String>[
                              if ((values['api_key'] as String? ?? '')
                                  .trim()
                                  .isEmpty)
                                'API Key',
                            ];
      if (missing.isNotEmpty) return '${missing.join('、')} 尚未保存';
    }
    if (_providerConfigured[id] ?? false) return '引擎参数已配置';
    if (_reloadingCore) return '正在重新加载引擎';
    if (_translationReloadPending) return '当前任务完成后自动生效';
    return '参数已保存，等待 Core 加载';
  }

  List<Widget> _providerParameterActions(String id) {
    final fields = (_providerConfigs[id] ?? <String, dynamic>{})
        .keys
        .where((key) => !_isInternalProviderField(key))
        .toList();
    return [
      OutlinedButton(
          onPressed: _testingConnection ? null : _testConnection,
          child: Text(_testingConnection ? '连接中…' : '测试连接')),
      OutlinedButton(
        onPressed: fields.isEmpty ? null : () => _editProviderParameters(id),
        child: const Text('配置参数'),
      ),
    ];
  }

  Widget _settingsCard(String title, String description, List<Widget> children,
          {required String section,
          required Key key,
          List<Widget> footerActions = const []}) =>
      _card(
          Padding(
              padding: const EdgeInsets.all(20),
              child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(children: [
                      Expanded(
                          child: Text(title,
                              style: const TextStyle(
                                  fontWeight: FontWeight.w700, fontSize: 15))),
                      FilledButton(
                          onPressed: _reloadingCore
                              ? null
                              : () => _saveSettings(section),
                          child: const Text('保存')),
                    ]),
                    const SizedBox(height: 4),
                    Text(description,
                        style: const TextStyle(color: _muted, fontSize: 12)),
                    const SizedBox(height: 17),
                    Expanded(
                        child: SingleChildScrollView(
                            child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: children
                                    .map((child) => Padding(
                                        padding:
                                            const EdgeInsets.only(bottom: 13),
                                        child: child))
                                    .toList()))),
                    if (footerActions.isNotEmpty) ...[
                      const SizedBox(height: 12),
                      Align(
                          alignment: Alignment.centerRight,
                          child: Wrap(
                              alignment: WrapAlignment.end,
                              spacing: 10,
                              runSpacing: 8,
                              children: footerActions)),
                    ]
                  ])),
          key: key);

  Widget _switchRow(String label, bool value, ValueChanged<bool> onChanged) =>
      Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
        Text(label, style: const TextStyle(color: _ink)),
        Switch(value: value, onChanged: onChanged)
      ]);

  Widget _card(Widget child, {Key? key}) => Container(
      key: key,
      clipBehavior: Clip.antiAlias,
      decoration: BoxDecoration(
          color: Colors.white,
          border: Border.all(color: const Color(0xffe8ebf2)),
          borderRadius: BorderRadius.circular(16),
          boxShadow: const [
            BoxShadow(
                color: Color(0x19293152), blurRadius: 30, offset: Offset(0, 12))
          ]),
      child: child);
}

class LogoMark extends StatelessWidget {
  const LogoMark({super.key, this.size = 38});
  final double size;

  @override
  Widget build(BuildContext context) => Image.asset(
        'assets/LOGO.png',
        width: size,
        height: size,
        fit: BoxFit.contain,
        filterQuality: FilterQuality.high,
      );
}

class LogoLockup extends StatelessWidget {
  const LogoLockup({super.key});

  @override
  Widget build(BuildContext context) =>
      const Row(mainAxisAlignment: MainAxisAlignment.center, children: [
        LogoMark(size: 38),
        SizedBox(width: 10),
        Text('ANTO',
            style: TextStyle(
                color: Colors.white,
                fontSize: 21,
                fontWeight: FontWeight.w700,
                letterSpacing: -.5))
      ]);
}
