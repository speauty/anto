# ANTO 字幕翻译助手

ANTO 是 Windows 字幕翻译应用。`desktop/` 是 Flutter 客户端，`cmd/anto-core/` 和 `core/` 提供本地 Go 服务。目前支持“机器翻译-阿里云”“机器翻译-腾讯云”“机器翻译-百度云”“彩云小译”“DeepL Pro”“AI 翻译-DeepSeek”和“批量翻译-有道智云”引擎。

## 使用

从发布目录运行 `anto-ui.exe`。`anto-ui.exe`、`anto-core.exe`、`cfg.yml`、`data/` 和 `flutter_windows.dll` 必须位于同一目录。客户端启动 Core 后，通过回环地址上的随机端口和一次性 Token 与其通信。

在设置页配置所选引擎：阿里云填写 AccessKey ID、AccessKey Secret 和地区；腾讯云填写 SecretId、SecretKey 和地区；百度云填写 App ID、密钥并选择已开通的服务版本；彩云小译需在开放平台单独申请 API Token；DeepL Pro 使用已有 Pro 或当前付费 API 套餐的 API Key；DeepSeek 填写 API Key，可调整服务地址、模型名和翻译提示词；有道智云需开通批量翻译服务并填写应用 ID、应用密钥。DeepSeek 的服务地址也可填写自建兼容 Chat Completions 的服务，例如 `http://127.0.0.1:8000/v1`；本机自建服务可不填 API Key。选择 SRT 文件或输入目录、语言和翻译模式后，点击“开始翻译”。任务成功加入队列后，表单会重置。默认输出目录为源文件所在目录的 `anto/` 子目录。

设置页的“测试连接”会检查所选引擎接入点的网络连通性，不调用翻译接口，也不验证密钥。百度云、彩云小译、DeepL Pro、DeepSeek 和有道智云无需选择地区。

## 开发

Go Core：在仓库根目录运行 `go test ./...`、`go build -o ./bin/anto-core.exe ./cmd/anto-core`。

Flutter：在 `desktop/` 运行 `flutter pub get`、`flutter analyze`、`flutter test`、`flutter build windows --release`。发布时将 Core 和 `cfg.yml` 复制到 `desktop/build/windows/x64/runner/Release/`。

本地 API 见 [docs/core-api.md](docs/core-api.md)。
