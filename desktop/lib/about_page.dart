import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

const _brand = Color(0xff5e61e9);
const _ink = Color(0xff252940);
const _muted = Color(0xff858ba7);

class AboutPage extends StatelessWidget {
  const AboutPage({super.key, required this.version});

  final String version;

  static const projectUrl = 'https://github.com/speauty/anto';
  static const feedbackUrl = '$projectUrl/issues';
  static const aliyunGuideUrl = 'https://www.aliyun.com/product/ai/alimt';
  static const tencentGuideUrl =
      'https://cloud.tencent.com/document/product/551/104415';
  static const baiduGuideUrl = 'https://fanyi-api.baidu.com/doc/21';
  static const caiyunGuideUrl = 'https://docs.caiyunapp.com/lingocloud-api/';
  static const deeplGuideUrl =
      'https://developers.deepl.com/api-reference/translate/request-translation';
  static const deepseekGuideUrl =
      'https://api-docs.deepseek.com/quick_start/pricing/';
  static const youdaoGuideUrl =
      'https://ai.youdao.com/DOCSIRMA/html/trans/api/plwbfy/index.html';

  Future<void> _copy(BuildContext context, String value, String label) async {
    var copied = false;
    try {
      await Clipboard.setData(ClipboardData(text: value));
      copied = true;
    } catch (_) {
      copied = false;
    }
    if (!context.mounted) return;
    await showDialog<void>(
      context: context,
      builder: (dialogContext) => AlertDialog(
        title: Row(children: [
          Icon(copied ? Icons.check_circle_outline : Icons.error_outline,
              color: copied ? _brand : const Color(0xffd9534f)),
          const SizedBox(width: 9),
          Text(copied ? '复制成功' : '复制失败'),
        ]),
        content: Text(copied ? '已复制$label' : '$label未能复制，请重试'),
        actions: [
          FilledButton(
            onPressed: () => Navigator.of(dialogContext).pop(),
            child: const Text('知道了'),
          ),
        ],
      ),
    );
  }

  Widget _card(Widget child, Key key) => Container(
        key: key,
        clipBehavior: Clip.antiAlias,
        decoration: BoxDecoration(
          color: Colors.white,
          border: Border.all(color: const Color(0xffe8ebf2)),
          borderRadius: BorderRadius.circular(16),
          boxShadow: const [
            BoxShadow(
              color: Color(0x19293152),
              blurRadius: 30,
              offset: Offset(0, 12),
            ),
          ],
        ),
        child: child,
      );

  Widget _introCard() => _card(
        SingleChildScrollView(
          child: Padding(
            padding: const EdgeInsets.all(24),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Wrap(
                  crossAxisAlignment: WrapCrossAlignment.center,
                  spacing: 18,
                  runSpacing: 14,
                  children: [
                    Image.asset('assets/LOGO.png', width: 68, height: 68),
                    Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        const Text(
                          'Anto-你的字幕翻译小助手',
                          style: TextStyle(
                            color: _ink,
                            fontSize: 22,
                            fontWeight: FontWeight.w700,
                          ),
                        ),
                        const SizedBox(height: 6),
                        Text('版本 $version',
                            style:
                                const TextStyle(color: _muted, fontSize: 13)),
                      ],
                    ),
                  ],
                ),
                const SizedBox(height: 22),
                const Text(
                  '让字幕翻译更轻松',
                  style: TextStyle(
                    color: _ink,
                    fontSize: 17,
                    fontWeight: FontWeight.w600,
                  ),
                ),
                const SizedBox(height: 8),
                const Text(
                  '选择 SRT 文件或目录，设置语言与翻译模式，Anto 会将字幕加入任务队列，并按你的输出设置导出翻译结果。',
                  style: TextStyle(color: _muted, height: 1.5),
                ),
                const SizedBox(height: 14),
                const Wrap(spacing: 8, runSpacing: 8, children: [
                  _FeatureTag('SRT 字幕'),
                  _FeatureTag('批量任务'),
                  _FeatureTag('多轨导出'),
                ]),
              ],
            ),
          ),
        ),
        const Key('about-intro-card'),
      );

  Widget _authorCard(BuildContext context) => _card(
        Padding(
          padding: const EdgeInsets.all(22),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text('作者与项目',
                  style: TextStyle(
                      color: _ink, fontSize: 16, fontWeight: FontWeight.w700)),
              const SizedBox(height: 5),
              const Text('感谢你使用 Anto',
                  style: TextStyle(color: _muted, fontSize: 12)),
              const SizedBox(height: 20),
              const Expanded(
                child: SingleChildScrollView(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text('speauty',
                          style: TextStyle(
                              color: _ink,
                              fontSize: 20,
                              fontWeight: FontWeight.w700)),
                      SizedBox(height: 9),
                      Text('Anto 由 speauty 开发与维护，希望让字幕翻译中重复的操作变得简单。',
                          style: TextStyle(color: _muted, height: 1.5)),
                      SizedBox(height: 20),
                      Text('项目地址',
                          style: TextStyle(color: _muted, fontSize: 12)),
                      SizedBox(height: 5),
                      SelectableText(projectUrl,
                          style: TextStyle(color: _brand, fontSize: 13)),
                    ],
                  ),
                ),
              ),
              const SizedBox(height: 12),
              Align(
                alignment: Alignment.centerRight,
                child: Wrap(
                  alignment: WrapAlignment.end,
                  spacing: 10,
                  runSpacing: 8,
                  children: [
                    OutlinedButton(
                      onPressed: () => _copy(context, feedbackUrl, '反馈地址'),
                      child: const Text('复制反馈地址'),
                    ),
                    OutlinedButton(
                      onPressed: () => _copy(context, projectUrl, '项目地址'),
                      child: const Text('复制项目地址'),
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
        const Key('about-author-card'),
      );

  Widget _donationCard() => _card(
        const SingleChildScrollView(
          child: Padding(
            padding: EdgeInsets.all(22),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('支持开发',
                    style: TextStyle(
                        color: _ink,
                        fontSize: 16,
                        fontWeight: FontWeight.w700)),
                SizedBox(height: 5),
                Text('一份鼓励，让 Anto 继续变得更好',
                    style: TextStyle(color: _muted, fontSize: 12)),
                SizedBox(height: 20),
                Text(
                  '如果 Anto 帮你节省了时间，欢迎支持它的持续维护与改进。无论是否打赏，你的建议和反馈都很珍贵。',
                  style: TextStyle(color: _muted, height: 1.5),
                ),
                SizedBox(height: 22),
                Wrap(
                  spacing: 16,
                  runSpacing: 16,
                  children: [
                    _DonationQr(
                        label: '微信赞赏',
                        asset: 'assets/donation_wechat.png',
                        key: Key('donation-wechat-qr')),
                    _DonationQr(
                        label: '支付宝赞赏',
                        asset: 'assets/donation_alipay.png',
                        key: Key('donation-alipay-qr')),
                  ],
                ),
              ],
            ),
          ),
        ),
        const Key('about-donation-card'),
      );

  Widget _engineLine(String title, String detail) => Padding(
        padding: const EdgeInsets.only(top: 9),
        child: Text.rich(
          TextSpan(children: [
            TextSpan(
                text: '$title  ',
                style:
                    const TextStyle(color: _ink, fontWeight: FontWeight.w600)),
            TextSpan(text: detail),
          ]),
          style: const TextStyle(color: _muted, height: 1.45, fontSize: 13),
        ),
      );

  Widget _engineDetails(
    BuildContext context, {
    required String name,
    required List<(String, String)> lines,
    required String guideUrl,
  }) =>
      Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(children: [
            Expanded(
              child: Text(name,
                  style: const TextStyle(
                      color: _ink, fontSize: 15, fontWeight: FontWeight.w700)),
            ),
            const SizedBox(width: 12),
            TextButton(
                onPressed: () => _copy(context, guideUrl, '官方网址'),
                child: const Text('复制官方网址')),
          ]),
          for (final line in lines) _engineLine(line.$1, line.$2),
        ],
      );

  Widget _engineCard(BuildContext context) => _card(
        Padding(
          padding: const EdgeInsets.all(22),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text('翻译引擎',
                  style: TextStyle(
                      color: _ink, fontSize: 16, fontWeight: FontWeight.w700)),
              const SizedBox(height: 5),
              const Text('开通、使用与费用参考；实际额度和价格以服务商公告为准',
                  style: TextStyle(color: _muted, fontSize: 12)),
              const SizedBox(height: 18),
              _engineDetails(
                context,
                name: '机器翻译-阿里云',
                lines: const [
                  ('开通', '在阿里云控制台开通机器翻译通用版，创建有调用权限的 AccessKey。'),
                  ('使用', '在“我的设置 → 翻译服务 → 配置参数”填写密钥和地域，保存后选择此引擎。'),
                  ('额度', '每月 100 万字符免费，主子账号共享，当月未用不结转。'),
                  ('计费', '按文本字符计费；依次抵扣免费额度、资源包，超出后按 50 元/百万字符后付费。'),
                ],
                guideUrl: aliyunGuideUrl,
              ),
              const Divider(height: 32),
              _engineDetails(
                context,
                name: '机器翻译-腾讯云',
                lines: const [
                  ('开通', '传统文本翻译已不支持新用户；现有用户可在腾讯云控制台查看服务与密钥。'),
                  (
                    '使用',
                    '在“我的设置 → 翻译服务 → 配置参数”填写 SecretId、SecretKey 和地域，保存后选择此引擎。'
                  ),
                  ('额度', '自 2026 年 10 月起不再发放新免费资源包，已有资源包仅在有效期内可用。'),
                  ('计费', '先抵扣可用资源包；手动开启后付费后，文本翻译低用量档为 58 元/百万字符。'),
                ],
                guideUrl: tencentGuideUrl,
              ),
              const Divider(height: 32),
              _engineDetails(
                context,
                name: '机器翻译-百度云',
                lines: const [
                  ('开通', '在百度翻译开放平台开通通用文本翻译 API，取得 App ID 和密钥。'),
                  ('使用', '在“我的设置 → 翻译服务 → 配置参数”填写凭据并选择实际开通的服务版本，保存后使用。'),
                  ('额度', '标准版每月 5 万字符、高级版 100 万、尊享版 200 万；QPS 分别为 1、10、100。'),
                  ('计费', '标准版和高级版超出免费额度按 49 元/百万字符计费；实际权益和价格以百度控制台为准。'),
                ],
                guideUrl: baiduGuideUrl,
              ),
              const Divider(height: 32),
              _engineDetails(
                context,
                name: '彩云小译',
                lines: const [
                  ('开通', '在彩云科技开放平台注册，申请小译 API Token；App 会员不能替代 API Token。'),
                  (
                    '使用',
                    '在“我的设置 → 翻译服务 → 配置参数”填写 Token，保存后选择彩云小译。支持中文与英、日、韩互译，以及其他已列语种译为中文。'
                  ),
                  ('额度', '官方文档说明新用户赠送 100 万字符，有效期一个月。'),
                  ('计费', '超出免费额度后按 39 元/百万字符计费，空格和标点也计入字符数。'),
                  ('来源', '数据来自彩云小译'),
                ],
                guideUrl: caiyunGuideUrl,
              ),
              const Divider(height: 32),
              _engineDetails(
                context,
                name: 'DeepL Pro',
                lines: const [
                  (
                    '开通',
                    '使用已有 DeepL API Pro 或当前付费 API 套餐的密钥；网页翻译会员不能替代 API 套餐。'
                  ),
                  ('使用', '在“我的设置 → 翻译服务 → 配置参数”填写 API Key，保存后选择 DeepL Pro。'),
                  ('额度', 'API Pro 无免费字符，可在 DeepL 账户设置费用或密钥用量上限。'),
                  ('计费', '按月基本费及成功请求的源文本字符数计费；具体价格以账户页面为准。'),
                ],
                guideUrl: deeplGuideUrl,
              ),
              const Divider(height: 32),
              _engineDetails(
                context,
                name: 'AI 翻译-DeepSeek',
                lines: const [
                  (
                    '开通',
                    '在 DeepSeek 开放平台创建 API Key；也可使用自建的兼容 Chat Completions 服务。'
                  ),
                  (
                    '使用',
                    '在“我的设置 → 翻译服务 → 配置参数”填写 API Key；可调整服务地址、模型名和翻译提示词，保存后选择此引擎。'
                  ),
                  (
                    '模型',
                    '官方模型可选 deepseek-flash 或 deepseek-v4-pro；自建服务填写其实际模型名。'
                  ),
                  ('计费', '官方 API 按输入和输出 Token 计费，具体价格以 DeepSeek 文档为准。'),
                ],
                guideUrl: deepseekGuideUrl,
              ),
              const Divider(height: 32),
              _engineDetails(
                context,
                name: '批量翻译-有道智云',
                lines: const [
                  ('开通', '在有道智云创建应用并绑定批量翻译服务，取得应用 ID 和应用密钥。'),
                  ('使用', '在“我的设置 → 翻译服务 → 配置参数”填写凭据，保存后选择此引擎。多条字幕会合并请求。'),
                  ('额度', '官方文档说明每个账户赠送 50 元体验金。'),
                  ('计费', '按批量翻译服务报价计费，实际价格以有道智云控制台为准。'),
                ],
                guideUrl: youdaoGuideUrl,
              ),
            ],
          ),
        ),
        const Key('about-engine-card'),
      );

  @override
  Widget build(BuildContext context) => LayoutBuilder(
        builder: (context, constraints) {
          final stacked =
              constraints.maxWidth < 700 || constraints.maxHeight < 500;
          return SingleChildScrollView(
            child: Column(children: [
              SizedBox(
                  height: stacked ? 300 : 250,
                  width: double.infinity,
                  child: _introCard()),
              const SizedBox(height: 16),
              if (stacked) ...[
                SizedBox(height: 290, child: _authorCard(context)),
                const SizedBox(height: 16),
                SizedBox(height: 390, child: _donationCard()),
              ] else
                SizedBox(
                  height: 360,
                  child: Row(
                    crossAxisAlignment: CrossAxisAlignment.stretch,
                    children: [
                      Expanded(child: _authorCard(context)),
                      const SizedBox(width: 16),
                      Expanded(child: _donationCard()),
                    ],
                  ),
                ),
              const SizedBox(height: 16),
              _engineCard(context),
            ]),
          );
        },
      );
}

class _FeatureTag extends StatelessWidget {
  const _FeatureTag(this.label);
  final String label;

  @override
  Widget build(BuildContext context) => Container(
        padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
        decoration: BoxDecoration(
          color: const Color(0xfff0f1ff),
          borderRadius: BorderRadius.circular(8),
        ),
        child: Text(label, style: const TextStyle(color: _brand, fontSize: 12)),
      );
}

class _DonationQr extends StatelessWidget {
  const _DonationQr({super.key, required this.label, required this.asset});
  final String label;
  final String asset;

  @override
  Widget build(BuildContext context) => Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Container(
            width: 196,
            height: 196,
            decoration: BoxDecoration(
              color: Colors.white,
              border: Border.all(color: const Color(0xffe8ebf2)),
              borderRadius: BorderRadius.circular(13),
            ),
            child: ClipRRect(
              borderRadius: BorderRadius.circular(12),
              child: Image.asset(asset,
                  fit: BoxFit.contain, filterQuality: FilterQuality.high),
            ),
          ),
          const SizedBox(height: 9),
          Text(label, style: const TextStyle(color: _ink, fontSize: 13)),
        ],
      );
}
