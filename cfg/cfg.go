package cfg

import (
	"anto/domain/service/translator/ali_cloud_mt"
	"anto/domain/service/translator/baidu"
	"anto/domain/service/translator/caiyunai"
	"anto/domain/service/translator/huawei_cloud_nlp"
	"anto/domain/service/translator/ling_va"
	"anto/domain/service/translator/niutrans"
	"anto/domain/service/translator/openapi_youdao"
	"anto/domain/service/translator/tencent_cloud_mt"
	"anto/domain/service/translator/volcengine"
	"anto/domain/service/translator/youdao"
	"fmt"
	"sync"
)

var (
	apiSingleton  *Cfg
	onceSingleton sync.Once
)

func Singleton() *Cfg {
	onceSingleton.Do(func() {
		apiSingleton = new(Cfg)
		apiSingleton.App = App{}.Default()
		apiSingleton.HuaweiCloudNlp = new(huawei_cloud_nlp.Cfg).Default().(*huawei_cloud_nlp.Cfg)
		apiSingleton.LingVA = new(ling_va.Cfg).Default().(*ling_va.Cfg)
		apiSingleton.Baidu = new(baidu.Cfg).Default().(*baidu.Cfg)
		apiSingleton.TencentCloudMT = new(tencent_cloud_mt.Cfg).Default().(*tencent_cloud_mt.Cfg)
		apiSingleton.OpenAPIYouDao = new(openapi_youdao.Cfg).Default().(*openapi_youdao.Cfg)
		apiSingleton.AliCloudMT = new(ali_cloud_mt.Cfg).Default().(*ali_cloud_mt.Cfg)
		apiSingleton.CaiYunAI = new(caiyunai.Cfg).Default().(*caiyunai.Cfg)
		apiSingleton.Niutrans = new(niutrans.Cfg).Default().(*niutrans.Cfg)
		apiSingleton.VolcEngine = new(volcengine.Cfg).Default().(*volcengine.Cfg)
		apiSingleton.YouDao = new(youdao.Cfg).Default().(*youdao.Cfg)
	})
	return apiSingleton
}

type Cfg struct {
	App            *App                  `mapstructure:"-"`
	HuaweiCloudNlp *huawei_cloud_nlp.Cfg `mapstructure:"huawei_cloud_nlp"`
	LingVA         *ling_va.Cfg          `mapstructure:"ling_va"`
	Baidu          *baidu.Cfg            `mapstructure:"baidu"`
	TencentCloudMT *tencent_cloud_mt.Cfg `mapstructure:"tencent_cloud_mt"`
	OpenAPIYouDao  *openapi_youdao.Cfg   `mapstructure:"openapi_youdao"`
	AliCloudMT     *ali_cloud_mt.Cfg     `mapstructure:"ali_cloud_mt"`
	CaiYunAI       *caiyunai.Cfg         `mapstructure:"caiyun_ai"`
	Niutrans       *niutrans.Cfg         `mapstructure:"niutrans"`
	VolcEngine     *volcengine.Cfg       `mapstructure:"volc_engine"`
	YouDao         *youdao.Cfg           `mapstructure:"youdao"`
}

func (customC *Cfg) NewUITitle() string {
	return fmt.Sprintf(
		"%s-%s(作者: %s 邮箱: speauty@163.com)",
		"aa", customC.App.Version, customC.App.Author,
	)
}
