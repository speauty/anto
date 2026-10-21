package deepl_pro

import (
	"anto/domain/service/translator"

	"github.com/spf13/viper"
)

type Config struct {
	*translator.DefaultConfig
	APIKey          string `mapstructure:"api_key"`
	QPS             int    `mapstructure:"qps"`
	MaxCharNum      int    `mapstructure:"max_single_text_length"`
	MaxCoroutineNum int    `mapstructure:"max_coroutine_num"`
}

func (config *Config) Default() translator.ImplConfig {
	// Local throttling defaults; DeepL does not publish a fixed QPS ceiling.
	return &Config{APIKey: "", QPS: 5, MaxCharNum: 3000, MaxCoroutineNum: 5}
}

func (config *Config) SyncDisk(currentViper *viper.Viper) error {
	for key, value := range config.JoinAllTagAndValue(API(), config, "mapstructure") {
		currentViper.Set(key, value)
	}
	return nil
}

func (config *Config) GetAK() string           { return config.APIKey }
func (config *Config) GetQPS() int             { return config.QPS }
func (config *Config) GetMaxCharNum() int      { return config.MaxCharNum }
func (config *Config) GetMaxCoroutineNum() int { return config.MaxCoroutineNum }

func (config *Config) SetAK(value string) error {
	if err := config.ValidatorStr(value); err != nil {
		return err
	}
	config.APIKey = value
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
