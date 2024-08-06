package cfg

import (
	"fmt"
	"github.com/spf13/viper"
)

func (customC *Cfg) GetViper() *viper.Viper {
	if customC.currentViper == nil {
		customC.currentViper = viper.New()
		customC.currentViper.AddConfigPath("./")
		customC.currentViper.SetConfigType("yml")
		customC.currentViper.SetConfigName("cfg")
	}
	return customC.currentViper
}

func (customC *Cfg) Load() error {
	currentViper := customC.GetViper()
	currentViper.AutomaticEnv()

	if err := currentViper.ReadInConfig(); err != nil {
		return fmt.Errorf("当前配置加载失败, 错误: %s", err)
	}

	if err := currentViper.Unmarshal(customC); err != nil {
		return fmt.Errorf("当前配置解析失败, 错误: %s", err)
	}
	return nil
}

func (customC *Cfg) InitConfig() error {
	currentViper := customC.GetViper()

	_ = customC.HuaweiCloudNlp.SyncDisk(currentViper)
	_ = customC.LingVA.SyncDisk(currentViper)
	_ = customC.Baidu.SyncDisk(currentViper)
	_ = customC.TencentCloudMT.SyncDisk(currentViper)
	_ = customC.OpenAPIYouDao.SyncDisk(currentViper)
	_ = customC.AliCloudMT.SyncDisk(currentViper)
	_ = customC.CaiYunAI.SyncDisk(currentViper)
	_ = customC.Niutrans.SyncDisk(currentViper)
	_ = customC.VolcEngine.SyncDisk(currentViper)
	_ = customC.YouDao.SyncDisk(currentViper)
	_ = customC.AiBaidu.SyncDisk(currentViper)
	_ = customC.MicrosoftEdge.SyncDisk(currentViper)
	_ = customC.XFYun.SyncDisk(currentViper)

	return customC.Sync()
}

func (customC *Cfg) Sync() error {
	return customC.GetViper().WriteConfig()
}
