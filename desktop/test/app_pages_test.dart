import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:anto_desktop/core_client.dart';
import 'package:anto_desktop/main.dart';
import 'package:anto_desktop/windows_picker.dart';

class _TestPicker extends WindowsPicker {
  const _TestPicker();

  @override
  Future<List<String>> pickFiles() async => [
        r'C:\subtitles\sample.srt',
        r'C:\subtitles\second.srt',
      ];

  @override
  Future<String?> pickDirectory() async => r'C:\subtitles';
}

class _TaskCore extends CoreClient {
  final _events = StreamController<Map<String, dynamic>>.broadcast();
  final tasks = <Map<String, dynamic>>[];
  Map<String, dynamic>? lastTaskRequest;
  bool failCreate = false;

  @override
  Future<void> start({String executable = 'anto-core.exe'}) async {}

  @override
  Future<dynamic> get(String path) async {
    switch (path) {
      case '/api/v1/translators':
        return <Map<String, dynamic>>[
          {
            'id': 'ali_cloud_mt',
            'name': '机器翻译-阿里云',
            'valid': true,
            'languages': [
              {'id': 'en', 'name': '英语'},
              {'id': 'zh', 'name': '中文'},
            ],
          },
          {
            'id': 'tencent_cloud_mt',
            'name': '机器翻译-腾讯云',
            'valid': true,
            'languages': [
              {'id': 'en', 'name': '英语'},
              {'id': 'zh', 'name': '中文'},
              {'id': 'jp', 'name': '日语'},
            ],
          },
          {
            'id': 'baidu_cloud_mt',
            'name': '机器翻译-百度云',
            'valid': true,
            'languages': [
              {'id': 'en', 'name': '英语'},
              {'id': 'zh', 'name': '中文'},
              {'id': 'jp', 'name': '日语'},
            ],
          },
          {
            'id': 'caiyunai',
            'name': '彩云小译',
            'valid': true,
            'languages': [
              {'id': 'zh', 'name': '中文'},
              {'id': 'zh-Hant', 'name': '繁体中文'},
              {'id': 'en', 'name': '英语'},
              {'id': 'ja', 'name': '日语'},
              {'id': 'de', 'name': '德语'},
            ],
          },
          {
            'id': 'deepl_pro',
            'name': 'DeepL Pro',
            'valid': true,
            'languages': [
              {'id': 'zh', 'name': '中文'},
              {'id': 'en', 'name': '英语'},
              {'id': 'de', 'name': '德语'},
            ],
          },
          {
            'id': 'youdao_batch',
            'name': '批量翻译-有道智云',
            'valid': true,
            'languages': [
              {'id': 'zh', 'name': '中文'},
              {'id': 'en', 'name': '英语'},
              {'id': 'ja', 'name': '日语'},
            ],
          },
        ];
      case '/api/v1/config':
        return <String, dynamic>{
          'app': <String, dynamic>{
            'default_translator': 'ali_cloud_mt',
            'default_from_lang': 'en',
            'default_to_lang': 'zh',
          },
          'providers': <String, dynamic>{},
        };
      case '/api/v1/tasks':
        return tasks;
    }
    throw StateError('unexpected GET $path');
  }

  @override
  Future<dynamic> post(String path, {Map<String, dynamic>? body}) async {
    if (path != '/api/v1/tasks') throw StateError('unexpected POST $path');
    lastTaskRequest = body;
    if (failCreate) throw StateError('任务创建失败');
    final task = <String, dynamic>{
      'id': 'task-1',
      'status': 'queued',
      'total': 2,
      'completed': 0,
      'created_at': DateTime.now().toIso8601String(),
    };
    tasks.add(task);
    return task;
  }

  @override
  Stream<Map<String, dynamic>> events() => _events.stream;

  @override
  Future<void> dispose() async {
    await _events.close();
    await super.dispose();
  }
}

void main() {
  testWidgets('renders the workspace, queue, and settings pages',
      (tester) async {
    await tester.pumpWidget(const AntoApp(autoConnect: false));
    await tester.pumpAndSettle();

    expect(find.text('ANTO'), findsOneWidget);
    final logo = tester.getRect(find.byType(LogoMark));
    final name = tester.getRect(find.text('ANTO'));
    expect((logo.left + name.right) / 2, closeTo(112, 1));

    expect(find.text('新建翻译任务'), findsOneWidget);
    expect(find.text('源语言'), findsOneWidget);
    expect(find.text('目标语言'), findsOneWidget);

    await tester.tap(find.text('任务队列'));
    await tester.pumpAndSettle();
    expect(find.text('任务记录'), findsOneWidget);

    await tester.tap(find.text('我的设置'));
    await tester.pumpAndSettle();
    expect(find.text('输出目录'), findsOneWidget);
    expect(find.byTooltip('选择输出目录'), findsOneWidget);
  });

  testWidgets('new task can choose Tencent Cloud translation', (tester) async {
    final core = _TaskCore();
    await tester.pumpWidget(AntoApp(core: core, picker: const _TestPicker()));
    await tester.pumpAndSettle();

    await tester.tap(find.text('机器翻译-阿里云'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('机器翻译-腾讯云').last);
    await tester.pumpAndSettle();
    await tester.tap(find.text('选择文件'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('开始翻译'));
    await tester.pumpAndSettle();

    expect(core.lastTaskRequest?['translator_id'], 'tencent_cloud_mt');
    await tester.pumpWidget(const SizedBox.shrink());
  });

  testWidgets('new task can choose Baidu Cloud translation', (tester) async {
    final core = _TaskCore();
    await tester.pumpWidget(AntoApp(core: core, picker: const _TestPicker()));
    await tester.pumpAndSettle();

    await tester.tap(find.text('机器翻译-阿里云'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('机器翻译-百度云').last);
    await tester.pumpAndSettle();
    await tester.tap(find.text('选择文件'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('开始翻译'));
    await tester.pumpAndSettle();

    expect(core.lastTaskRequest?['translator_id'], 'baidu_cloud_mt');
    await tester.pumpWidget(const SizedBox.shrink());
  });

  testWidgets('Caiyun direction choices and attribution are visible',
      (tester) async {
    final core = _TaskCore();
    await tester.pumpWidget(AntoApp(core: core, picker: const _TestPicker()));
    await tester.pumpAndSettle();
    await tester.tap(find.text('机器翻译-阿里云'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('彩云小译').last);
    await tester.pumpAndSettle();
    expect(find.text('数据来自彩云小译'), findsOneWidget);
    await tester.tap(find.text('选择文件'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('开始翻译'));
    await tester.pumpAndSettle();
    expect(core.lastTaskRequest?['translator_id'], 'caiyunai');
    await tester.pumpWidget(const SizedBox.shrink());
  });

  testWidgets('new task can choose DeepL Pro translation', (tester) async {
    final core = _TaskCore();
    await tester.pumpWidget(AntoApp(core: core, picker: const _TestPicker()));
    await tester.pumpAndSettle();
    await tester.tap(find.text('机器翻译-阿里云'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('DeepL Pro').last);
    await tester.pumpAndSettle();
    await tester.tap(find.text('选择文件'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('开始翻译'));
    await tester.pumpAndSettle();
    expect(core.lastTaskRequest?['translator_id'], 'deepl_pro');
    await tester.pumpWidget(const SizedBox.shrink());
  });

  testWidgets('new task can choose Youdao batch translation', (tester) async {
    final core = _TaskCore();
    await tester.pumpWidget(AntoApp(core: core, picker: const _TestPicker()));
    await tester.pumpAndSettle();
    await tester.tap(find.text('机器翻译-阿里云'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('批量翻译-有道智云').last);
    await tester.pumpAndSettle();
    await tester.tap(find.text('选择文件'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('开始翻译'));
    await tester.pumpAndSettle();
    expect(core.lastTaskRequest?['translator_id'], 'youdao_batch');
    await tester.pumpWidget(const SizedBox.shrink());
  });

  testWidgets('about page shows introduction, donation, and engine guidance',
      (tester) async {
    await tester.binding.setSurfaceSize(const Size(1264, 789));
    addTearDown(() => tester.binding.setSurfaceSize(null));
    await tester.pumpWidget(const AntoApp(autoConnect: false));
    await tester.pumpAndSettle();
    await tester.tap(find.text('关于应用'));
    await tester.pumpAndSettle();

    final intro = tester.getRect(find.byKey(const Key('about-intro-card')));
    final author = tester.getRect(find.byKey(const Key('about-author-card')));
    final donation =
        tester.getRect(find.byKey(const Key('about-donation-card')));
    final engine = tester.getRect(find.byKey(const Key('about-engine-card')));
    expect(intro.top, 30);
    expect(intro.right, donation.right,
        reason: 'intro=$intro author=$author donation=$donation');
    expect(author.top, donation.top);
    expect(author.bottom, donation.bottom);
    expect(engine.top, greaterThan(donation.bottom));
    expect(find.text('Anto-你的字幕翻译小助手'), findsOneWidget);
    expect(find.text('speauty'), findsOneWidget);
    expect(find.text('微信赞赏'), findsOneWidget);
    expect(find.text('支付宝赞赏'), findsOneWidget);
    expect(find.text('翻译引擎'), findsOneWidget);
    expect(find.text('机器翻译-阿里云'), findsOneWidget);
    expect(find.text('机器翻译-腾讯云'), findsOneWidget);
    expect(find.text('机器翻译-百度云'), findsOneWidget);
    expect(find.text('彩云小译'), findsOneWidget);
    expect(find.text('DeepL Pro'), findsOneWidget);
    expect(find.text('AI 翻译-DeepSeek'), findsOneWidget);
    expect(find.text('批量翻译-有道智云'), findsOneWidget);
    final aliyunTitle = tester.getRect(find.text('机器翻译-阿里云'));
    final tencentTitle = tester.getRect(find.text('机器翻译-腾讯云'));
    final baiduTitle = tester.getRect(find.text('机器翻译-百度云'));
    expect(tencentTitle.top, greaterThan(aliyunTitle.bottom));
    expect(baiduTitle.top, greaterThan(tencentTitle.bottom));
    expect(aliyunTitle.left, tencentTitle.left);
    expect(find.text('复制计费说明'), findsNothing);
    expect(find.text('复制官方网址'), findsNWidgets(7));
    final copyButtons = find.ancestor(
        of: find.text('复制官方网址'), matching: find.byType(TextButton));
    for (final button in copyButtons.evaluate()) {
      final buttonRect = tester.getRect(find.byWidget(button.widget));
      expect(buttonRect.right, closeTo(engine.right - 22, 1));
    }
    expect(find.text('二维码待替换'), findsNothing);
    for (final entry in const [
      (Key('donation-wechat-qr'), 'assets/donation_wechat.png'),
      (Key('donation-alipay-qr'), 'assets/donation_alipay.png'),
    ]) {
      final image = tester.widget<Image>(find.descendant(
          of: find.byKey(entry.$1), matching: find.byType(Image)));
      expect((image.image as AssetImage).assetName, entry.$2);
    }
    await tester.pumpWidget(const SizedBox.shrink());
  });

  testWidgets('about cards stack in a narrow window', (tester) async {
    await tester.binding.setSurfaceSize(const Size(470, 700));
    addTearDown(() => tester.binding.setSurfaceSize(null));
    await tester.pumpWidget(const AntoApp(autoConnect: false));
    await tester.pumpAndSettle();
    await tester.tap(find.byIcon(Icons.info_outline));
    await tester.pumpAndSettle();

    final intro = tester.getRect(find.byKey(const Key('about-intro-card')));
    final author = tester.getRect(find.byKey(const Key('about-author-card')));
    final donation =
        tester.getRect(find.byKey(const Key('about-donation-card')));
    final engine = tester.getRect(find.byKey(const Key('about-engine-card')));
    expect(intro.top, lessThan(author.top));
    expect(author.top, lessThan(donation.top));
    expect(donation.top, lessThan(engine.top));
    expect(tester.takeException(), isNull);
    await tester.pumpWidget(const SizedBox.shrink());
  });

  testWidgets('about address buttons show copy success and failure dialogs',
      (tester) async {
    await tester.binding.setSurfaceSize(const Size(1264, 789));
    addTearDown(() => tester.binding.setSurfaceSize(null));
    await tester.pumpWidget(const AntoApp(autoConnect: false));
    await tester.pumpAndSettle();
    await tester.tap(find.text('关于应用'));
    await tester.pumpAndSettle();

    var fail = false;
    tester.binding.defaultBinaryMessenger.setMockMethodCallHandler(
      SystemChannels.platform,
      (call) async {
        if (call.method == 'Clipboard.setData' && fail) {
          throw PlatformException(code: 'clipboard-unavailable');
        }
        return null;
      },
    );
    addTearDown(() => tester.binding.defaultBinaryMessenger
        .setMockMethodCallHandler(SystemChannels.platform, null));

    await tester.tap(find.text('复制项目地址'));
    await tester.pumpAndSettle();
    expect(find.text('复制成功'), findsOneWidget);
    expect(find.text('已复制项目地址'), findsOneWidget);
    await tester.tap(find.text('知道了'));
    await tester.pumpAndSettle();

    fail = true;
    await tester.tap(find.text('复制反馈地址'));
    await tester.pumpAndSettle();
    expect(find.text('复制失败'), findsOneWidget);
    expect(find.text('反馈地址未能复制，请重试'), findsOneWidget);
    await tester.tap(find.text('知道了'));
    await tester.pumpAndSettle();
  });

  testWidgets('workspace and queue cards fill the available height',
      (tester) async {
    await tester.binding.setSurfaceSize(const Size(1264, 789));
    addTearDown(() => tester.binding.setSurfaceSize(null));
    await tester.pumpWidget(const AntoApp(autoConnect: false));
    await tester.pumpAndSettle();

    final composer =
        tester.getRect(find.byKey(const Key('workspace-composer-card')));
    final recent =
        tester.getRect(find.byKey(const Key('workspace-recent-card')));
    expect(composer.top, 30);
    expect(composer.bottom, 745);
    expect(recent.top, composer.top);
    expect(recent.bottom, composer.bottom);
    final clear = tester.getRect(find.ancestor(
        of: find.text('重置'), matching: find.byType(OutlinedButton)));
    final submit = tester.getRect(find.ancestor(
        of: find.text('开始翻译'), matching: find.byType(FilledButton)));
    expect(submit.left - clear.right, 9);
    expect(composer.right - submit.right, inInclusiveRange(22, 24));
    expect(composer.bottom - submit.bottom, inInclusiveRange(22, 24));
    expect(
        find.descendant(
            of: find.byKey(const Key('workspace-composer-card')),
            matching: find.byType(Scrollable)),
        findsWidgets);
    expect(
        find.descendant(
            of: find.byKey(const Key('workspace-recent-card')),
            matching: find.byType(Scrollable)),
        findsWidgets);

    await tester.tap(find.text('任务队列'));
    await tester.pumpAndSettle();
    final queue = tester.getRect(find.byKey(const Key('queue-tasks-card')));
    expect(queue.top, 30);
    expect(queue.bottom, 745);
    final title = tester.getRect(find.text('任务记录'));
    final refresh = tester.getRect(find.text('自动刷新'));
    final newTask = tester.getRect(find.ancestor(
        of: find.text('新建任务'), matching: find.byType(FilledButton)));
    expect(refresh.left - title.right, 8);
    expect(refresh.center.dy, title.center.dy);
    expect(queue.right - newTask.right, inInclusiveRange(20, 22));
    expect(newTask.center.dy, title.center.dy);
    expect(
        find.descendant(
            of: find.byKey(const Key('queue-tasks-card')),
            matching: find.byType(Scrollable)),
        findsWidgets);
    await tester.tap(find.text('新建任务'));
    await tester.pumpAndSettle();
    expect(find.text('新建翻译任务'), findsOneWidget);
  });

  testWidgets('overflow scrolls inside each task card', (tester) async {
    await tester.binding.setSurfaceSize(const Size(1264, 360));
    addTearDown(() => tester.binding.setSurfaceSize(null));
    final core = _TaskCore();
    for (var index = 0; index < 20; index++) {
      core.tasks.add({
        'id': 'task-$index',
        'status': 'failed',
        'total': 1,
        'completed': 1,
        'error': '字幕翻译失败，请检查翻译服务连接和参数配置后重新提交任务。',
        'output_files': ['C:\\subtitles\\anto\\task-$index.srt'],
        'created_at': DateTime.now().toIso8601String(),
      });
    }
    await tester.pumpWidget(AntoApp(core: core));
    await tester.pumpAndSettle();

    for (final key in const [
      Key('workspace-composer-card'),
      Key('workspace-recent-card'),
    ]) {
      final card = find.byKey(key);
      final scrollable =
          find.descendant(of: card, matching: find.byType(Scrollable)).first;
      final position = tester.state<ScrollableState>(scrollable).position;
      await tester.drag(card, const Offset(0, -200));
      await tester.pumpAndSettle();
      expect(position.pixels, greaterThan(0));
      expect(tester.getRect(card).bottom, 316);
    }

    await tester.tap(find.text('任务队列'));
    await tester.pumpAndSettle();
    final queue = find.byKey(const Key('queue-tasks-card'));
    final scrollable =
        find.descendant(of: queue, matching: find.byType(Scrollable)).first;
    final position = tester.state<ScrollableState>(scrollable).position;
    await tester.drag(queue, const Offset(0, -200));
    await tester.pumpAndSettle();
    expect(position.pixels, greaterThan(0));
    expect(tester.getRect(queue).bottom, 316);
  });

  testWidgets('shows selected SRT files and allows removing one',
      (tester) async {
    await tester
        .pumpWidget(const AntoApp(autoConnect: false, picker: _TestPicker()));
    await tester.tap(find.text('选择文件'));
    await tester.pumpAndSettle();

    expect(find.text('已选择 2 个 SRT 文件'), findsOneWidget);
    expect(find.text('sample.srt'), findsOneWidget);
    expect(find.text('second.srt'), findsOneWidget);
    expect(find.text(r'C:\subtitles\sample.srt'), findsOneWidget);

    await tester.tap(find.byTooltip('移除 sample.srt'));
    await tester.pumpAndSettle();
    expect(find.text('已选择 1 个 SRT 文件'), findsOneWidget);
    expect(find.text('sample.srt'), findsNothing);
    expect(find.text('second.srt'), findsOneWidget);
  });

  testWidgets('keeps language and translation controls in paired rows',
      (tester) async {
    await tester.binding.setSurfaceSize(const Size(1200, 900));
    addTearDown(() => tester.binding.setSurfaceSize(null));
    await tester.pumpWidget(const AntoApp(autoConnect: false));
    await tester.pumpAndSettle();

    final sourceY = tester.getTopLeft(find.text('源语言')).dy;
    final targetY = tester.getTopLeft(find.text('目标语言')).dy;
    final engineY = tester.getTopLeft(find.text('翻译引擎')).dy;
    final modeY = tester.getTopLeft(find.text('翻译模式')).dy;
    final directoryY = tester.getTopLeft(find.text('输入目录')).dy;
    expect(sourceY, targetY);
    expect(engineY, modeY);
    expect(engineY, greaterThan(sourceY));
    expect(directoryY, greaterThan(engineY));
  });

  testWidgets('reset restores the task form defaults', (tester) async {
    final core = _TaskCore();
    await tester.pumpWidget(AntoApp(core: core, picker: const _TestPicker()));
    await tester.pumpAndSettle();
    await tester.tap(find.text('选择文件'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('选择目录'));
    await tester.pumpAndSettle();

    final dropdowns = find.byType(DropdownButtonFormField<String>);
    await tester.ensureVisible(dropdowns.at(0));
    await tester.pumpAndSettle();
    await tester.tap(dropdowns.at(0));
    await tester.pumpAndSettle();
    await tester.tap(find.text('中文').last);
    await tester.pumpAndSettle();
    await tester.ensureVisible(dropdowns.at(1));
    await tester.pumpAndSettle();
    await tester.tap(dropdowns.at(1));
    await tester.pumpAndSettle();
    await tester.tap(find.text('英语').last);
    await tester.pumpAndSettle();
    await tester.ensureVisible(dropdowns.at(3));
    await tester.pumpAndSettle();
    await tester.tap(dropdowns.at(3));
    await tester.pumpAndSettle();
    await tester.tap(find.text('增量翻译').last);
    await tester.pumpAndSettle();

    expect(tester.state<FormFieldState<String>>(dropdowns.at(0)).value, 'zh');
    expect(tester.state<FormFieldState<String>>(dropdowns.at(1)).value, 'en');
    expect(tester.state<FormFieldState<String>>(dropdowns.at(3)).value, '增量翻译');
    await tester.tap(find.text('重置'));
    await tester.pumpAndSettle();

    expect(find.text('已选择 2 个 SRT 文件'), findsNothing);
    expect(find.text('选择 SRT 文件或目录'), findsOneWidget);
    final inputField = tester.widget<TextField>(find.byWidgetPredicate(
        (widget) =>
            widget is TextField &&
            widget.decoration?.hintText == '可选，扫描目录中的 .srt'));
    expect(inputField.controller?.text, isEmpty);
    expect(tester.state<FormFieldState<String>>(dropdowns.at(0)).value, 'en');
    expect(tester.state<FormFieldState<String>>(dropdowns.at(1)).value, 'zh');
    expect(tester.state<FormFieldState<String>>(dropdowns.at(2)).value,
        'ali_cloud_mt');
    expect(tester.state<FormFieldState<String>>(dropdowns.at(3)).value, '全量翻译');
    await tester.pumpWidget(const SizedBox.shrink());
  });

  testWidgets('resets task inputs after successful submission', (tester) async {
    final core = _TaskCore();
    await tester.pumpWidget(AntoApp(core: core, picker: const _TestPicker()));
    await tester.pumpAndSettle();

    await tester.tap(find.text('选择文件'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('选择目录'));
    await tester.pumpAndSettle();
    expect(find.text('已选择 2 个 SRT 文件'), findsOneWidget);

    await tester.ensureVisible(find.text('开始翻译'));
    await tester.tap(find.text('开始翻译'));
    await tester.pumpAndSettle();

    expect(core.lastTaskRequest?['files'], hasLength(2));
    expect(core.lastTaskRequest?['directory'], r'C:\subtitles');
    expect(find.text('已选择 2 个 SRT 文件'), findsNothing);
    expect(find.text('选择 SRT 文件或目录'), findsOneWidget);
    final inputField = tester.widget<TextField>(find.byWidgetPredicate(
        (widget) =>
            widget is TextField &&
            widget.decoration?.hintText == '可选，扫描目录中的 .srt'));
    expect(inputField.controller?.text, isEmpty);

    await tester.tap(find.text('知道了'));
    await tester.pumpAndSettle();
    await tester.pumpWidget(const SizedBox.shrink());
  });

  testWidgets('completed task shows its output filename', (tester) async {
    final core = _TaskCore();
    core.tasks.add({
      'id': 'completed-1',
      'status': 'completed',
      'total': 1,
      'completed': 1,
      'output_files': [r'C:\subtitles\anto\sample.srt'],
      'created_at': DateTime.now().toIso8601String(),
    });
    await tester.pumpWidget(AntoApp(core: core));
    await tester.pumpAndSettle();

    expect(find.text('sample.srt'), findsOneWidget);
    expect(find.text('等待处理'), findsNothing);
    await tester.pumpWidget(const SizedBox.shrink());
  });
}
