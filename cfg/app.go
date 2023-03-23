package cfg

type App struct {
	Env     string `mapstructure:"env"`
	Author  string `mapstructure:"-"`
	Version string `mapstructure:"-"`
}
