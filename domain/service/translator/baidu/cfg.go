package baidu

import "anto/domain/service/translator"

type Cfg struct {
	*translator.DefaultConfig
	AppId               string `mapstructure:"app_id"`  // 应用ID
	AppKey              string `mapstructure:"app_key"` // 应用密钥
	MaxSingleTextLength int    `mapstructure:"max_single_text_length"`
}

func (customC *Cfg) GetAK() string { return customC.AppId }
func (customC *Cfg) GetSK() string { return customC.AppKey }
func (customC *Cfg) GetTML() int   { return customC.MaxSingleTextLength }

func (customC *Cfg) Default() translator.ImplConfig {
	return &Cfg{AppId: "", AppKey: "", MaxSingleTextLength: 1000}
}

func (customC *Cfg) SetAK(ak string) error {
	customC.AppId = ak
	return nil
}

func (customC *Cfg) SetSK(sk string) error {
	customC.AppKey = sk
	return nil
}

func (customC *Cfg) SetTML(tml int) error {
	customC.MaxSingleTextLength = tml
	return nil
}
