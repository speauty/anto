# anto Windows 客户端

这是 Flutter Windows 客户端。客户端启动同目录下的 `anto-core.exe`，读取 `ANTO_CORE_READY` 握手信息后，通过 `127.0.0.1` 的随机端口和一次性 Token 调用 Go Core。

本目录包含 Flutter 业务代码、通信层和 Windows runner。发布后运行 `anto-ui.exe`，客户端会从同目录启动 `anto-core.exe`。发布目录必须同时包含 `cfg.yml`、`data/` 和 `flutter_windows.dll`。

当前仅保留“机器翻译-阿里云”引擎。用户在设置页填写 AccessKey ID、AccessKey Secret 和地区。
