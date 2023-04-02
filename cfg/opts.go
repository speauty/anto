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
	}

	{ // sync baidu
		viper.Set("baidu.app_id", customC.Baidu.AppId)
		viper.Set("baidu.app_key", customC.Baidu.AppKey)
	}

	{ // sync ling_va
		viper.Set("ling_va.data_id", customC.LingVA.DataId)
	}

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("写入配置失败, 错误: %s", err)
	}
	return nil
}
