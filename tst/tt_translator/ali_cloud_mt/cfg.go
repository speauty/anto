package ali_cloud_mt

type Cfg struct {
	AKId                string `mapstructure:"ak_id"`
	AKSecret            string `mapstructure:"ak_secret"`
	Region              string `mapstructure:"region"`
	MaxSingleTextLength int    `mapstructure:"max_single_text_length"`
}

func (customC Cfg) Default() *Cfg {
	return &Cfg{
		AKId: "", AKSecret: "", Region: "", MaxSingleTextLength: 3000,
	}
}
