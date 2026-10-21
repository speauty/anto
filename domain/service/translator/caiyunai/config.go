package caiyunai

import (
	"anto/domain/service/translator"

	"github.com/spf13/viper"
)

type Config struct {
	*translator.DefaultConfig
	Token           string `mapstructure:"token"`
	QPS             int    `mapstructure:"qps"`
	MaxCharNum      int    `mapstructure:"max_single_text_length"`
	MaxCoroutineNum int    `mapstructure:"max_coroutine_num"`
}

func (config *Config) Default() translator.ImplConfig {
	// These are conservative application defaults, not provider limits.
	return &Config{Token: "", QPS: 1, MaxCharNum: 1000, MaxCoroutineNum: 1}
}

func (config *Config) SyncDisk(currentViper *viper.Viper) error {
	for key, value := range config.JoinAllTagAndValue(API(), config, "mapstructure") {
		currentViper.Set(key, value)
	}
	return nil
}

func (config *Config) GetAK() string           { return config.Token }
func (config *Config) GetQPS() int             { return config.QPS }
func (config *Config) GetMaxCharNum() int      { return config.MaxCharNum }
func (config *Config) GetMaxCoroutineNum() int { return config.MaxCoroutineNum }

func (config *Config) SetAK(value string) error {
	if err := config.ValidatorStr(value); err != nil {
		return err
	}
	config.Token = value
	return nil
}
func (config *Config) SetQPS(value int) error {
	if err := config.ValidatorNum(value); err != nil {
		return err
	}
	config.QPS = value
	return nil
}
func (config *Config) SetMaxCharNum(value int) error {
	if err := config.ValidatorNum(value); err != nil {
		return err
	}
	config.MaxCharNum = value
	return nil
}
func (config *Config) SetMaxCoroutineNum(value int) error {
	if err := config.ValidatorNum(value); err != nil {
		return err
	}
	config.MaxCoroutineNum = value
	return nil
}
