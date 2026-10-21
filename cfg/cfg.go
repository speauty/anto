package cfg

import (
	"anto/domain/service/translator/ali_cloud_mt"
	"anto/domain/service/translator/baidu_cloud_mt"
	"anto/domain/service/translator/caiyunai"
	"anto/domain/service/translator/deepl_pro"
	"anto/domain/service/translator/deepseek"
	"anto/domain/service/translator/tencent_cloud_mt"
	"anto/domain/service/translator/youdao_batch"
	"github.com/spf13/viper"
	"sync"
)

var (
	apiSingleton  *Cfg
	onceSingleton sync.Once
)

func Singleton() *Cfg {
	onceSingleton.Do(func() {
		apiSingleton = new(Cfg)
		apiSingleton.AliCloudMT = new(ali_cloud_mt.Config).Default().(*ali_cloud_mt.Config)
		apiSingleton.TencentCloudMT = new(tencent_cloud_mt.Config).Default().(*tencent_cloud_mt.Config)
		apiSingleton.BaiduCloudMT = new(baidu_cloud_mt.Config).Default().(*baidu_cloud_mt.Config)
		apiSingleton.CaiyunAI = new(caiyunai.Config).Default().(*caiyunai.Config)
		apiSingleton.DeepLPro = new(deepl_pro.Config).Default().(*deepl_pro.Config)
		apiSingleton.DeepSeek = new(deepseek.Config).Default().(*deepseek.Config)
		apiSingleton.YoudaoBatch = new(youdao_batch.Config).Default().(*youdao_batch.Config)
	})
	return apiSingleton
}

type Cfg struct {
	AliCloudMT     *ali_cloud_mt.Config     `mapstructure:"ali_cloud_mt"`
	TencentCloudMT *tencent_cloud_mt.Config `mapstructure:"tencent_cloud_mt"`
	BaiduCloudMT   *baidu_cloud_mt.Config   `mapstructure:"baidu_cloud_mt"`
	CaiyunAI       *caiyunai.Config         `mapstructure:"caiyunai"`
	DeepLPro       *deepl_pro.Config        `mapstructure:"deepl_pro"`
	DeepSeek       *deepseek.Config         `mapstructure:"deepseek"`
	YoudaoBatch    *youdao_batch.Config     `mapstructure:"youdao_batch"`

	currentViper *viper.Viper `mapstructure:"-"`
}
