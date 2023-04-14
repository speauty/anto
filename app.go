package main

import (
	"anto/bootstrap"
	"anto/cfg"
	"anto/common"
	"anto/dependency/repository"
	"anto/dependency/service/translator/ali_cloud_mt"
	"anto/dependency/service/translator/baidu"
	"anto/dependency/service/translator/caiyunai"
	"anto/dependency/service/translator/huawei_cloud_nlp"
	"anto/dependency/service/translator/ling_va"
	"anto/dependency/service/translator/niutrans"
	"anto/dependency/service/translator/openapi_youdao"
	"anto/dependency/service/translator/tencent_cloud_mt"
	"anto/dependency/service/translator/volcengine"
	"anto/dependency/service/translator/youdao"
	"anto/lib/log"
	"anto/platform/win"
	"context"
)

func main() {
	ctx := context.Background()

	new(bootstrap.ResourceBuilder).Install()

	if err := cfg.Singleton().Load(""); err != nil {
		panic(err)
	}

	cfg.Singleton().App.Author = common.Author
	cfg.Singleton().App.Version = common.Version
	log.Singleton()

	cfg.Singleton().UI.Title = cfg.Singleton().NewUITitle()

	huawei_cloud_nlp.Singleton().Init(cfg.Singleton().HuaweiCloudNlp)
	ling_va.Singleton().Init(cfg.Singleton().LingVA)
	baidu.Singleton().Init(cfg.Singleton().Baidu)
	tencent_cloud_mt.Singleton().Init(cfg.Singleton().TencentCloudMT)
	openapi_youdao.Singleton().Init(cfg.Singleton().OpenAPIYouDao)
	ali_cloud_mt.Singleton().Init(cfg.Singleton().AliCloudMT)
	caiyunai.Singleton().Init(cfg.Singleton().CaiYunAI)
	niutrans.Singleton().Init(cfg.Singleton().Niutrans)
	volcengine.Singleton().Init(cfg.Singleton().VolcEngine)

	repository.GetTranslators().Register(
		huawei_cloud_nlp.Singleton(),
		youdao.Singleton(), ling_va.Singleton(), baidu.Singleton(),
		tencent_cloud_mt.Singleton(), openapi_youdao.Singleton(),
		ali_cloud_mt.Singleton(), caiyunai.Singleton(), niutrans.Singleton(),
		volcengine.Singleton(),
	)

	win.Run(ctx)
}
