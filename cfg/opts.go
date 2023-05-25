package cfg

import (
	"fmt"
	"github.com/spf13/viper"
)

func (customC *Cfg) GetViper() *viper.Viper {
	if customC.currentViper == nil {
		customC.currentViper = viper.New()
	}
	return customC.currentViper
}

func (customC *Cfg) Load(cfgFilePath string) error {
	if cfgFilePath == "" {
		cfgFilePath = "./"
	}
	currentViper := customC.GetViper()
	currentViper.AddConfigPath(cfgFilePath)
	currentViper.SetConfigType("yml")
	currentViper.SetConfigName("cfg")
	currentViper.AutomaticEnv()

	if err := currentViper.ReadInConfig(); err != nil {
		return fmt.Errorf("当前配置加载失败, 错误: %s", err)
	}

	if err := currentViper.Unmarshal(customC); err != nil {
		return fmt.Errorf("当前配置解析失败, 错误: %s", err)
	}
	return nil
}

func (customC *Cfg) Sync() error {
	return customC.GetViper().WriteConfig()
}
