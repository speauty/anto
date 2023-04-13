package cfg

import (
	"anto/tst/tt_translator/ali_cloud_mt"
	"anto/tst/tt_translator/baidu"
	"anto/tst/tt_translator/caiyunai"
	"anto/tst/tt_translator/huawei_cloud_nlp"
	"anto/tst/tt_translator/ling_va"
	"anto/tst/tt_translator/openapi_youdao"
	"anto/tst/tt_translator/tencent_cloud_mt"
	"anto/tst/tt_ui"
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
		apiCfg.UI = tt_ui.Cfg{}.Default()
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
	UI             *tt_ui.Cfg            `mapstructure:"-"`
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
