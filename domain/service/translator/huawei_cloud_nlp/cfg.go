package huawei_cloud_nlp

import (
	"anto/domain/service/translator"
	"fmt"
	"github.com/spf13/viper"
)

type Cfg struct {
	*translator.DefaultConfig
	AKId                string `mapstructure:"ak_id"`      // Access Key
	SkKey               string `mapstructure:"sk_key"`     // Secret Access Key
	Region              string `mapstructure:"region"`     // 当前接口开发的区域, 目前仅支持华北-北京四终端节点
	ProjectId           string `mapstructure:"project_id"` // 项目ID
	MaxSingleTextLength int    `mapstructure:"max_single_text_length"`
	QPS                 int    `mapstructure:"qps"`
	MaxCoroutineNum     int    `mapstructure:"max_coroutine_num"`
}

func (customC *Cfg) GetAK() string               { return customC.AKId }
func (customC *Cfg) GetSK() string               { return customC.SkKey }
func (customC *Cfg) GetPK() string               { return customC.ProjectId }
func (customC *Cfg) GetRegion() string           { return customC.Region }
func (customC *Cfg) GetMaxSingleTextLength() int { return customC.MaxSingleTextLength }
func (customC *Cfg) GetQPS() int                 { return customC.QPS }
func (customC *Cfg) GetMaxCoroutineNum() int     { return customC.MaxCoroutineNum }

func (customC *Cfg) Default() translator.ImplConfig {
	return &Cfg{
		AKId: "", SkKey: "", Region: "cn-north-4", ProjectId: "",
		MaxSingleTextLength: 2000, QPS: 20, MaxCoroutineNum: 10,
	}
}

func (customC *Cfg) SetAK(ak string) error {
	customC.AKId = ak
	return nil
}

func (customC *Cfg) SetSK(sk string) error {
	customC.SkKey = sk
	return nil
}

func (customC *Cfg) SetPK(pk string) error {
	customC.ProjectId = pk
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

	viper.Set("huawei_cloud_nlp.ak_id", customC.AKId)
	viper.Set("huawei_cloud_nlp.sk_key", customC.SkKey)
	viper.Set("huawei_cloud_nlp.project_id", customC.ProjectId)
	viper.Set("huawei_cloud_nlp.region", customC.Region)
	viper.Set("huawei_cloud_nlp.max_single_text_length", customC.MaxSingleTextLength)
	viper.Set("huawei_cloud_nlp.qps", customC.QPS)
	viper.Set("huawei_cloud_nlp.max_coroutine_num", customC.MaxCoroutineNum)

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("写入配置[%s]失败, 错误: %s", Singleton().GetName(), err)
	}
	return nil
}
