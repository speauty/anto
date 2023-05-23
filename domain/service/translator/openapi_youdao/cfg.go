package openapi_youdao

import (
	"anto/domain/service/translator"
	"fmt"
	"github.com/spf13/viper"
)

type Cfg struct {
	*translator.DefaultConfig
	AppKey              string `mapstructure:"app_key"`                // 应用ID
	AppSecret           string `mapstructure:"app_secret"`             // 应用密钥
	MaxSingleTextLength int    `mapstructure:"max_single_text_length"` // 单次翻译最大长度
	QPS                 int    `mapstructure:"qps"`
	MaxCoroutineNum     int    `mapstructure:"max_coroutine_num"`
}

func (customC *Cfg) GetAK() string               { return customC.AppKey }
func (customC *Cfg) GetSK() string               { return customC.AppSecret }
func (customC *Cfg) GetMaxSingleTextLength() int { return customC.MaxSingleTextLength }
func (customC *Cfg) GetQPS() int                 { return customC.QPS }
func (customC *Cfg) GetMaxCoroutineNum() int     { return customC.MaxCoroutineNum }

func (customC *Cfg) Default() translator.ImplConfig {
	return &Cfg{
		AppSecret: "", AppKey: "",
		MaxSingleTextLength: 5000, QPS: 1, MaxCoroutineNum: 20,
	}
}

func (customC *Cfg) SetAK(ak string) error {
	customC.AppKey = ak
	return nil
}

func (customC *Cfg) SetSK(sk string) error {
	customC.AppSecret = sk
	return nil
}

func (customC *Cfg) SetMaxSingleTextLength(textLen int) error {
	customC.MaxSingleTextLength = textLen
	return nil
}

func (customC *Cfg) SetQPS(qps int) error {
	customC.QPS = qps
	return nil
}

func (customC *Cfg) SetMaxCoroutineNum(coroutineNum int) error {
	customC.MaxCoroutineNum = coroutineNum
	return nil
}

func (customC *Cfg) Sync() error {

	viper.Set("openapi_youdao.app_key", customC.AppKey)
	viper.Set("openapi_youdao.app_secret", customC.AppSecret)
	viper.Set("openapi_youdao.max_single_text_length", customC.MaxSingleTextLength)
	viper.Set("openapi_youdao.qps", customC.QPS)
	viper.Set("openapi_youdao.max_coroutine_num", customC.MaxCoroutineNum)

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("写入配置[%s]失败, 错误: %s", Singleton().GetName(), err)
	}
	return nil
}
