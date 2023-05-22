package niutrans

type Cfg struct {
	AppKey string `mapstructure:"app_key"`
}

func (customC Cfg) Default() *Cfg {
	return &Cfg{
		AppKey: "",
	}
}
