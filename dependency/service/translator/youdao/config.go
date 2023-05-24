package youdao

import (
	"anto/dependency/service/translator"
	"github.com/spf13/viper"
)

type Config struct {
	*translator.DefaultConfig `mapstructure:"-"`
	QPS                       int `mapstructure:"qps"`
	MaxCharNum                int `mapstructure:"max_single_text_length"`
	MaxCoroutineNum           int `mapstructure:"max_coroutine_num"`
}

func (config *Config) Default() translator.ImplConfig {
	return &Config{
		MaxCharNum: 2000, QPS: 50, MaxCoroutineNum: 20,
	}
}

func (config *Config) SyncDisk(currentViper *viper.Viper) error {
	tagAndVal := config.JoinAllTagAndValue(API(), config, "mapstructure")

	for tag, val := range tagAndVal {
		currentViper.Set(tag, val)
	}
	return nil
}
