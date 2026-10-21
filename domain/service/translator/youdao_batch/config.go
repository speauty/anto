package youdao_batch

import (
	"anto/domain/service/translator"
	"github.com/spf13/viper"
)

type Config struct {
	*translator.DefaultConfig
	AppKey          string `mapstructure:"app_key"`
	AppSecret       string `mapstructure:"app_secret"`
	QPS             int    `mapstructure:"qps"`
	MaxCharNum      int    `mapstructure:"max_single_text_length"`
	MaxCoroutineNum int    `mapstructure:"max_coroutine_num"`
}

func (config *Config) Default() translator.ImplConfig {
	return &Config{AppKey: "", AppSecret: "", QPS: 1, MaxCharNum: 5000, MaxCoroutineNum: 1}
}

func (config *Config) SyncDisk(currentViper *viper.Viper) error {
	for key, value := range config.JoinAllTagAndValue(API(), config, "mapstructure") {
		currentViper.Set(key, value)
	}
	return nil
}

func (config *Config) GetAK() string           { return config.AppKey }
func (config *Config) GetSK() string           { return config.AppSecret }
func (config *Config) GetQPS() int             { return config.QPS }
func (config *Config) GetMaxCharNum() int      { return config.MaxCharNum }
func (config *Config) GetMaxCoroutineNum() int { return config.MaxCoroutineNum }

func (config *Config) SetAK(value string) error {
	if err := config.ValidatorStr(value); err != nil {
		return err
	}
	config.AppKey = value
	return nil
}
func (config *Config) SetSK(value string) error {
	if err := config.ValidatorStr(value); err != nil {
		return err
	}
	config.AppSecret = value
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
	if value > 5000 {
		return &limitError{}
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

type limitError struct{}

func (*limitError) Error() string { return "有道智云单次最大长度不能超过 5000 个字符" }
