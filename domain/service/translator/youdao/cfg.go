package youdao

import (
	"anto/domain/service/translator"
	"fmt"
	"github.com/spf13/viper"
)

type Cfg struct {
	*translator.DefaultConfig
	MaxSingleTextLength int `mapstructure:"max_single_text_length"` // 单次翻译最大长度
	QPS                 int `mapstructure:"qps"`
	MaxCoroutineNum     int `mapstructure:"max_coroutine_num"`
}

func (customC *Cfg) GetMaxSingleTextLength() int { return customC.MaxSingleTextLength }
func (customC *Cfg) GetQPS() int                 { return customC.QPS }
func (customC *Cfg) GetMaxCoroutineNum() int     { return customC.MaxCoroutineNum }

func (customC *Cfg) Default() translator.ImplConfig {
	return &Cfg{
		MaxSingleTextLength: 5000, QPS: 50, MaxCoroutineNum: 20,
	}
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

	viper.Set("youdao.max_single_text_length", customC.MaxSingleTextLength)
	viper.Set("youdao.qps", customC.QPS)
	viper.Set("youdao.max_coroutine_num", customC.MaxCoroutineNum)

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("写入配置[%s]失败, 错误: %s", Singleton().GetName(), err)
	}
	return nil
}
