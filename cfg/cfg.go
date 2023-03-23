package cfg

import "sync"

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
	App *App `mapstructure:"app"`
}
