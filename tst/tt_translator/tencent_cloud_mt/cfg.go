package tencent_cloud_mt

type Cfg struct {
	SecretId            string `mapstructure:"secret_id"`              // 用于标识接口调用者身份
	SecretKey           string `mapstructure:"secret_key"`             // 用于验证接口调用者的身份
	Region              string `mapstructure:"-"`                      // 地域参数
	ProjectId           int64  `mapstructure:"-"`                      // 项目ID, 默认值为0
	MaxSingleTextLength int    `mapstructure:"max_single_text_length"` // 单次翻译最大长度
}

func (customC Cfg) Default() *Cfg {
	return &Cfg{
		SecretId:            "",
		SecretKey:           "",
		Region:              "ap-chengdu",
		ProjectId:           0,
		MaxSingleTextLength: 2000,
	}
}
