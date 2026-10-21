package cfg

import (
	"anto/domain/service/translator/ali_cloud_mt"
	"anto/domain/service/translator/baidu_cloud_mt"
	"anto/domain/service/translator/caiyunai"
	"anto/domain/service/translator/deepl_pro"
	"anto/domain/service/translator/deepseek"
	"anto/domain/service/translator/tencent_cloud_mt"
	"anto/domain/service/translator/youdao_batch"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestExistingConfigWithoutTencentSectionKeepsTencentDefaults(t *testing.T) {
	values := viper.New()
	values.SetConfigType("yaml")
	if err := values.ReadConfig(strings.NewReader("ali_cloud_mt:\n  region: cn-hangzhou\n")); err != nil {
		t.Fatal(err)
	}
	config := &Cfg{
		AliCloudMT:     new(ali_cloud_mt.Config).Default().(*ali_cloud_mt.Config),
		TencentCloudMT: new(tencent_cloud_mt.Config).Default().(*tencent_cloud_mt.Config),
		BaiduCloudMT:   new(baidu_cloud_mt.Config).Default().(*baidu_cloud_mt.Config),
		CaiyunAI:       new(caiyunai.Config).Default().(*caiyunai.Config),
		DeepLPro:       new(deepl_pro.Config).Default().(*deepl_pro.Config),
		DeepSeek:       new(deepseek.Config).Default().(*deepseek.Config),
		YoudaoBatch:    new(youdao_batch.Config).Default().(*youdao_batch.Config),
	}
	if err := values.Unmarshal(config); err != nil {
		t.Fatal(err)
	}
	if config.TencentCloudMT == nil || config.TencentCloudMT.GetRegion() != "ap-guangzhou" || config.TencentCloudMT.GetQPS() != 5 {
		t.Fatalf("Tencent defaults missing after loading older config: %+v", config.TencentCloudMT)
	}
	if config.BaiduCloudMT == nil || config.BaiduCloudMT.GetQPS() != 1 || config.BaiduCloudMT.GetMaxCharNum() != 1000 {
		t.Fatalf("Baidu defaults missing after loading older config: %+v", config.BaiduCloudMT)
	}
	if config.CaiyunAI == nil || config.CaiyunAI.GetQPS() != 1 || config.CaiyunAI.GetMaxCharNum() != 1000 {
		t.Fatalf("Caiyun defaults missing after loading older config: %+v", config.CaiyunAI)
	}
	if config.DeepLPro == nil || config.DeepLPro.GetQPS() != 5 || config.DeepLPro.GetMaxCharNum() != 3000 {
		t.Fatalf("DeepL defaults missing after loading older config: %+v", config.DeepLPro)
	}
	if config.DeepSeek == nil || config.DeepSeek.Model != "deepseek-flash" || config.DeepSeek.BaseURL != "https://api.deepseek.com" {
		t.Fatalf("DeepSeek defaults missing after loading older config: %+v", config.DeepSeek)
	}
	if config.YoudaoBatch == nil || config.YoudaoBatch.GetQPS() != 1 || config.YoudaoBatch.GetMaxCharNum() != 5000 {
		t.Fatalf("Youdao defaults missing after loading older config: %+v", config.YoudaoBatch)
	}
}
