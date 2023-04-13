package cfg

import _const "anto/const"

type App struct {
	Env     string `mapstructure:"env"`
	Author  string `mapstructure:"author"`
	Version string `mapstructure:"version"`
}

func (customA App) Default() *App {
	return &App{
		Env:     "release",
		Author:  _const.Author,
		Version: _const.Version,
	}
}
