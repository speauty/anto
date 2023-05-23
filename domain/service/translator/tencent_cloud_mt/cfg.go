package tencent_cloud_mt

import (
	"anto/domain/service/translator"
	"fmt"
	"github.com/spf13/viper"
)

type Cfg struct {
	*translator.DefaultConfig
	SecretId            string `mapstructure:"secret_id"`              // 用于标识接口调用者身份
	SecretKey           string `mapstructure:"secret_key"`             // 用于验证接口调用者的身份
	Region              string `mapstructure:"-"`                      // 地域参数
	ProjectId           int64  `mapstructure:"-"`                      // 项目ID, 默认值为0
	MaxSingleTextLength int    `mapstructure:"max_single_text_length"` // 单次翻译最大长度
	QPS                 int    `mapstructure:"qps"`
	MaxCoroutineNum     int    `mapstructure:"max_coroutine_num"`
}

func (customC *Cfg) GetAK() string               { return customC.SecretId }
func (customC *Cfg) GetSK() string               { return customC.SecretKey }
func (customC *Cfg) GetMaxSingleTextLength() int { return customC.MaxSingleTextLength }
func (customC *Cfg) GetQPS() int                 { return customC.QPS }
func (customC *Cfg) GetMaxCoroutineNum() int     { return customC.MaxCoroutineNum }
func (customC *Cfg) GetRegion() string           { return customC.Region }

func (customC *Cfg) Default() translator.ImplConfig {
	return &Cfg{
		SecretId: "", SecretKey: "", Region: "ap-chengdu", ProjectId: 0,
		MaxSingleTextLength: 2000, QPS: 5, MaxCoroutineNum: 20,
	}
}

func (customC *Cfg) SetAK(ak string) error {
	customC.SecretId = ak
	return nil
}

func (customC *Cfg) SetSK(sk string) error {
	customC.SecretKey = sk
	return nil
}

func (customC *Cfg) SetRegion(region string) error {
	customC.Region = region
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

	viper.Set("tencent_cloud_mt.secret_id", customC.SecretId)
	viper.Set("tencent_cloud_mt.secret_key", customC.SecretKey)
	viper.Set("tencent_cloud_mt.max_single_text_length", customC.MaxSingleTextLength)
	viper.Set("tencent_cloud_mt.qps", customC.QPS)
	viper.Set("tencent_cloud_mt.max_coroutine_num", customC.MaxCoroutineNum)

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("写入配置[%s]失败, 错误: %s", Singleton().GetName(), err)
	}
	return nil
}
