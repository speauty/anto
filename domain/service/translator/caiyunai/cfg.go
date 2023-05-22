package caiyunai

import "anto/domain/service/translator"

type Cfg struct {
	*translator.DefaultConfig
	Token               string `mapstructure:"token"`
	MaxSingleTextLength int    `mapstructure:"max_single_text_length"`
}

func (customC *Cfg) GetAK() string { return customC.Token }

func (customC *Cfg) Default() translator.ImplConfig {
	// 这是官网给出的测试Token, 随便造
	return &Cfg{Token: "3975l6lr5pcbvidl6jl2", MaxSingleTextLength: 5000}
}

func (customC *Cfg) SetAK(ak string) error {
	customC.Token = ak
	return nil
}
