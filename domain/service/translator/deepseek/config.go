package deepseek

import (
	"anto/domain/service/translator"

	"github.com/spf13/viper"
)

type Config struct {
	*translator.DefaultConfig
	APIKey          string `mapstructure:"api_key"`
	BaseURL         string `mapstructure:"base_url"`
	Model           string `mapstructure:"model"`
	Prompt          string `mapstructure:"prompt"`
	QPS             int    `mapstructure:"qps"`
	MaxCharNum      int    `mapstructure:"max_single_text_length"`
	MaxCoroutineNum int    `mapstructure:"max_coroutine_num"`
}

func (config *Config) Default() translator.ImplConfig {
	return &Config{
		BaseURL: "https://api.deepseek.com", Model: "deepseek-flash",
		Prompt: "请忠实翻译字幕，保留口吻、专有名词和原有换行，不要添加说明。",
		QPS:    1, MaxCharNum: 2000, MaxCoroutineNum: 1,
	}
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
