package volcengine

import (
	"anto/domain/service/translator"
	"fmt"
	"github.com/spf13/viper"
)

type Cfg struct {
	*translator.DefaultConfig
	AccessKey           string `mapstructure:"access_key"`
	SecretKey           string `mapstructure:"secret_key"`
	MaxSingleTextLength int    `mapstructure:"max_single_text_length"` // 单次翻译最大长度
	QPS                 int    `mapstructure:"qps"`
	MaxCoroutineNum     int    `mapstructure:"max_coroutine_num"`
}

func (customC *Cfg) GetAK() string               { return customC.AccessKey }
func (customC *Cfg) GetSK() string               { return customC.SecretKey }
func (customC *Cfg) GetMaxSingleTextLength() int { return customC.MaxSingleTextLength }
func (customC *Cfg) GetQPS() int                 { return customC.QPS }
func (customC *Cfg) GetMaxCoroutineNum() int     { return customC.MaxCoroutineNum }

func (customC *Cfg) Default() translator.ImplConfig {
	return &Cfg{
		AccessKey: "", SecretKey: "",
		MaxSingleTextLength: 5000, QPS: 10, MaxCoroutineNum: 20,
	}
}

func (customC *Cfg) SetAK(ak string) error {
	customC.AccessKey = ak
	return nil
}

func (customC *Cfg) SetSK(sk string) error {
	customC.SecretKey = sk
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

	viper.Set("volc_engine.access_key", customC.AccessKey)
	viper.Set("volc_engine.secret_key", customC.SecretKey)
	viper.Set("volc_engine.max_single_text_length", customC.MaxSingleTextLength)
	viper.Set("volc_engine.qps", customC.QPS)
	viper.Set("volc_engine.max_coroutine_num", customC.MaxCoroutineNum)

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("写入配置[%s]失败, 错误: %s", Singleton().GetName(), err)
	}
	return nil
}
