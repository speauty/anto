package main

import (
	"translator/cfg"
	_const "translator/const"
	"translator/tst/tt_log"
	"translator/tst/tt_ui"
)

func main() {

	if err := cfg.GetInstance().Load(""); err != nil {
		panic(err)
	}
	cfg.GetInstance().App.Author = _const.Author
	cfg.GetInstance().App.Version = _const.Version
	tt_log.GetInstance()

	cfg.GetInstance().UI.Title = cfg.GetInstance().NewUITitle()
	if err := tt_ui.GetInstance().Init(cfg.GetInstance().UI); err != nil {
		panic(err)
	}

	tt_ui.GetInstance().Run()

}
