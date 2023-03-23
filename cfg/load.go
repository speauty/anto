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
