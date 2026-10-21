# Flutter Windows 客户端状态

更新：2026-10-19

Flutter Windows 客户端是唯一桌面界面，构建产物名为 `anto-ui.exe`。Go Core 独立运行在 `anto-core.exe` 中，提供配置、翻译任务和 SSE 事件 API。旧 Walk 客户端已移除。

当前接入“机器翻译-阿里云”“机器翻译-腾讯云”和“机器翻译-百度云”。设置页的地区下拉框对应阿里云和腾讯云支持的地区与公网接入地址；百度云选择服务版本以控制请求长度和 QPS。用户提交任务成功后，工作台表单重置；已完成任务显示输出文件名。

待验收：使用真实供应商凭据分别验证翻译、地区切换、全量/增量、轨道导出和异常恢复。自动化测试不调用真实供应商服务。Core 任务保存在内存中，进程重启后历史任务不会保留。

开发验证：`go test ./...`、`flutter analyze`、`flutter test`、`flutter build windows --release`。本地协议详见 [core-api.md](core-api.md)。
