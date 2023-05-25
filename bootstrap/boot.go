package bootstrap

import (
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
	"context"
)

func Boot(_ context.Context) {
	new(ResourceBuilder).Install()

	if err := cfg.Singleton().Load(""); err != nil {
		panic(err)
	}

	cfg.Singleton().App.Author = common.Author
	cfg.Singleton().App.Version = common.Version
	log.Singleton()

	cfg.Singleton().UI.Title = cfg.Singleton().NewUITitle()

	huawei_cloud_nlp.API().Init(cfg.Singleton().HuaweiCloudNlp)
	ling_va.API().Init(cfg.Singleton().LingVA)
	baidu.API().Init(cfg.Singleton().Baidu)
	tencent_cloud_mt.API().Init(cfg.Singleton().TencentCloudMT)
	openapi_youdao.API().Init(cfg.Singleton().OpenAPIYouDao)
	ali_cloud_mt.API().Init(cfg.Singleton().AliCloudMT)
	caiyunai.API().Init(cfg.Singleton().CaiYunAI)
	niutrans.API().Init(cfg.Singleton().Niutrans)
	volcengine.API().Init(cfg.Singleton().VolcEngine)
	youdao.API().Init(cfg.Singleton().YouDao)

	repository.GetTranslators().Register(
		huawei_cloud_nlp.API(),
		youdao.API(), ling_va.API(), baidu.API(),
		tencent_cloud_mt.API(), openapi_youdao.API(),
		ali_cloud_mt.API(), caiyunai.API(), niutrans.API(),
		volcengine.API(),
	)
}
