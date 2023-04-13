package cfg

import (
	"anto/dependency/service/translator/ali_cloud_mt"
	"anto/dependency/service/translator/baidu"
	"anto/dependency/service/translator/caiyunai"
	"anto/dependency/service/translator/huawei_cloud_nlp"
	"anto/dependency/service/translator/ling_va"
	"anto/dependency/service/translator/openapi_youdao"
	"anto/dependency/service/translator/tencent_cloud_mt"
	"anto/lib/ui"
	"fmt"
	"sync"
)

var (
	apiCfg  *Cfg
	onceCfg sync.Once
)

func GetInstance() *Cfg {
	onceCfg.Do(func() {
		apiCfg = new(Cfg)
		apiCfg.App = App{}.Default()
		apiCfg.UI = ui.Cfg{}.Default()
		apiCfg.HuaweiCloudNlp = huawei_cloud_nlp.Cfg{}.Default()
		apiCfg.LingVA = ling_va.Cfg{}.Default()
		apiCfg.Baidu = baidu.Cfg{}.Default()
		apiCfg.TencentCloudMT = tencent_cloud_mt.Cfg{}.Default()
		apiCfg.OpenAPIYouDao = openapi_youdao.Cfg{}.Default()
		apiCfg.AliCloudMT = ali_cloud_mt.Cfg{}.Default()
		apiCfg.CaiYunAI = caiyunai.Cfg{}.Default()
	})
	return apiCfg
}

type Cfg struct {
	App            *App                  `mapstructure:"-"`
	UI             *ui.Cfg               `mapstructure:"-"`
	HuaweiCloudNlp *huawei_cloud_nlp.Cfg `mapstructure:"huawei_cloud_nlp"`
	LingVA         *ling_va.Cfg          `mapstructure:"ling_va"`
	Baidu          *baidu.Cfg            `mapstructure:"baidu"`
	TencentCloudMT *tencent_cloud_mt.Cfg `mapstructure:"tencent_cloud_mt"`
	OpenAPIYouDao  *openapi_youdao.Cfg   `mapstructure:"openapi_youdao"`
	AliCloudMT     *ali_cloud_mt.Cfg     `mapstructure:"ali_cloud_mt"`
	CaiYunAI       *caiyunai.Cfg         `mapstructure:"caiyun_ai"`
}

func (customC *Cfg) NewUITitle() string {
	return fmt.Sprintf(
		"%s-%s-%s@%s",
		customC.UI.Title, customC.App.Version,
		customC.App.Env, customC.App.Author,
	)
}
