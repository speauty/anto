import 'dart:async';
import 'dart:io';

import 'package:anto_desktop/core_client.dart';
import 'package:anto_desktop/aliyun_regions.dart';
import 'package:anto_desktop/main.dart';
import 'package:anto_desktop/tencent_regions.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

class FakeSettingsCore extends CoreClient {
  final eventsController = StreamController<Map<String, dynamic>>.broadcast();
  Map<String, dynamic>? lastUpdate;
  String? lastConnectionPath;
  Map<String, dynamic>? lastConnectionBody;
  bool failSave = false;
  bool failConnection = false;
  int startCount = 0;
  int stopCount = 0;
  List<dynamic> tasks = <dynamic>[];

  final config = <String, dynamic>{
    'app': <String, dynamic>{
      'default_translator': 'ali_cloud_mt',
      'default_from_lang': 'en',
      'default_to_lang': 'zh',
      'output_directory': '',
      'track_export': 2,
      'keep_original': true,
      'notify_on_complete': true,
      'auto_connect': true,
      'minimize_to_tray': false,
      'log_level': '信息',
      'concurrency': 5,
    },
    'providers': <String, dynamic>{
      'ali_cloud_mt': <String, dynamic>{
        'configured': false,
        'values': <String, dynamic>{
          'ak_id': '',
          'ak_secret': '',
          'region': '',
          'max_coroutine_num': 25,
          'max_single_text_length': 3000,
          'qps': 10,
        },
      },
      'tencent_cloud_mt': <String, dynamic>{
        'configured': false,
        'values': <String, dynamic>{
          'secret_id': '',
          'secret_key': '',
          'region': 'ap-guangzhou',
          'max_coroutine_num': 5,
          'max_single_text_length': 2000,
          'qps': 5,
        },
      },
      'baidu_cloud_mt': <String, dynamic>{
        'configured': false,
        'values': <String, dynamic>{
          'app_id': '',
          'app_secret': '',
          'service_tier': 'standard',
          'max_coroutine_num': 1,
        },
      },
      'caiyunai': <String, dynamic>{
        'configured': false,
        'values': <String, dynamic>{
          'token': '',
          'max_coroutine_num': 1,
          'max_single_text_length': 1000,
          'qps': 1,
        },
      },
      'deepl_pro': <String, dynamic>{
        'configured': false,
        'values': <String, dynamic>{
          'api_key': '',
          'max_coroutine_num': 5,
          'max_single_text_length': 3000,
          'qps': 5,
        },
      },
      'deepseek': <String, dynamic>{
        'configured': false,
        'values': <String, dynamic>{
          'api_key': '',
          'base_url': 'https://api.deepseek.com',
          'model': 'deepseek-flash',
          'prompt': '请忠实翻译字幕。',
          'max_coroutine_num': 1,
          'max_single_text_length': 2000,
          'qps': 1,
        },
      },
      'youdao_batch': <String, dynamic>{
        'configured': false,
        'values': <String, dynamic>{
          'app_key': '',
          'app_secret': '',
          'max_coroutine_num': 1,
          'max_single_text_length': 5000,
          'qps': 1,
        },
      },
    },
  };

  @override
  Future<void> start({String executable = 'anto-core.exe'}) async {
    startCount++;
    if (startCount > 1) {
      final provider = (config['providers']
          as Map<String, dynamic>)['ali_cloud_mt'] as Map<String, dynamic>;
      provider['configured'] = true;
    }
  }

  @override
  Future<void> stop() async {
    stopCount++;
  }

  @override
  Future<dynamic> get(String path) async {
    switch (path) {
      case '/api/v1/config':
        return config;
      case '/api/v1/translators':
        return <Map<String, dynamic>>[
          {'id': 'ali_cloud_mt', 'name': '机器翻译-阿里云', 'valid': false},
          {'id': 'tencent_cloud_mt', 'name': '机器翻译-腾讯云', 'valid': false},
          {'id': 'baidu_cloud_mt', 'name': '机器翻译-百度云', 'valid': false},
          {'id': 'caiyunai', 'name': '彩云小译', 'valid': false},
          {'id': 'deepl_pro', 'name': 'DeepL Pro', 'valid': false},
          {'id': 'deepseek', 'name': 'AI 翻译-DeepSeek', 'valid': false},
          {'id': 'youdao_batch', 'name': '批量翻译-有道智云', 'valid': false},
        ];
      case '/api/v1/tasks':
        return tasks;
    }
    throw StateError('unexpected GET $path');
  }

  @override
  Future<dynamic> put(String path, Map<String, dynamic> body) async {
    lastUpdate = body;
    if (failSave) {
      throw const HttpException(
        '{"error":{"code":"config_write_failed","message":"配置文件没有写入权限"}}',
      );
    }
    return config;
  }

  @override
  Future<dynamic> post(String path, {Map<String, dynamic>? body}) async {
    lastConnectionPath = path;
    lastConnectionBody = body;
    if (failConnection) {
      throw const HttpException(
        '{"error":{"code":"translator_connection_failed","message":"连接公网接入点失败"}}',
      );
    }
    return <String, dynamic>{
      'endpoint': path.contains('/deepseek/')
          ? '127.0.0.1:8000'
          : 'mt.cn-hangzhou.aliyuncs.com',
      'latency_ms': 42,
    };
  }

  @override
  Stream<Map<String, dynamic>> events() => eventsController.stream;

  @override
  Future<void> dispose() async {
    await eventsController.close();
    await super.dispose();
  }
}

class StopDependentEventsCore extends FakeSettingsCore {
  final _stopped = Completer<void>();

  @override
  Stream<Map<String, dynamic>> events() {
    if (startCount > 1) return super.events();
    return StreamController<Map<String, dynamic>>(
      onCancel: () => _stopped.future,
    ).stream;
  }

  @override
  Future<void> stop() async {
    if (!_stopped.isCompleted) _stopped.complete();
    await super.stop();
  }
}

void main() {
  test('Aliyun region options use the documented public endpoints', () {
    expect(aliyunRegionById('cn-hangzhou')?.publicEndpoint,
        'mt.cn-hangzhou.aliyuncs.com');
    expect(aliyunRegionById('ap-southeast-1')?.publicEndpoint,
        'mt.ap-southeast-1.aliyuncs.com');
    expect(aliyunRegionById('me-central-1')?.publicEndpoint,
        'mt.me-central-1.aliyuncs.com');
    expect(aliyunRegionById('me-east-1')?.publicEndpoint,
        'alimt.me-east-1.aliyuncs.com');
    expect(aliyunRegionById('ap-southeast-2'), isNull);
  });

  test('Tencent regions include the official default region', () {
    expect(tencentRegionById('ap-guangzhou')?.name, '广州');
    expect(tencentRegionById('invalid'), isNull);
  });

  testWidgets('Tencent settings expose credentials and region dropdown',
      (tester) async {
    final core = FakeSettingsCore();
    await tester.pumpWidget(AntoApp(core: core));
    await tester.pump();
    await tester.tap(find.text('我的设置'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('机器翻译-阿里云'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('机器翻译-腾讯云').last);
    await tester.pumpAndSettle();
    await tester.tap(find.text('配置参数'));
    await tester.pumpAndSettle();

    expect(find.text('Secret ID、Secret Key 尚未保存'), findsWidgets);
    expect(find.text('公网接入地址：tmt.tencentcloudapi.com'), findsOneWidget);
    expect(find.text('广州 (ap-guangzhou)'), findsOneWidget);
    await tester.tap(find.text('完成'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('测试连接'));
    await tester.pumpAndSettle();
    expect(core.lastConnectionPath,
        '/api/v1/translators/tencent_cloud_mt/test-connection');
    expect(
        core.lastConnectionBody, <String, dynamic>{'region': 'ap-guangzhou'});
    await tester.pumpWidget(const SizedBox.shrink());
  });

  testWidgets('Baidu version selects request limits and uses regionless test',
      (tester) async {
    final core = FakeSettingsCore();
    await tester.pumpWidget(AntoApp(core: core));
    await tester.pump();
    await tester.tap(find.text('我的设置'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('机器翻译-阿里云'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('机器翻译-百度云').last);
    await tester.pumpAndSettle();
    await tester.tap(find.text('配置参数'));
    await tester.pumpAndSettle();

    expect(find.text('应用 ID、应用密钥 尚未保存'), findsWidgets);
    expect(find.text('单次 1000 字符 · 1 QPS；单次请求还受 6000 字节限制'), findsOneWidget);
    await tester.tap(find.text('标准版'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('高级版').last);
    await tester.pumpAndSettle();
    expect(find.text('单次 6000 字符 · 10 QPS；单次请求还受 6000 字节限制'), findsOneWidget);
    await tester.tap(find.text('完成'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('测试连接'));
    await tester.pumpAndSettle();
    expect(core.lastConnectionPath,
        '/api/v1/translators/baidu_cloud_mt/test-connection');
    expect(core.lastConnectionBody, isEmpty);
    await tester.tap(find.text('知道了'));
    await tester.pumpAndSettle();
    await tester.tap(find.descendant(
        of: find.byKey(const Key('settings-translation-card')),
        matching: find.text('保存')));
    await tester.pumpAndSettle();
    await tester
        .runAsync(() => Future<void>.delayed(const Duration(milliseconds: 20)));
    await tester.pumpAndSettle();
    expect(core.lastUpdate?['providers'], <String, dynamic>{
      'baidu_cloud_mt': <String, dynamic>{'service_tier': 'advanced'},
    });
    await tester.pumpWidget(const SizedBox.shrink());
  });

  testWidgets('Caiyun token is configured and connection test has no region',
      (tester) async {
    final core = FakeSettingsCore();
    await tester.pumpWidget(AntoApp(core: core));
    await tester.pump();
    await tester.tap(find.text('我的设置'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('机器翻译-阿里云'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('彩云小译').last);
    await tester.pumpAndSettle();
    await tester.tap(find.text('配置参数'));
    await tester.pumpAndSettle();
    expect(find.text('API Token 尚未保存'), findsWidgets);
    expect(find.text('访问令牌（Token）'), findsOneWidget);
    await tester.enterText(
        find
            .descendant(
                of: find.byType(AlertDialog), matching: find.byType(TextField))
            .first,
        'test-token');
    await tester.tap(find.text('完成'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('测试连接'));
    await tester.pumpAndSettle();
    expect(core.lastConnectionPath,
        '/api/v1/translators/caiyunai/test-connection');
    expect(core.lastConnectionBody, isEmpty);
    await tester.tap(find.text('知道了'));
    await tester.pumpAndSettle();
    await tester.tap(find.descendant(
        of: find.byKey(const Key('settings-translation-card')),
        matching: find.text('保存')));
    await tester.pumpAndSettle();
    expect(core.lastUpdate?['providers'], <String, dynamic>{
      'caiyunai': <String, dynamic>{'token': 'test-token'},
    });
    await tester.pumpWidget(const SizedBox.shrink());
  });

  testWidgets('DeepL Pro saves its API key and tests the paid endpoint',
      (tester) async {
    final core = FakeSettingsCore();
    await tester.pumpWidget(AntoApp(core: core));
    await tester.pump();
    await tester.tap(find.text('我的设置'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('机器翻译-阿里云'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('DeepL Pro').last);
    await tester.pumpAndSettle();
    await tester.tap(find.text('配置参数'));
    await tester.pumpAndSettle();
    expect(find.text('API Key 尚未保存'), findsWidgets);
    expect(find.text('API Key'), findsWidgets);
    await tester.enterText(
        find
            .descendant(
                of: find.byType(AlertDialog), matching: find.byType(TextField))
            .first,
        'test-api-key');
    await tester.tap(find.text('完成'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('测试连接'));
    await tester.pumpAndSettle();
    expect(core.lastConnectionPath,
        '/api/v1/translators/deepl_pro/test-connection');
    expect(core.lastConnectionBody, isEmpty);
    await tester.tap(find.text('知道了'));
    await tester.pumpAndSettle();
    await tester.tap(find.descendant(
        of: find.byKey(const Key('settings-translation-card')),
        matching: find.text('保存')));
    await tester.pumpAndSettle();
    expect(core.lastUpdate?['providers'], <String, dynamic>{
      'deepl_pro': <String, dynamic>{'api_key': 'test-api-key'},
    });
    await tester.pumpWidget(const SizedBox.shrink());
  });

  testWidgets('DeepSeek saves endpoint, model and custom prompt',
      (tester) async {
    final core = FakeSettingsCore();
    await tester.pumpWidget(AntoApp(core: core));
    await tester.pump();
    await tester.tap(find.text('我的设置'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('机器翻译-阿里云'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('AI 翻译-DeepSeek').last);
    await tester.pumpAndSettle();
    await tester.tap(find.text('配置参数'));
    await tester.pumpAndSettle();
    final fields = find.descendant(
        of: find.byType(AlertDialog), matching: find.byType(TextField));
    expect(find.text('API Key 尚未保存'), findsWidgets);
    await tester.enterText(fields.at(0), 'test-api-key');
    await tester.enterText(fields.at(1), 'http://127.0.0.1:8000/v1');
    await tester.enterText(fields.at(2), 'local-model');
    await tester.enterText(fields.at(3), '自然、简洁');
    await tester.tap(find.text('完成'));
    await tester.pumpAndSettle();
    await tester.tap(find.descendant(
        of: find.byKey(const Key('settings-translation-card')),
        matching: find.text('保存')));
    await tester.pumpAndSettle();
    expect(core.lastUpdate?['providers'], <String, dynamic>{
      'deepseek': <String, dynamic>{
        'api_key': 'test-api-key',
        'base_url': 'http://127.0.0.1:8000/v1',
        'model': 'local-model',
        'prompt': '自然、简洁',
      },
    });
    await tester.tap(find.text('测试连接'));
    await tester.pumpAndSettle();
    expect(core.lastConnectionPath,
        '/api/v1/translators/deepseek/test-connection');
    expect(core.lastConnectionBody, isEmpty);
    expect(find.textContaining('127.0.0.1:8000'), findsOneWidget);
    expect(find.textContaining('127.0.0.1:8000:443'), findsNothing);
    await tester.tap(find.text('知道了'));
    await tester.pumpAndSettle();
    await tester.pumpWidget(const SizedBox.shrink());
  });

  testWidgets('Youdao batch saves app credentials and tests its endpoint',
      (tester) async {
    final core = FakeSettingsCore();
    await tester.pumpWidget(AntoApp(core: core));
    await tester.pump();
    await tester.tap(find.text('我的设置'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('机器翻译-阿里云'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('批量翻译-有道智云').last);
    await tester.pumpAndSettle();
    await tester.tap(find.text('配置参数'));
    await tester.pumpAndSettle();
    expect(find.text('应用 ID、应用密钥 尚未保存'), findsWidgets);
    final fields = find.descendant(
        of: find.byType(AlertDialog), matching: find.byType(TextField));
    await tester.enterText(fields.at(0), 'test-app-id');
    await tester.enterText(fields.at(1), 'test-secret');
    await tester.tap(find.text('完成'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('测试连接'));
    await tester.pumpAndSettle();
    expect(core.lastConnectionPath,
        '/api/v1/translators/youdao_batch/test-connection');
    expect(core.lastConnectionBody, isEmpty);
    await tester.tap(find.text('知道了'));
    await tester.pumpAndSettle();
    await tester.tap(find.descendant(
        of: find.byKey(const Key('settings-translation-card')),
        matching: find.text('保存')));
    await tester.pumpAndSettle();
    expect(core.lastUpdate?['providers'], <String, dynamic>{
      'youdao_batch': <String, dynamic>{
        'app_key': 'test-app-id',
        'app_secret': 'test-secret',
      },
    });
    await tester.pumpWidget(const SizedBox.shrink());
  });

  testWidgets('Aliyun region is a dropdown with its public endpoint',
      (tester) async {
    final core = FakeSettingsCore();
    await tester.pumpWidget(AntoApp(core: core));
    await tester.pump();
    await tester.tap(find.byIcon(Icons.settings_outlined));
    await tester.pumpAndSettle();
    await tester.tap(find.text('配置参数'));
    await tester.pumpAndSettle();

    expect(find.text('Access Key、Access Secret 尚未保存'), findsWidgets);
    expect(tester.getTopLeft(find.text('区域')).dy,
        lessThan(tester.getTopLeft(find.text('协程数量')).dy));
    expect(find.text('公网接入地址：mt.cn-hangzhou.aliyuncs.com'), findsOneWidget);
    expect(find.text('华东1（杭州） (cn-hangzhou)'), findsOneWidget);
    await tester.tap(find.text('华东1（杭州） (cn-hangzhou)'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('华东2（上海） (cn-shanghai)').last);
    await tester.pumpAndSettle();
    expect(find.text('公网接入地址：mt.aliyuncs.com'), findsOneWidget);

    await tester.tap(find.text('完成'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('测试连接'));
    await tester.pumpAndSettle();
    expect(core.lastConnectionBody, <String, dynamic>{'region': 'cn-shanghai'});
    await tester.tap(find.text('知道了'));
    await tester.pumpAndSettle();
    await tester.tap(find.descendant(
        of: find.byKey(const Key('settings-translation-card')),
        matching: find.text('保存')));
    await tester.pumpAndSettle();
    expect(core.lastUpdate?['providers'], <String, dynamic>{
      'ali_cloud_mt': <String, dynamic>{'region': 'cn-shanghai'},
    });
    await tester.pumpWidget(const SizedBox.shrink());
  });

  testWidgets('connection and parameter buttons sit together at card bottom',
      (tester) async {
    final core = FakeSettingsCore();
    await tester.pumpWidget(AntoApp(core: core));
    await tester.pump();
    await tester.tap(find.text('我的设置'));
    await tester.pumpAndSettle();

    expect(find.text('配置默认引擎和引擎参数'), findsOneWidget);
    expect(find.text('控制翻译后的字幕文件格式'), findsOneWidget);
    expect(find.text('调整客户端运行方式'), findsOneWidget);

    final connection = tester.getRect(find.ancestor(
        of: find.text('测试连接'), matching: find.byType(OutlinedButton)));
    final parameters = tester.getRect(find.ancestor(
        of: find.text('配置参数'),
        matching:
            find.byWidgetPredicate((widget) => widget is OutlinedButton)));
    final card =
        tester.getRect(find.byKey(const Key('settings-translation-card')));
    expect(connection.center.dy, parameters.center.dy);
    expect(parameters.left - connection.right, 10);
    expect(card.right - parameters.right, inInclusiveRange(20, 22));
    expect(card.bottom - parameters.bottom, inInclusiveRange(20, 22));

    await tester.tap(find.text('测试连接'));
    await tester.pumpAndSettle();
    expect(core.lastConnectionPath,
        '/api/v1/translators/ali_cloud_mt/test-connection');
    expect(core.lastConnectionBody, <String, dynamic>{'region': 'cn-hangzhou'});
    expect(find.text('引擎网络可达'), findsOneWidget);
    expect(
        find.textContaining('mt.cn-hangzhou.aliyuncs.com:443'), findsOneWidget);
    await tester.pumpWidget(const SizedBox.shrink());
  });

  testWidgets('saved credentials show that Core has not loaded them yet',
      (tester) async {
    final core = FakeSettingsCore();
    final provider = (core.config['providers']
        as Map<String, dynamic>)['ali_cloud_mt'] as Map<String, dynamic>;
    final values = provider['values'] as Map<String, dynamic>;
    values['ak_id'] = '********';
    values['ak_secret'] = '********';
    await tester.pumpWidget(AntoApp(core: core));
    await tester.pump();
    await tester.tap(find.text('我的设置'));
    await tester.pumpAndSettle();

    expect(find.text('参数已保存，等待 Core 加载'), findsOneWidget);
    await tester.pumpWidget(const SizedBox.shrink());
  });

  testWidgets('settings cards and related fields use compact paired layout',
      (tester) async {
    await tester.binding.setSurfaceSize(const Size(1300, 900));
    addTearDown(() => tester.binding.setSurfaceSize(null));
    final core = FakeSettingsCore();
    await tester.pumpWidget(AntoApp(core: core));
    await tester.pump();
    await tester.tap(find.text('我的设置'));
    await tester.pumpAndSettle();

    expect(tester.getTopLeft(find.text('输出目录')).dy,
        tester.getTopLeft(find.text('导出轨道')).dy);
    expect(tester.getTopLeft(find.text('并发数量')).dy,
        tester.getTopLeft(find.text('日志级别')).dy);
    final outputTop =
        tester.getTopLeft(find.byKey(const Key('settings-output-card'))).dy;
    final translationBottom = tester
        .getBottomRight(find.byKey(const Key('settings-translation-card')))
        .dy;
    expect(outputTop - translationBottom, 16);
    expect(tester.takeException(), isNull);
    await tester.pumpWidget(const SizedBox.shrink());
  });

  testWidgets('settings cards fill the height and save their own fields',
      (tester) async {
    await tester.binding.setSurfaceSize(const Size(1300, 900));
    addTearDown(() => tester.binding.setSurfaceSize(null));
    final core = FakeSettingsCore();
    await tester.pumpWidget(AntoApp(core: core));
    await tester.pump();
    await tester.tap(find.text('我的设置'));
    await tester.pumpAndSettle();

    final translation =
        tester.getRect(find.byKey(const Key('settings-translation-card')));
    final app = tester.getRect(find.byKey(const Key('settings-app-card')));
    final output =
        tester.getRect(find.byKey(const Key('settings-output-card')));
    expect(translation.top, 30);
    expect(app.top, 30);
    expect(app.bottom, 856);
    expect(output.bottom, 856);
    expect(output.left, translation.left);
    expect(app.left, greaterThan(output.right));
    final openLogs = tester.getRect(find.ancestor(
        of: find.text('打开日志目录'),
        matching:
            find.byWidgetPredicate((widget) => widget is OutlinedButton)));
    final clearLogs = tester.getRect(find.ancestor(
        of: find.text('清理历史日志'),
        matching:
            find.byWidgetPredicate((widget) => widget is OutlinedButton)));
    expect(clearLogs.left - openLogs.right, 10);
    expect(app.right - clearLogs.right, inInclusiveRange(20, 22));
    expect(app.bottom - clearLogs.bottom, inInclusiveRange(20, 22));
    expect(find.text('保存设置'), findsNothing);
    expect(find.text('保存'), findsNWidgets(3));
    for (final key in const [
      Key('settings-translation-card'),
      Key('settings-app-card'),
      Key('settings-output-card'),
    ]) {
      expect(
          find.descendant(
              of: find.byKey(key), matching: find.byType(Scrollable)),
          findsWidgets);
    }

    await tester.enterText(
        find.descendant(
            of: find.byKey(const Key('settings-app-card')),
            matching: find.byType(TextField)),
        '9');
    await tester.tap(find.descendant(
        of: find.byKey(const Key('settings-output-card')),
        matching: find.text('保存')));
    await tester.pumpAndSettle();
    expect(core.lastUpdate?['app'], containsPair('track_export', 2));
    expect(core.lastUpdate?['app'], isNot(contains('concurrency')));
    await tester.tap(find.text('知道了'));
    await tester.pumpAndSettle();

    final saveApp = find.descendant(
        of: find.byKey(const Key('settings-app-card')),
        matching: find.text('保存'));
    await tester.ensureVisible(saveApp);
    await tester.pumpAndSettle();
    await tester.tap(saveApp);
    await tester.pumpAndSettle();
    expect(core.lastUpdate?['app'], containsPair('concurrency', 9));
    expect(core.lastUpdate?['app'], isNot(contains('track_export')));
    await tester.pumpWidget(const SizedBox.shrink());
  });

  testWidgets('settings fields stack without overflow in a narrow window',
      (tester) async {
    await tester.binding.setSurfaceSize(const Size(470, 900));
    addTearDown(() => tester.binding.setSurfaceSize(null));
    final core = FakeSettingsCore();
    await tester.pumpWidget(AntoApp(core: core));
    await tester.pump();
    await tester.tap(find.byIcon(Icons.settings_outlined));
    await tester.pumpAndSettle();

    final translation =
        tester.getRect(find.byKey(const Key('settings-translation-card')));
    final output =
        tester.getRect(find.byKey(const Key('settings-output-card')));
    final app = tester.getRect(find.byKey(const Key('settings-app-card')));
    expect(translation.top, lessThan(output.top));
    expect(output.top, lessThan(app.top));

    expect(tester.getTopLeft(find.text('并发数量')).dy,
        lessThan(tester.getTopLeft(find.text('日志级别')).dy));
    expect(tester.getTopLeft(find.text('输出目录')).dy,
        lessThan(tester.getTopLeft(find.text('导出轨道')).dy));
    expect(tester.takeException(), isNull);
    await tester.pumpWidget(const SizedBox.shrink());
  });

  testWidgets('connection failure displays provider network error',
      (tester) async {
    final core = FakeSettingsCore()..failConnection = true;
    await tester.pumpWidget(AntoApp(core: core));
    await tester.pump();
    await tester.tap(find.text('我的设置'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('测试连接'));
    await tester.pumpAndSettle();
    expect(find.text('引擎连接失败'), findsOneWidget);
    expect(find.text('连接公网接入点失败'), findsOneWidget);
    await tester.pumpWidget(const SizedBox.shrink());
  });

  testWidgets('engine parameters open in a dialog and cancel restores values',
      (tester) async {
    final core = FakeSettingsCore();
    await tester.pumpWidget(AntoApp(core: core));
    await tester.pump();
    await tester.tap(find.text('我的设置'));
    await tester.pumpAndSettle();

    expect(find.text('配置参数'), findsOneWidget);
    expect(find.text('QPS'), findsNothing);
    await tester.tap(find.text('配置参数'));
    await tester.pumpAndSettle();
    expect(find.text('机器翻译-阿里云 参数'), findsOneWidget);
    expect(find.text('QPS'), findsOneWidget);
    await tester.enterText(find.byType(TextField).last, '20');
    await tester.tap(find.text('取消'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('配置参数'));
    await tester.pumpAndSettle();
    expect(find.text('10'), findsOneWidget);
    await tester.tap(find.text('取消'));
    await tester.pumpAndSettle();
    await tester.pumpWidget(const SizedBox.shrink());
  });

  testWidgets('saving ignores untouched engines and shows Core error details',
      (tester) async {
    final core = FakeSettingsCore()..failSave = true;
    await tester.pumpWidget(AntoApp(core: core));
    await tester.pump();
    await tester.tap(find.text('我的设置'));
    await tester.pumpAndSettle();
    final saveApp = find.descendant(
        of: find.byKey(const Key('settings-app-card')),
        matching: find.text('保存'));
    await tester.ensureVisible(saveApp);
    await tester.pumpAndSettle();
    await tester.tap(saveApp);
    await tester.pumpAndSettle();
    await tester.pump(const Duration(milliseconds: 100));
    await tester.pumpAndSettle();

    expect(core.lastUpdate?['providers'], isEmpty);
    expect(find.text('设置保存失败'), findsWidgets);
    expect(find.text('配置文件没有写入权限'), findsOneWidget);
    await tester.pumpWidget(const SizedBox.shrink());
  });

  testWidgets('saving sends only parameters changed in the dialog',
      (tester) async {
    final core = FakeSettingsCore();
    await tester.pumpWidget(AntoApp(core: core));
    await tester.pump();
    await tester.tap(find.text('我的设置'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('配置参数'));
    await tester.pumpAndSettle();
    await tester.enterText(find.byType(TextField).last, '20');
    await tester.tap(find.text('完成'));
    await tester.pumpAndSettle();
    await tester.tap(find.descendant(
        of: find.byKey(const Key('settings-translation-card')),
        matching: find.text('保存')));
    await tester.pumpAndSettle();
    await tester
        .runAsync(() => Future<void>.delayed(const Duration(milliseconds: 20)));
    await tester.pumpAndSettle();
    await tester.pump(const Duration(milliseconds: 100));
    await tester.pumpAndSettle();

    expect(core.lastUpdate?['providers'], <String, dynamic>{
      'ali_cloud_mt': <String, dynamic>{'qps': 20},
    });
    expect(core.startCount, 2);
    expect(core.stopCount, 1);
    expect(find.text('翻译服务 已保存并生效。'), findsOneWidget);
    expect(
      find.text('设置已保存'),
      findsWidgets,
      reason: tester
          .widgetList<Text>(find.byType(Text))
          .map((widget) => widget.data)
          .whereType<String>()
          .join(' | '),
    );
    await tester.pumpWidget(const SizedBox.shrink());
  });

  testWidgets('translation reload keeps unsaved app settings', (tester) async {
    final core = FakeSettingsCore();
    await tester.pumpWidget(AntoApp(core: core));
    await tester.pump();
    await tester.tap(find.text('我的设置'));
    await tester.pumpAndSettle();
    final concurrency = find.descendant(
        of: find.byKey(const Key('settings-app-card')),
        matching: find.byType(TextField));
    await tester.enterText(concurrency, '9');
    await tester.tap(find.descendant(
        of: find.byKey(const Key('settings-translation-card')),
        matching: find.text('保存')));
    await tester.pumpAndSettle();

    await tester
        .runAsync(() => Future<void>.delayed(const Duration(milliseconds: 20)));
    await tester.pumpAndSettle();

    expect(core.startCount, 2);
    expect(core.stopCount, 1);
    expect(tester.widget<TextField>(concurrency).controller?.text, '9');
    await tester.pumpWidget(const SizedBox.shrink());
  });

  testWidgets('translation reload closes Core before waiting for events',
      (tester) async {
    final core = StopDependentEventsCore();
    await tester.pumpWidget(AntoApp(core: core));
    await tester.pump();
    await tester.tap(find.text('我的设置'));
    await tester.pumpAndSettle();
    await tester.tap(find.descendant(
        of: find.byKey(const Key('settings-translation-card')),
        matching: find.text('保存')));
    await tester.pumpAndSettle();
    await tester
        .runAsync(() => Future<void>.delayed(const Duration(milliseconds: 20)));
    await tester.pumpAndSettle();

    expect(core.stopCount, 1);
    expect(core.startCount, 2);
    expect(find.text('正在重新加载引擎'), findsNothing);
    expect(find.text('翻译服务 已保存并生效。'), findsOneWidget);
    await tester.pumpWidget(const SizedBox.shrink());
  });

  testWidgets('translation reload waits for active tasks', (tester) async {
    final core = FakeSettingsCore()
      ..tasks = <dynamic>[
        <String, dynamic>{'id': 'task-1', 'status': 'translating', 'total': 1}
      ];
    (core.config['app'] as Map<String, dynamic>)['notify_on_complete'] = false;
    await tester.pumpWidget(AntoApp(core: core));
    await tester.pump();
    await tester.tap(find.text('我的设置'));
    await tester.pumpAndSettle();
    await tester.tap(find.descendant(
        of: find.byKey(const Key('settings-translation-card')),
        matching: find.text('保存')));
    await tester.pumpAndSettle();

    expect(core.startCount, 1);
    expect(core.stopCount, 0);
    expect(find.text('翻译服务 已保存，当前任务完成后会自动生效。'), findsOneWidget);
    await tester.tap(find.text('知道了'));
    await tester.pumpAndSettle();
    core.tasks = <dynamic>[
      <String, dynamic>{'id': 'task-1', 'status': 'completed', 'total': 1}
    ];
    core.eventsController.add(<String, dynamic>{
      'type': 'task.completed',
      'task': <String, dynamic>{
        'id': 'task-1',
        'status': 'completed',
        'total': 1
      },
    });
    await tester.pumpAndSettle();
    await tester
        .runAsync(() => Future<void>.delayed(const Duration(milliseconds: 20)));
    await tester.pumpAndSettle();

    expect(core.startCount, 2);
    expect(core.stopCount, 1);
    await tester.pumpWidget(const SizedBox.shrink());
  });
}
