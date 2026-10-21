package bootstrap

import (
	"anto/cfg"
	"anto/domain/repository"
	"anto/domain/service/translator/ali_cloud_mt"
	"anto/domain/service/translator/baidu_cloud_mt"
	"anto/domain/service/translator/caiyunai"
	"anto/domain/service/translator/deepl_pro"
	"anto/domain/service/translator/deepseek"
	"anto/domain/service/translator/tencent_cloud_mt"
	"anto/domain/service/translator/youdao_batch"
	"anto/lib/log"
	"context"
)

func Boot(_ context.Context) {
	new(ResourceBuilder).Install()

	if err := cfg.Singleton().Load(); err != nil {
		panic(err)
	}

	log.Singleton()

	ali_cloud_mt.API().Init(cfg.Singleton().AliCloudMT)
	tencent_cloud_mt.API().Init(cfg.Singleton().TencentCloudMT)
	baidu_cloud_mt.API().Init(cfg.Singleton().BaiduCloudMT)
	caiyunai.API().Init(cfg.Singleton().CaiyunAI)
	deepl_pro.API().Init(cfg.Singleton().DeepLPro)
	deepseek.API().Init(cfg.Singleton().DeepSeek)
	youdao_batch.API().Init(cfg.Singleton().YoudaoBatch)
	repository.GetTranslators().Register(ali_cloud_mt.API(), tencent_cloud_mt.API(), baidu_cloud_mt.API(), caiyunai.API(), deepl_pro.API(), deepseek.API(), youdao_batch.API())
}
