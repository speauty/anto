package cfg

import (
	"anto/domain/service/translator/ai_baidu"
	"anto/domain/service/translator/ali_cloud_mt"
	"anto/domain/service/translator/baidu"
	"anto/domain/service/translator/caiyunai"
	"anto/domain/service/translator/deepl"
	"anto/domain/service/translator/deepl_pro"
	"anto/domain/service/translator/g_deepl_x"
	"anto/domain/service/translator/google_cloud"
	"anto/domain/service/translator/huawei_cloud_nlp"
	"anto/domain/service/translator/ling_va"
	"anto/domain/service/translator/niutrans"
	"anto/domain/service/translator/openai"
	"anto/domain/service/translator/openai_sweet"
	"anto/domain/service/translator/openapi_youdao"
	"anto/domain/service/translator/tencent_cloud_mt"
	"anto/domain/service/translator/volcengine"
	"anto/domain/service/translator/youdao"
	"anto/platform/win/ui"
	"fmt"
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
		apiSingleton.App = App{}.Default()
		apiSingleton.UI = ui.Cfg{}.Default()
		apiSingleton.HuaweiCloudNlp = new(huawei_cloud_nlp.Config).Default().(*huawei_cloud_nlp.Config)
		apiSingleton.LingVA = new(ling_va.Config).Default().(*ling_va.Config)
		apiSingleton.Baidu = new(baidu.Config).Default().(*baidu.Config)
		apiSingleton.TencentCloudMT = new(tencent_cloud_mt.Config).Default().(*tencent_cloud_mt.Config)
		apiSingleton.OpenAPIYouDao = new(openapi_youdao.Config).Default().(*openapi_youdao.Config)
		apiSingleton.AliCloudMT = new(ali_cloud_mt.Config).Default().(*ali_cloud_mt.Config)
		apiSingleton.CaiYunAI = new(caiyunai.Config).Default().(*caiyunai.Config)
		apiSingleton.Niutrans = new(niutrans.Config).Default().(*niutrans.Config)
		apiSingleton.VolcEngine = new(volcengine.Config).Default().(*volcengine.Config)
		apiSingleton.YouDao = new(youdao.Config).Default().(*youdao.Config)
		apiSingleton.GoogleCloud = new(google_cloud.Config).Default().(*google_cloud.Config)
		apiSingleton.OpenAI = new(openai.Config).Default().(*openai.Config)
		apiSingleton.OpenAISweet = new(openai_sweet.Config).Default().(*openai_sweet.Config)
		apiSingleton.DeepL = new(deepl.Config).Default().(*deepl.Config)
		apiSingleton.DeepLPro = new(deepl_pro.Config).Default().(*deepl_pro.Config)
		apiSingleton.AiBaidu = new(ai_baidu.Config).Default().(*ai_baidu.Config)
	})
	return apiSingleton
}

type Cfg struct {
	App            *App                     `mapstructure:"-"`
	UI             *ui.Cfg                  `mapstructure:"-"`
	HuaweiCloudNlp *huawei_cloud_nlp.Config `mapstructure:"huawei_cloud_nlp"`
	LingVA         *ling_va.Config          `mapstructure:"ling_va"`
	Baidu          *baidu.Config            `mapstructure:"baidu"`
	TencentCloudMT *tencent_cloud_mt.Config `mapstructure:"tencent_cloud_mt"`
	OpenAPIYouDao  *openapi_youdao.Config   `mapstructure:"openapi_youdao"`
	AliCloudMT     *ali_cloud_mt.Config     `mapstructure:"ali_cloud_mt"`
	CaiYunAI       *caiyunai.Config         `mapstructure:"caiyun_ai"`
	Niutrans       *niutrans.Config         `mapstructure:"niutrans"`
	VolcEngine     *volcengine.Config       `mapstructure:"volc_engine"`
	YouDao         *youdao.Config           `mapstructure:"youdao"`
	GDeeplX        *g_deepl_x.Config        `mapstructure:"g_deepl_x"`
	GoogleCloud    *google_cloud.Config     `mapstructure:"google_cloud"`
	OpenAI         *openai.Config           `mapstructure:"openai"`
	OpenAISweet    *openai_sweet.Config     `mapstructure:"openai_sweet"`
	DeepL          *deepl.Config            `mapstructure:"deepl"`
	DeepLPro       *deepl_pro.Config        `mapstructure:"deepl_pro"`
	AiBaidu        *ai_baidu.Config         `mapstructure:"ai_baidu"`

	currentViper *viper.Viper `mapstructure:"-"`
}

func (customC *Cfg) NewUITitle() string {
	return fmt.Sprintf(
		"%s-%s speauty@163.com",
		customC.UI.Title, customC.App.Version,
	)
}
