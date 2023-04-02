package cfg

import (
	"fmt"
	"sync"
	"translator/tst/tt_translator/huawei_cloud_nlp"
	"translator/tst/tt_translator/ling_va"
	"translator/tst/tt_ui"
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
	})
	return apiCfg
}

type Cfg struct {
	App            *App                  `mapstructure:"-"`
	UI             *tt_ui.Cfg            `mapstructure:"-"`
	HuaweiCloudNlp *huawei_cloud_nlp.Cfg `mapstructure:"huawei_cloud_nlp"`
	LingVA         *ling_va.Cfg          `mapstructure:"ling_va"`
}

func (customC *Cfg) NewUITitle() string {
	return fmt.Sprintf(
		"%s-%s-%s@%s",
		customC.UI.Title, customC.App.Version,
		customC.App.Env, customC.App.Author,
	)
}
