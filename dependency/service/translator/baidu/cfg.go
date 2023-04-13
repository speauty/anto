package baidu

type Cfg struct {
	AppId               string `mapstructure:"app_id"`  // 应用ID
	AppKey              string `mapstructure:"app_key"` // 应用密钥
	MaxSingleTextLength int    `mapstructure:"max_single_text_length"`
}

func (customC Cfg) Default() *Cfg {
	return &Cfg{AppId: "", AppKey: "", MaxSingleTextLength: 1000}
}
