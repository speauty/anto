package ui

import (
	_const "anto/common"
)

type Cfg struct {
	Title       string `mapstructure:"title"`
	Icon        string `mapstructure:"icon"`
	ResourceDir string `mapstructure:"resource_dir"`
}

func (customC Cfg) Default() *Cfg {
	return &Cfg{Title: _const.UITitle, Icon: _const.UIIcon, ResourceDir: _const.UIResourceDir}
}
