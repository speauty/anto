package baidu

import (
	"anto/domain/service/translator"
	"fmt"
	"github.com/spf13/viper"
)

type Cfg struct {
	*translator.DefaultConfig
	AppId               string `mapstructure:"app_id"`  // 应用ID
	AppKey              string `mapstructure:"app_key"` // 应用密钥
	MaxSingleTextLength int    `mapstructure:"max_single_text_length"`
	QPS                 int    `mapstructure:"qps"`
	MaxCoroutineNum     int    `mapstructure:"max_coroutine_num"`
}

func (customC *Cfg) GetAK() string               { return customC.AppId }
func (customC *Cfg) GetSK() string               { return customC.AppKey }
func (customC *Cfg) GetMaxSingleTextLength() int { return customC.MaxSingleTextLength }
func (customC *Cfg) GetQPS() int                 { return customC.QPS }
func (customC *Cfg) GetMaxCoroutineNum() int     { return customC.MaxCoroutineNum }

func (customC *Cfg) Default() translator.ImplConfig {
	return &Cfg{
		AppId: "", AppKey: "",
		MaxSingleTextLength: 1000, QPS: 1, MaxCoroutineNum: 1,
	}
}

func (customC *Cfg) SetAK(ak string) error {
	customC.AppId = ak
	return nil
}

func (customC *Cfg) SetSK(sk string) error {
	customC.AppKey = sk
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

	viper.Set("baidu.app_id", customC.AppId)
	viper.Set("baidu.app_key", customC.AppKey)
	viper.Set("baidu.max_single_text_length", customC.MaxSingleTextLength)
	viper.Set("baidu.qps", customC.QPS)
	viper.Set("baidu.max_coroutine_num", customC.MaxCoroutineNum)

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("写入配置[%s]失败, 错误: %s", Singleton().GetName(), err)
	}
	return nil
}
