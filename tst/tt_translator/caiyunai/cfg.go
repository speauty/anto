package caiyunai

type Cfg struct {
	Token               string `mapstructure:"token"`
	MaxSingleTextLength int    `mapstructure:"max_single_text_length"`
}

func (customC Cfg) Default() *Cfg {
	// 这是官网给出的测试Token, 随便造
	return &Cfg{Token: "3975l6lr5pcbvidl6jl2", MaxSingleTextLength: 5000}
}
