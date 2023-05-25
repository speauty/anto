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

func (config *Config) GetQPS() int             { return config.QPS }
func (config *Config) GetMaxCharNum() int      { return config.MaxCharNum }
func (config *Config) GetMaxCoroutineNum() int { return config.MaxCoroutineNum }

func (config *Config) SetQPS(num int) error {
	if err := config.ValidatorNum(num); err != nil {
		return err
	}
	config.QPS = num
	return nil
}

func (config *Config) SetMaxCharNum(num int) error {
	if err := config.ValidatorNum(num); err != nil {
		return err
	}
	config.MaxCharNum = num
	return nil
}

func (config *Config) SetMaxCoroutineNum(num int) error {
	if err := config.ValidatorNum(num); err != nil {
		return err
	}
	config.MaxCoroutineNum = num
	return nil
}
