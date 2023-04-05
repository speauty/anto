package caiyunai

type Cfg struct {
	Token string `mapstructure:"token"`
}

func (customC Cfg) Default() *Cfg {
	// 这是官网给出的测试Token, 随便造
	return &Cfg{Token: "3975l6lr5pcbvidl6jl2"}
}
