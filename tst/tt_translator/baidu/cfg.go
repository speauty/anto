package baidu

type Cfg struct {
	AppId  string `mapstructure:"app_id"`  // 应用ID
	AppKey string `mapstructure:"app_key"` // 应用密钥
}

func (customC Cfg) Default() *Cfg {
	return &Cfg{AppId: "", AppKey: ""}
}
