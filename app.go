package main

import (
	"translator/boot"
	"translator/cfg"
	_const "translator/const"
	"translator/domain"
	"translator/menu"
	"translator/page"
	"translator/tst/tt_log"
	"translator/tst/tt_translator/huawei_cloud_nlp"
	"translator/tst/tt_translator/ling_va"
	"translator/tst/tt_translator/youdao"
	"translator/tst/tt_ui"
)

func main() {
	new(boot.ResourceBuilder).Install()

	if err := cfg.GetInstance().Load(""); err != nil {
		panic(err)
	}
	cfg.GetInstance().App.Author = _const.Author
	cfg.GetInstance().App.Version = _const.Version
	tt_log.GetInstance()

	cfg.GetInstance().UI.Title = cfg.GetInstance().NewUITitle()

	huawei_cloud_nlp.GetInstance().Init(cfg.GetInstance().HuaweiCloudNlp)
	ling_va.GetInstance().Init(cfg.GetInstance().LingVA)

	domain.GetTranslators().Register(
		huawei_cloud_nlp.GetInstance(),
		youdao.GetInstance(),
		ling_va.GetInstance(),
	)

	tt_ui.GetInstance().RegisterMenus(menu.GetInstance().GetMenus())

	tt_ui.GetInstance().RegisterPages(
		page.GetAboutUs(), page.GetSettings(), page.GetUsage(), page.GetSubripTranslate(),
	)

	if err := tt_ui.GetInstance().Init(cfg.GetInstance().UI); err != nil {
		panic(err)
	}

	_ = tt_ui.GetInstance().GoPage(page.GetSettings().GetId())

	tt_ui.GetInstance().Run()
}
