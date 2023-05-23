package ali_cloud_mt

import (
	"anto/domain/service/translator"
	"fmt"
	"github.com/spf13/viper"
)

type Cfg struct {
	*translator.DefaultConfig `mapstructure:"-"`
	AKId                      string `mapstructure:"ak_id"`
	AKSecret                  string `mapstructure:"ak_secret"`
	Region                    string `mapstructure:"region"`
	MaxSingleTextLength       int    `mapstructure:"max_single_text_length"`
	QPS                       int    `mapstructure:"qps"`
	MaxCoroutineNum           int    `mapstructure:"max_coroutine_num"`
}

func (customC *Cfg) GetAK() string               { return customC.AKId }
func (customC *Cfg) GetSK() string               { return customC.AKSecret }
func (customC *Cfg) GetMaxSingleTextLength() int { return customC.MaxSingleTextLength }
func (customC *Cfg) GetQPS() int                 { return customC.QPS }
func (customC *Cfg) GetMaxCoroutineNum() int     { return customC.MaxCoroutineNum }
func (customC *Cfg) GetRegion() string           { return customC.Region }

func (customC *Cfg) Default() translator.ImplConfig {
	return &Cfg{
		AKId: "", AKSecret: "", Region: "",
		MaxSingleTextLength: 3000, QPS: 50, MaxCoroutineNum: 20,
	}
}

func (customC *Cfg) SetAK(ak string) error {
	customC.AKId = ak
	return nil
}

func (customC *Cfg) SetSK(sk string) error {
	customC.AKSecret = sk
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

func (customC *Cfg) SetRegion(region string) error {
	customC.Region = region
	return nil
}

func (customC *Cfg) Sync() error {

	viper.Set("ali_cloud_mt.ak_id", customC.AKId)
	viper.Set("ali_cloud_mt.ak_secret", customC.AKSecret)
	viper.Set("ali_cloud_mt.region", customC.Region)
	viper.Set("ali_cloud_mt.max_single_text_length", customC.MaxSingleTextLength)
	viper.Set("ali_cloud_mt.qps", customC.QPS)
	viper.Set("ali_cloud_mt.max_coroutine_num", customC.MaxCoroutineNum)

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("写入配置[%s]失败, 错误: %s", Singleton().GetName(), err)
	}
	return nil
}
