# ANTO 项目开发说明

本文件适用于仓库根目录下的全部代码。项目是 Windows 优先的字幕翻译应用，包含 Go 后端和两套 Windows 客户端：旧版 Walk GUI，以及正在迁移中的 Flutter 客户端。

## 项目结构

- `cmd/anto/`：旧版桌面程序入口。启动后执行 `bootstrap.Boot`，再运行 `platform/win` 中的 Walk GUI。
- `cmd/anto-core/`：独立 Go Core 入口。启动时初始化配置和翻译引擎，在 `127.0.0.1:0` 上监听随机端口，并输出一行 `ANTO_CORE_READY {"port":...,"token":"..."}`。
- `core/`：与平台无关的任务服务，负责收集 SRT 文件、读取和解析字幕、调用翻译引擎、写出结果、取消任务和发布任务事件。
- `core/api/`：Core 的 HTTP API。端点和 SSE 事件格式见 `docs/core-api.md`。
- `domain/service/translator/`：翻译引擎接口、请求工具和各厂商实现。每个引擎通常包含 `config.go`、`lang.go` 和 `translator.go`。
- `domain/repository/`：翻译引擎注册表；引擎只有在 `bootstrap/boot.go` 初始化并注册后才可用。
- `cron/`、`lib/srt/`：旧版字幕读取、检测、翻译、写出流程和 SRT 编解码基础库。
- `cfg/`、`cfg.yml`：Viper 配置模型及根目录配置文件。配置加载使用当前工作目录下的 `cfg.yml`。
- `desktop/`：Flutter Windows 客户端。`desktop/lib/core_client.dart` 启动同目录的 `anto-core.exe`，读取随机端口和一次性 Token 后访问 Core。
- `platform/win/`、`assets/`、`resource/`：Walk UI、Windows 资源和旧版界面素材。

## 常用命令

所有 Go 命令应在仓库根目录执行，Go 版本按 `go.mod` 使用 Go 1.21。

```powershell
go mod tidy
go test ./...
go vet ./...
go build -o .\bin\anto-core.exe anto/cmd/anto-core
make build       # 旧版 Walk GUI，输出 bin/anto-v3.7.0-windows.exe
make deploy      # 生成 Windows 资源后构建旧版 GUI
make compress    # 使用 upx 压缩旧版 GUI（需要本机安装 upx）
```

`make rs` 依赖 `rsrc`，只在需要更新 `cmd/anto/rsrc.syso` 时运行。Windows GUI 构建依赖 Windows 图形环境；Core 和大部分业务代码可以单独编译测试。

Flutter 命令在 `desktop/` 目录执行：

```powershell
flutter pub get
flutter analyze
flutter test
flutter build windows
```

发布 Flutter 客户端时，`anto-core.exe` 必须和生成的客户端可执行文件放在同一目录。客户端不应假设固定端口，也不应绕过 `CoreClient` 自行拼接认证请求。

如果 Go 因机器全局缓存目录权限失败，可先为当前 PowerShell 会话指定可写缓存，例如：

```powershell
$env:GOCACHE = Join-Path $env:TEMP 'anto-gocache'
go test ./...
```

## 配置和运行约定

- `bootstrap.Boot` 是两个 Go 入口共用的初始化路径：安装资源、加载 `cfg.yml`、初始化各翻译引擎并注册到 repository。新增或修改翻译引擎时，需要同步检查配置结构、默认值、初始化和注册代码。
- `cfg.yml` 中的 AK/SK、Token、API Key 等字段属于敏感信息。提交前只保留空值或示例值，不要把真实凭据写入仓库、测试夹具或日志。
- `cfg.Load` 从当前工作目录查找 `cfg.yml`；启动 Go 程序时应保证工作目录包含该文件，或者在文档和脚本中明确设置工作目录。
- Core API 除 `/api/v1/health` 外都要求 `X-Anto-Token`。Core 只能绑定回环地址，Token 每次进程启动重新生成；修改 API 时要保持这两个边界。
- 任务状态包括 `queued`、`reading`、`translating`、`writing`、`completed`、`failed`、`cancelled`。长任务必须传递 `context.Context`，及时响应取消并发布对应事件。
- 请求和响应字段沿用 `docs/core-api.md` 中的 snake_case 命名。修改 API 时同时更新该文档、Go 测试和 Flutter 客户端。
- 任务只收集 `.srt` 文件；目录输入只扫描目录的直接子项。默认输出目录是源文件目录下的 `anto/`，除非请求显式提供 `output_directory`。增量翻译会跳过已有副轨字幕的块。

## 编码约定

- Go 文件提交前运行 `gofmt`；错误应携带上下文并向上返回，避免无理由吞掉错误或使用裸 `panic`。新增并发代码时保护任务、取消器和订阅者集合的共享状态。
- 翻译引擎实现 `domain/service/translator.ImplTranslator`。网络调用应复用 `domain/service/translator/request.go` 的限流、重试和错误处理能力，并遵守各引擎的最大文本长度和 QPS 配置。
- 新增翻译引擎至少补齐配置、语言映射、实现、`cfg.Cfg` 字段、`cfg.yml` 节点、`bootstrap.Boot` 初始化和 repository 注册；同时补充可验证的单元测试或 API 测试。
- Flutter 代码遵循 `desktop/analysis_options.yaml` 的 lint 规则，提交前运行 `dart format lib test`（存在测试目录时）和 `flutter analyze`。与 Core 的进程生命周期、认证和 SSE 订阅集中在通信层处理。
- 变更优先保持旧版 Walk GUI、Core API 和 Flutter 客户端之间的行为兼容；涉及协议字段、任务状态或输出文件格式时，先更新 `docs/core-api.md` 再改实现。

## 测试重点

- Core API 至少覆盖随机回环端口、Token 鉴权、HTTP 方法、任务创建和取消、任务状态转换及 SSE 事件。
- 字幕逻辑重点覆盖 UTF-8/带 BOM 输入、空字幕块、增量模式、分段长度限制、双轨/主轨/副轨输出和输出目录创建。
- 翻译服务的网络测试应避免依赖真实供应商额度；使用接口替身或固定响应验证请求参数、错误处理和限流行为。
- Flutter 改动至少运行 `flutter analyze`；涉及 CoreClient 或页面交互时补充 widget/unit 测试，避免依赖真实 `anto-core.exe` 和外部翻译服务。
