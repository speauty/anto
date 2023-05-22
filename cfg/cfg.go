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
	"anto/platform/win/ui"
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
		apiSingleton.UI = ui.Cfg{}.Default()
		apiSingleton.HuaweiCloudNlp = huawei_cloud_nlp.Cfg{}.Default()
		apiSingleton.LingVA = ling_va.Cfg{}.Default()
		apiSingleton.Baidu = baidu.Cfg{}.Default()
		apiSingleton.TencentCloudMT = tencent_cloud_mt.Cfg{}.Default()
		apiSingleton.OpenAPIYouDao = openapi_youdao.Cfg{}.Default()
		apiSingleton.AliCloudMT = ali_cloud_mt.Cfg{}.Default()
		apiSingleton.CaiYunAI = caiyunai.Cfg{}.Default()
		apiSingleton.Niutrans = niutrans.Cfg{}.Default()
		apiSingleton.VolcEngine = volcengine.Cfg{}.Default()
	})
	return apiSingleton
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
	Niutrans       *niutrans.Cfg         `mapstructure:"niutrans"`
	VolcEngine     *volcengine.Cfg       `mapstructure:"volc_engine"`
}

func (customC *Cfg) NewUITitle() string {
	return fmt.Sprintf(
		"%s-%s(作者: %s 邮箱: speauty@163.com)",
		customC.UI.Title, customC.App.Version, customC.App.Author,
	)
}
