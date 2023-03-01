package main

import (
	"gui.subtitle/src/ui"
)

func main() {
	cfg := new(ui.AppWindowCfg)
	cfg.Title = "字幕先生-Speauty出品"
	err := new(ui.AppWindow).Start(cfg)
	if err != nil {
		return
	}
}
