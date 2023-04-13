package main

import (
	"anto/boot"
	"anto/cfg"
	_const "anto/const"
	"anto/cron/detector"
	"anto/cron/reader"
	"anto/cron/translate"
	"anto/cron/writer"
	"anto/dependency/service/translator/ali_cloud_mt"
	"anto/dependency/service/translator/baidu"
	"anto/dependency/service/translator/caiyunai"
	"anto/dependency/service/translator/huawei_cloud_nlp"
	"anto/dependency/service/translator/ling_va"
	"anto/dependency/service/translator/openapi_youdao"
	"anto/dependency/service/translator/tencent_cloud_mt"
	"anto/dependency/service/translator/youdao"
	"anto/domain"
	"anto/lib/log"
	"anto/lib/nohup"
	"anto/lib/ui"
	"anto/menu"
	"anto/page"
	"context"
)

func main() {
	ctx := context.Background()

	new(boot.ResourceBuilder).Install()

	if err := cfg.GetInstance().Load(""); err != nil {
		panic(err)
	}
	cfg.GetInstance().App.Author = _const.Author
	cfg.GetInstance().App.Version = _const.Version
	log.Singleton()

	cfg.GetInstance().UI.Title = cfg.GetInstance().NewUITitle()

	huawei_cloud_nlp.Singleton().Init(cfg.GetInstance().HuaweiCloudNlp)
	ling_va.Singleton().Init(cfg.GetInstance().LingVA)
	baidu.Singleton().Init(cfg.GetInstance().Baidu)
	tencent_cloud_mt.Singleton().Init(cfg.GetInstance().TencentCloudMT)
	openapi_youdao.Singleton().Init(cfg.GetInstance().OpenAPIYouDao)
	ali_cloud_mt.Singleton().Init(cfg.GetInstance().AliCloudMT)
	caiyunai.Singleton().Init(cfg.GetInstance().CaiYunAI)

	domain.GetTranslators().Register(
		huawei_cloud_nlp.Singleton(),
		youdao.Singleton(), ling_va.Singleton(), baidu.Singleton(),
		tencent_cloud_mt.Singleton(), openapi_youdao.Singleton(),
		ali_cloud_mt.Singleton(), caiyunai.Singleton(),
	)

	ui.GetInstance().RegisterMenus(menu.GetInstance().GetMenus())

	ui.GetInstance().RegisterPages(
		page.GetAboutUs(), page.GetSettings(), page.GetUsage(), page.GetSubripTranslate(),
	)

	if err := ui.GetInstance().Init(cfg.GetInstance().UI); err != nil {
		panic(err)
	}

	_ = ui.GetInstance().GoPage(page.GetAboutUs().GetId())

	nohup.NewResident(
		ctx,
		detector.Singleton(), reader.Singleton(), translate.Singleton(), writer.Singleton(),
		ui.GetInstance(),
	)
}
