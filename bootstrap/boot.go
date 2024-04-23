package bootstrap

import (
	"anto/cfg"
	"anto/common"
	"anto/domain/repository"
	"anto/domain/service/translator/ai_baidu"
	"anto/domain/service/translator/ali_cloud_mt"
	"anto/domain/service/translator/baidu"
	"anto/domain/service/translator/caiyunai"
	"anto/domain/service/translator/deepl"
	"anto/domain/service/translator/deepl_pro"
	"anto/domain/service/translator/g_deepl_x"
	"anto/domain/service/translator/google_cloud"
	"anto/domain/service/translator/huawei_cloud_nlp"
	"anto/domain/service/translator/microsoft_edge"
	"anto/domain/service/translator/niutrans"
	"anto/domain/service/translator/openai"
	"anto/domain/service/translator/openai_sweet"
	"anto/domain/service/translator/openapi_youdao"
	"anto/domain/service/translator/tencent_cloud_mt"
	"anto/domain/service/translator/volcengine"
	"anto/lib/log"
	"context"
)

func Boot(_ context.Context) {
	new(ResourceBuilder).Install()

	if err := cfg.Singleton().Load(); err != nil {
		panic(err)
	}

	cfg.Singleton().App.Author = common.Author
	cfg.Singleton().App.Version = common.Version
	log.Singleton()

	cfg.Singleton().UI.Title = cfg.Singleton().NewUITitle()

	huawei_cloud_nlp.API().Init(cfg.Singleton().HuaweiCloudNlp)
	baidu.API().Init(cfg.Singleton().Baidu)
	tencent_cloud_mt.API().Init(cfg.Singleton().TencentCloudMT)
	openapi_youdao.API().Init(cfg.Singleton().OpenAPIYouDao)
	ali_cloud_mt.API().Init(cfg.Singleton().AliCloudMT)
	caiyunai.API().Init(cfg.Singleton().CaiYunAI)
	niutrans.API().Init(cfg.Singleton().Niutrans)
	volcengine.API().Init(cfg.Singleton().VolcEngine)
	g_deepl_x.API().Init(cfg.Singleton().YouDao)
	google_cloud.API().Init(cfg.Singleton().GoogleCloud)
	openai.API().Init(cfg.Singleton().OpenAI)
	openai_sweet.API().Init(cfg.Singleton().OpenAISweet)
	deepl.API().Init(cfg.Singleton().DeepL)
	deepl_pro.API().Init(cfg.Singleton().DeepLPro)
	ai_baidu.API().Init(cfg.Singleton().AiBaidu)
	microsoft_edge.API().Init(cfg.Singleton().MicrosoftEdge)

	repository.GetTranslators().Register(
		huawei_cloud_nlp.API(), baidu.API(),
		tencent_cloud_mt.API(), openapi_youdao.API(),
		ali_cloud_mt.API(), caiyunai.API(), niutrans.API(),
		volcengine.API(), g_deepl_x.API(),
		google_cloud.API(), openai.API(), deepl.API(), deepl_pro.API(),
		openai_sweet.API(), ai_baidu.API(), microsoft_edge.API(),
	)
}
