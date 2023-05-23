package ling_va

import (
	"anto/domain/service/translator"
	"fmt"
	"github.com/spf13/viper"
)

type Cfg struct {
	*translator.DefaultConfig
	DataId              string `mapstructure:"data_id"`
	MaxSingleTextLength int    `mapstructure:"max_single_text_length"`
	QPS                 int    `mapstructure:"qps"`
	MaxCoroutineNum     int    `mapstructure:"max_coroutine_num"`
}

func (customC *Cfg) GetAK() string { return customC.DataId }

func (customC *Cfg) GetMaxSingleTextLength() int { return customC.MaxSingleTextLength }
func (customC *Cfg) GetQPS() int                 { return customC.QPS }
func (customC *Cfg) GetMaxCoroutineNum() int     { return customC.MaxCoroutineNum }

func (customC *Cfg) Default() translator.ImplConfig {
	return &Cfg{
		DataId:              "3qnDcUVykFKnSC3cdRX2t",
		MaxSingleTextLength: 1000, QPS: 10, MaxCoroutineNum: 20,
	}
}

func (customC *Cfg) SetAK(ak string) error {
	customC.DataId = ak
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

	viper.Set("ling_va.data_id", customC.DataId)
	viper.Set("ling_va.max_single_text_length", customC.MaxSingleTextLength)
	viper.Set("ling_va.qps", customC.QPS)
	viper.Set("ling_va.max_coroutine_num", customC.MaxCoroutineNum)

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("写入配置[%s]失败, 错误: %s", Singleton().GetName(), err)
	}
	return nil
}
