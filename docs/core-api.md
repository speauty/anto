# Go Core 本地 API

`anto-core.exe` 只监听 `127.0.0.1` 的随机 TCP 端口。启动后向标准输出写一行：

```text
ANTO_CORE_READY {"port":43127,"token":"..."}
```

除健康检查外，所有请求必须携带 `X-Anto-Token` 请求头。Token 每次进程启动时重新生成。请求和响应字段使用 `snake_case`。

错误响应统一为 JSON：

```json
{"error":{"code":"task_not_found","message":"任务不存在"}}
```

## 端点

```text
GET  /api/v1/health
GET  /api/v1/translators
POST /api/v1/translators/{id}/test-connection
GET  /api/v1/config
PUT  /api/v1/config
POST /api/v1/tasks
GET  /api/v1/tasks
GET  /api/v1/tasks/{id}
POST /api/v1/tasks/{id}/cancel
GET  /api/v1/events       # Server-Sent Events
POST /api/v1/shutdown
```

`GET /api/v1/health` 不需要 Token，其余端点需要 Token。认证失败返回 `401`，不支持的方法返回 `405`，业务错误返回对应状态码和上述错误 JSON。

`GET /api/v1/translators` 返回引擎的 `id`、`name`、`valid` 和可用语言 `languages`。每个语言包含 `id` 和 `name`。客户端应使用返回的 `id`，并优先选择 `valid=true` 的引擎；选择引擎后应使用其语言列表校验源语言和目标语言。

`POST /api/v1/translators/{id}/test-connection` 测试引擎接入点的网络连通性。请求可传 `{"region":"cn-hangzhou"}`（阿里云）或 `{"region":"ap-guangzhou"}`（腾讯云）来测试设置页当前选中的地区；百度云、彩云小译、DeepL Pro、有道智云和 DeepSeek 无需地区参数。DeepSeek 测试配置的服务地址，HTTPS 地址进行 TLS 握手，本机 HTTP 地址进行 TCP 连接。成功返回接入点和耗时，例如 `{"endpoint":"api.deepseek.com","latency_ms":123}`。此检查不调用翻译 API，不验证密钥，也不消耗翻译额度。地区无效返回 `400`，网络或 TLS 失败返回 `502`，引擎不存在返回 `404`。该端点同样要求 `X-Anto-Token`。

`GET /api/v1/config` 返回稳定的应用和供应商配置 DTO，不暴露 Viper 的任意嵌套结构：

```json
{
  "app": {
    "default_translator":"ali_cloud_mt",
    "default_from_lang":"en",
    "default_to_lang":"zh",
    "output_directory":"",
    "track_export":2,
    "keep_original":true,
    "notify_on_complete":true,
    "auto_connect":true,
    "minimize_to_tray":false,
    "log_level":"信息",
    "concurrency":5
  },
  "providers": {
    "ali_cloud_mt": {
      "configured":false,
      "values":{"ak_id":"","ak_secret":"","region":"cn-hangzhou","qps":50,"max_single_text_length":3000,"max_coroutine_num":25}
    },
    "tencent_cloud_mt": {
      "configured":false,
      "values":{"secret_id":"","secret_key":"","region":"ap-guangzhou","qps":5,"max_single_text_length":2000,"max_coroutine_num":5}
    },
    "baidu_cloud_mt": {
      "configured":false,
      "values":{"app_id":"","app_secret":"","service_tier":"standard","max_coroutine_num":1}
    },
    "caiyunai": {
      "configured":false,
      "values":{"token":"","qps":1,"max_single_text_length":1000,"max_coroutine_num":1}
    },
    "deepl_pro": {
      "configured":false,
      "values":{"api_key":"","qps":5,"max_single_text_length":3000,"max_coroutine_num":5}
    },
    "deepseek": {
      "configured":false,
      "values":{"api_key":"","base_url":"https://api.deepseek.com","model":"deepseek-flash","prompt":"请忠实翻译字幕，保留口吻、专有名词和原有换行，不要添加说明。","qps":1,"max_single_text_length":2000,"max_coroutine_num":1}
    }
  }
}
```

供应商 `id` 包括 `ali_cloud_mt`（机器翻译-阿里云）、`tencent_cloud_mt`（机器翻译-腾讯云）、`baidu_cloud_mt`（机器翻译-百度云）、`caiyunai`（彩云小译）、`deepl_pro`（DeepL Pro）、`deepseek`（AI 翻译-DeepSeek）和 `youdao_batch`（批量翻译-有道智云）。敏感值读取时返回 `********`；更新时原样提交该占位值表示保留当前凭据。有道智云使用 `app_key`（应用 ID）、`app_secret`（应用密钥）、`qps`、`max_single_text_length`（不超过 5000）和 `max_coroutine_num`。Core 会用当前配置默认值补齐缺失字段，因此客户端可以直接按返回的 `values` 渲染表单。应用设置通过 `{"app":{...}}` 提交，供应商配置通过 `{"providers":{"deepseek":{"api_key":"...","model":"deepseek-flash"}}}` 提交。`app` 可以只包含本次修改的字段，省略的字段保留当前值。

应用设置中的 `track_export` 只能是 `0`、`1` 或 `2`，`concurrency` 必须在 `1` 到 `256` 之间。供应商的 `qps`、`max_single_text_length` 和 `max_coroutine_num` 必须是正整数，分别不能超过 `100000`、`1000000` 和 `1024`；腾讯云的单次最大长度不能超过 `2000` 个字符。百度云用 `service_tier` 选择 `standard`、`advanced` 或 `premium`，分别限制为 `1000` 字符/1 QPS、`6000` 字符/10 QPS、`6000` 字符/100 QPS；每次请求还须小于等于 `6000` 个 UTF-8 字节。百度云不接受单独设置 `qps` 或 `max_single_text_length`。校验失败返回 `400` 和 `invalid_config`；配置文件不可写返回 `500` 和 `config_write_failed`。

彩云小译使用独立申请的 API Token。`qps:1` 和单次 `1000` 字符是本应用的保守默认值，并非服务商公布的上限。源语言为中文或繁体中文时，目标可选另一种中文、英语、日语、韩语；英语、日语、韩语、德语、西班牙语、法语、意大利语、葡萄牙语、俄语、土耳其语、越南语只能译为中文或繁体中文。不支持的方向在任务创建时被拒绝。

DeepL Pro 使用 `api.deepl.com/v2/translate` 和 `DeepL-Auth-Key` 请求头，适用于已有 Pro 或当前付费 API 套餐的密钥，不适用于 API Free 端点的密钥。默认 `qps:5`、单次 `3000` 字符是应用本地限速与分段值，不是 DeepL 的固定上限。Core 保证请求体不超过 DeepL 的 `128 KiB` 限制，并对 `429`、`529`、`500`、`503`、`504` 做有限次退避重试；`456` 额度不足直接报错。英语目标默认美式英语，葡萄牙语目标默认葡萄牙葡语。

DeepSeek 默认使用 `https://api.deepseek.com/chat/completions` 和 `deepseek-flash`。`base_url` 可指向自建、兼容 Chat Completions 的服务；可填写服务根地址（例如 `http://127.0.0.1:8000/v1`）或完整的 `/chat/completions` 地址。远端地址必须是 HTTPS，本机回环地址可用 HTTP。官方地址必须设置 `api_key`；自建服务可留空。`model` 和 `prompt` 可由用户调整；`prompt` 是翻译风格要求，Core 固定附加 JSON 数组的条数、顺序和输出格式约束。Core 会按 `max_single_text_length`（不超过 6000 个字符）分段与分组，核对模型输出条数后写入字幕；对 `429`、`500`、`503` 有限次重试。官方地址使用 JSON 输出模式，自建兼容服务只需返回含 JSON 对象的标准 Chat Completions 响应。保存配置后，客户端重新加载 Core 才会启用新设置。

创建任务的请求示例：

```json
{
  "files": ["C:/video/source.srt"],
  "translator_id": "ali_cloud_mt",
  "from_lang": "en",
  "to_lang": "zh",
  "translate_mode": "全量翻译",
  "main_track_report": "来源语种",
  "flag_track_export": 2,
  "output_directory": ""
}
```

`files` 和 `directory` 至少提供一个。目录只扫描直接子项并忽略非 `.srt` 文件。未提供 `output_directory` 时，每个源文件写入源文件目录下的 `anto/`。

任务快照包含 `id`、`status`、`total`、`completed`、`current_file`、`output_directory`、`output_files`、`error`、`created_at` 和 `updated_at`。状态为 `queued`、`reading`、`translating`、`writing`、`completed`、`failed` 或 `cancelled`。Core 在每个翻译块调用供应商前发布 `translating`，写文件前发布 `writing`。

任务事件使用 JSON 作为 SSE 的 `data`，事件类型包括 `tasks.snapshot`、`task.created`、`task.progress`、`task.reading`、`task.translating`、`task.writing`、`task.completed`、`task.failed` 和 `task.cancelled`。

客户端首次连接或重连时，服务端先发送当前任务快照：

```text
event: tasks.snapshot
data: {"type":"tasks.snapshot","tasks":[]}
```

增量事件携带 `task_id` 和最新 `task` 快照。Core 为慢客户端保留有限缓冲，缓冲满时丢弃最旧事件；客户端应在重连后用快照校准列表。
