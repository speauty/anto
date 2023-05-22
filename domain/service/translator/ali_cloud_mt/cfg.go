package ali_cloud_mt

import "anto/domain/service/translator"

type Cfg struct {
	*translator.DefaultConfig
	AKId                string `mapstructure:"ak_id"`
	AKSecret            string `mapstructure:"ak_secret"`
	Region              string `mapstructure:"region"`
	MaxSingleTextLength int    `mapstructure:"max_single_text_length"`
}

func (customC *Cfg) GetAK() string     { return customC.AKId }
func (customC *Cfg) GetSK() string     { return customC.AKSecret }
func (customC *Cfg) GetTML() int       { return customC.MaxSingleTextLength }
func (customC *Cfg) GetRegion() string { return customC.Region }

func (customC *Cfg) Default() translator.ImplConfig {
	return &Cfg{
		AKId: "", AKSecret: "", Region: "", MaxSingleTextLength: 3000,
	}
}

func (customC *Cfg) SetAK(ak string) error {
	customC.AKId = ak
	return nil
}

func (customC *Cfg) SetSK(sk string) error {
	customC.AKSecret = sk
	return nil
}

func (customC *Cfg) SetRegion(region string) error {
	customC.Region = region
	return nil
}

func (customC *Cfg) SetTML(tml int) error {
	customC.MaxSingleTextLength = tml
	return nil
}
