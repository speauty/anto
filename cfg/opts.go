package cfg

import (
	"fmt"
	"github.com/spf13/viper"
)

func (customC *Cfg) Load(cfgFilePath string) error {
	if cfgFilePath == "" {
		cfgFilePath = "./"
	}

	viper.AddConfigPath(cfgFilePath)
	viper.SetConfigType("yml")
	viper.SetConfigName("cfg")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("当前配置加载失败, 错误: %s", err)
	}

	if err := viper.Unmarshal(customC); err != nil {
		return fmt.Errorf("当前配置解析失败, 错误: %s", err)
	}
	return nil
}

func (customC *Cfg) Sync() error {

	{ // sync huawei_cloud_nlp
		viper.Set("huawei_cloud_nlp.ak_id", customC.HuaweiCloudNlp.AKId)
		viper.Set("huawei_cloud_nlp.sk_key", customC.HuaweiCloudNlp.SkKey)
		viper.Set("huawei_cloud_nlp.project_id", customC.HuaweiCloudNlp.ProjectId)
		viper.Set("huawei_cloud_nlp.region", customC.HuaweiCloudNlp.Region)
		viper.Set("huawei_cloud_nlp.max_single_text_length", customC.HuaweiCloudNlp.MaxSingleTextLength)
	}

	{ // sync baidu
		viper.Set("baidu.app_id", customC.Baidu.AppId)
		viper.Set("baidu.app_key", customC.Baidu.AppKey)
		viper.Set("baidu.max_single_text_length", customC.Baidu.MaxSingleTextLength)
	}

	{ // sync ling_va
		viper.Set("ling_va.data_id", customC.LingVA.DataId)
		viper.Set("ling_va.max_single_text_length", customC.LingVA.MaxSingleTextLength)
	}

	{ // sync tencent_cloud_mt
		viper.Set("tencent_cloud_mt.secret_id", customC.TencentCloudMT.SecretId)
		viper.Set("tencent_cloud_mt.secret_key", customC.TencentCloudMT.SecretKey)
		viper.Set("tencent_cloud_mt.max_single_text_length", customC.TencentCloudMT.MaxSingleTextLength)
	}

	{ // sync openapi_youdao
		viper.Set("openapi_youdao.app_key", customC.OpenAPIYouDao.AppKey)
		viper.Set("openapi_youdao.app_secret", customC.OpenAPIYouDao.AppSecret)
		viper.Set("openapi_youdao.max_single_text_length", customC.OpenAPIYouDao.MaxSingleTextLength)
	}

	{ // sync ali_cloud_mt
		viper.Set("ali_cloud_mt.ak_id", customC.AliCloudMT.AKId)
		viper.Set("ali_cloud_mt.ak_secret", customC.AliCloudMT.AKSecret)
		viper.Set("ali_cloud_mt.region", customC.AliCloudMT.Region)
		viper.Set("ali_cloud_mt.max_single_text_length", customC.AliCloudMT.MaxSingleTextLength)
	}

	{ // sync caiyun_ai
		viper.Set("caiyun_ai.token", customC.CaiYunAI.Token)
		viper.Set("caiyun_ai.max_single_text_length", customC.CaiYunAI.MaxSingleTextLength)
	}

	{ // sync niutrans
		viper.Set("niutrans.app_key", customC.Niutrans.AppKey)
	}

	{ // sync volcengine
		viper.Set("volc_engine.access_key", customC.VolcEngine.AccessKey)
		viper.Set("volc_engine.secret_key", customC.VolcEngine.SecretKey)
	}

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("写入配置失败, 错误: %s", err)
	}
	return nil
}
