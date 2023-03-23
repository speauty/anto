package cfg

import (
	"fmt"
	"sync"
	"translator/tst/tt_ui"
)

var (
	apiCfg  *Cfg
	onceCfg sync.Once
)

func GetInstance() *Cfg {
	onceCfg.Do(func() {
		apiCfg = new(Cfg)
	})
	return apiCfg
}

type Cfg struct {
	App *App       `mapstructure:"app"`
	UI  *tt_ui.Cfg `mapstructure:"ui"`
}

func (customC *Cfg) NewUITitle() string {
	return fmt.Sprintf(
		"%s-%s-%s@%s",
		customC.UI.Title, customC.App.Version,
		customC.App.Env, customC.App.Author,
	)
}
