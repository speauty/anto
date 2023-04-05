package _const

const (
	Author     = "speauty"
	Version    = "v2.5.0"
	EnvDebug   = "debug"
	EnvRelease = "release"

	GoUidLen = 8

	UITitle       = "Anto"
	UIIcon        = "favicon.ico"
	UIResourceDir = "./assets"

	CfgYmlDefaultContent = `
huawei_cloud_nlp: # 华为云机器翻译
  ak_id: ""             # Access Key
  sk_key: ""            # Secret Access Key
  region: "cn-north-4"  # 当前接口开发的区域, 目前仅支持华北-北京四终端节点
  project_id: ""        # 项目ID

ling_va:
  data_id: "3qnDcUVykFKnSC3cdRX2t"	# 数据ID

baidu:
  app_id: ""            # 应用ID
  app_key: ""           # 应用密钥

tencent_cloud_mt:       # 腾讯云机器翻译
  secret_id: ""         # 密钥ID 用于标识接口调用者身份
  secret_key: ""        # 密钥关键字 用于验证接口调用者的身份

openapi_youdao:         # 有道智云翻译
  app_key: ""           # 应用ID
  app_secret: ""        # 应用密钥

ali_cloud_mt:           # 阿里云翻译
  ak_id: ""             # 应用ID
  ak_secret: ""         # 应用密钥
  region: ""            # 区域
caiyun_ai:                      # 彩云小翻译
  token: "3975l6lr5pcbvidl6jl2" # 密钥
`
)
