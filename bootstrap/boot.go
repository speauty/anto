package bootstrap

import (
	"anto/cfg"
	"anto/common"
	"anto/domain/repository"
	"anto/domain/service/translator/ali_cloud_mt"
	"anto/domain/service/translator/baidu"
	"anto/domain/service/translator/caiyunai"
	"anto/domain/service/translator/huawei_cloud_nlp"
	"anto/domain/service/translator/ling_va"
	"anto/domain/service/translator/niutrans"
	"anto/domain/service/translator/openapi_youdao"
	"anto/domain/service/translator/tencent_cloud_mt"
	"anto/domain/service/translator/volcengine"
	"anto/domain/service/translator/youdao"
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
}
