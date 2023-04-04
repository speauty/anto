package openapi_youdao

type Cfg struct {
	AppKey    string `mapstructure:"app_key"`    // 应用ID
	AppSecret string `mapstructure:"app_secret"` // 应用密钥
}

func (customC Cfg) Default() *Cfg {
	return &Cfg{AppSecret: "", AppKey: ""}
}