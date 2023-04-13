package main

import (
	"anto/boot"
	"anto/cfg"
	_const "anto/const"
	"anto/domain"
	"anto/menu"
	"anto/page"
	"anto/tst/tt_log"
	"anto/tst/tt_translator/ali_cloud_mt"
	"anto/tst/tt_translator/baidu"
	"anto/tst/tt_translator/caiyunai"
	"anto/tst/tt_translator/huawei_cloud_nlp"
	"anto/tst/tt_translator/ling_va"
	"anto/tst/tt_translator/openapi_youdao"
	"anto/tst/tt_translator/tencent_cloud_mt"
	"anto/tst/tt_translator/youdao"
	"anto/tst/tt_ui"
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
	baidu.GetInstance().Init(cfg.GetInstance().Baidu)
	tencent_cloud_mt.GetInstance().Init(cfg.GetInstance().TencentCloudMT)
	openapi_youdao.GetInstance().Init(cfg.GetInstance().OpenAPIYouDao)
	ali_cloud_mt.GetInstance().Init(cfg.GetInstance().AliCloudMT)
	caiyunai.GetInstance().Init(cfg.GetInstance().CaiYunAI)

	domain.GetTranslators().Register(
		huawei_cloud_nlp.GetInstance(),
		youdao.GetInstance(), ling_va.GetInstance(), baidu.GetInstance(),
		tencent_cloud_mt.GetInstance(), openapi_youdao.GetInstance(),
		ali_cloud_mt.GetInstance(), caiyunai.GetInstance(),
	)

	tt_ui.GetInstance().RegisterMenus(menu.GetInstance().GetMenus())

	tt_ui.GetInstance().RegisterPages(
		page.GetAboutUs(), page.GetSettings(), page.GetUsage(), page.GetSubripTranslate(),
	)

	if err := tt_ui.GetInstance().Init(cfg.GetInstance().UI); err != nil {
		panic(err)
	}

	_ = tt_ui.GetInstance().GoPage(page.GetAboutUs().GetId())

	tt_ui.GetInstance().Run()
}
