package cfg

import (
	_const "anto/common"
)

type App struct {
	Author  string `mapstructure:"author"`
	Version string `mapstructure:"version"`
}

func (customA App) Default() *App {
	return &App{
		Author:  _const.Author,
		Version: _const.Version,
	}
}
