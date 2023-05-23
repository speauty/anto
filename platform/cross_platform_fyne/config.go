package cross_platform_fyne

import (
	"anto/common"
	"fmt"
	"fyne.io/fyne/v2"
)

var (
	appWindowDefaultWidth      float32 = 400
	appWindowDefaultHeight     float32 = 300
	appMainWindowDefaultWidth  float32 = 800
	appMainWindowDefaultHeight float32 = 600
	appChanConsoleCnt          int     = 50
)

type Config struct {
	AppName                    string  `mapstructure:"-"`
	AppVersion                 string  `mapstructure:"-"` // v1.0.1
	AppMode                    string  `mapstructure:"-"` // debug release
	AppWindowDefaultWidth      float32 `mapstructure:"-"` // 应用窗口默认宽度
	AppWindowDefaultHeight     float32 `mapstructure:"-"` // 应用窗口默认高度
	AppMainWindowDefaultWidth  float32 `mapstructure:"-"` // 应用主窗口默认宽度
	AppMainWindowDefaultHeight float32 `mapstructure:"-"` // 应用主窗口默认高度
}

func (c *Config) Default() *Config {
	return &Config{
		AppName:                    common.AppName,
		AppVersion:                 common.Version,
		AppMode:                    "release",
		AppWindowDefaultWidth:      appWindowDefaultWidth,
		AppWindowDefaultHeight:     appWindowDefaultHeight,
		AppMainWindowDefaultWidth:  appMainWindowDefaultWidth,
		AppMainWindowDefaultHeight: appMainWindowDefaultHeight,
	}
}

func (c *Config) MainWindowSize() fyne.Size {
	return fyne.NewSize(c.AppMainWindowDefaultWidth, c.AppMainWindowDefaultHeight)
}

func (c *Config) WindowSize() fyne.Size {
	return fyne.NewSize(c.AppWindowDefaultWidth, c.AppWindowDefaultHeight)
}

func (c *Config) IsDebug() bool {
	return c.AppMode == "debug"
}

func (c *Config) IsRelease() bool {
	return c.AppMode == "release"
}

func (c *Config) GetMainWindowTitle() string {
	return fmt.Sprintf("%s@%s %s (免费开源应用)", c.AppName, common.Author, c.AppVersion)
}
