package volcengine

type Cfg struct {
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
}

func (customC Cfg) Default() *Cfg {
	return &Cfg{AccessKey: "", SecretKey: ""}
}
