package cfg

import (
	_const "anto/common"
	"fmt"
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

func (customA App) GetDownloadUrl() string {
	return fmt.Sprintf(_const.DownloadLatestVersionUrlFormat, customA.Version, customA.Version)
}
