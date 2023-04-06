package ling_va

type Cfg struct {
	DataId              string `mapstructure:"data_id"`
	MaxSingleTextLength int    `mapstructure:"max_single_text_length"`
}

func (customC Cfg) Default() *Cfg {
	return &Cfg{
		DataId:              "3qnDcUVykFKnSC3cdRX2t",
		MaxSingleTextLength: 1000,
	}
}
