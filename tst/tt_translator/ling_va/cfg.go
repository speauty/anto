package ling_va

type Cfg struct {
	DataId string `mapstructure:"data_id"`
}

func (customC Cfg) DefaultCfg() *Cfg {
	return &Cfg{
		DataId: "3qnDcUVykFKnSC3cdRX2t",
	}
}
