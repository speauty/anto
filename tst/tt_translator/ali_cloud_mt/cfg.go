package ali_cloud_mt

type Cfg struct {
	AKId     string `mapstructure:"ak_id"`
	AKSecret string `mapstructure:"ak_secret"`
	Region   string `mapstructure:"region"`
}

func (customC Cfg) Default() *Cfg {
	return &Cfg{
		AKId: "", AKSecret: "", Region: "",
	}
}
