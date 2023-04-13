package openapi_youdao

type Cfg struct {
	AppKey              string `mapstructure:"app_key"`                // 应用ID
	AppSecret           string `mapstructure:"app_secret"`             // 应用密钥
	MaxSingleTextLength int    `mapstructure:"max_single_text_length"` // 单次翻译最大长度
}

func (customC Cfg) Default() *Cfg {
	return &Cfg{AppSecret: "", AppKey: "", MaxSingleTextLength: 5000}
}
