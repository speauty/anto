package cross_platform_fyne

import "fyne.io/fyne/v2"

var (
	appWindowDefaultWidth      float32 = 400
	appWindowDefaultHeight     float32 = 300
	appMainWindowDefaultWidth  float32 = 800
	appMainWindowDefaultHeight float32 = 600
	appChanConsoleCnt          int     = 50
)

type Config struct {
	AppName                    string  `yaml:"app_name"`
	AppVersion                 string  `yaml:"app_version"`                    // v1.0.1
	AppMode                    string  `yaml:"app_mode"`                       // debug release
	AppWindowDefaultWidth      float32 `yaml:"app_window_default_width"`       // 应用窗口默认宽度
	AppWindowDefaultHeight     float32 `yaml:"app_window_default_height"`      // 应用窗口默认高度
	AppMainWindowDefaultWidth  float32 `yaml:"app_main_window_default_width"`  // 应用主窗口默认宽度
	AppMainWindowDefaultHeight float32 `yaml:"app_main_window_default_height"` // 应用主窗口默认高度
}

func (c *Config) Default() *Config {
	return &Config{
		AppName:                    "桌面应用(ANTO)",
		AppVersion:                 "v1.0.0",
		AppMode:                    "debug",
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
