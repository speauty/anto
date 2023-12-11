package openai

import (
	"anto/domain/service/translator"
	"github.com/spf13/viper"
)

type Config struct {
	*translator.DefaultConfig
	AppKey          string `mapstructure:"app_key"`
	ProjectKey      string `mapstructure:"project_key"`
	QPS             int    `mapstructure:"qps"`
	MaxCharNum      int    `mapstructure:"max_single_text_length"`
	MaxCoroutineNum int    `mapstructure:"max_coroutine_num"`
}

func (config *Config) Default() translator.ImplConfig {
	return &Config{
		AppKey: "", ProjectKey: "gpt-3.5-turbo",
		MaxCharNum: 2000, QPS: 1, MaxCoroutineNum: 1,
	}
}

func (config *Config) SyncDisk(currentViper *viper.Viper) error {
	tagAndVal := config.JoinAllTagAndValue(API(), config, "mapstructure")

	for tag, val := range tagAndVal {
		currentViper.Set(tag, val)
	}
	return nil
}

func (config *Config) GetAK() string           { return config.AppKey }
func (config *Config) GetProjectKey() string   { return config.ProjectKey }
func (config *Config) GetQPS() int             { return config.QPS }
func (config *Config) GetMaxCharNum() int      { return config.MaxCharNum }
func (config *Config) GetMaxCoroutineNum() int { return config.MaxCoroutineNum }

func (config *Config) SetAK(str string) error {
	if err := config.ValidatorStr(str); err != nil {
		return err
	}
	config.AppKey = str
	return nil
}

func (config *Config) SetProjectKey(projectKey string) error {
	if projectKey == "" {
		projectKey = "gpt-3.5-turbo"
	}
	config.ProjectKey = projectKey
	return nil
}

func (config *Config) SetQPS(num int) error {
	config.QPS = 1
	return nil
}

func (config *Config) SetMaxCharNum(num int) error {
	if err := config.ValidatorNum(num); err != nil {
		return err
	}
	if num > 2000 {
		num = 2000
	}
	config.MaxCharNum = num
	return nil
}

func (config *Config) SetMaxCoroutineNum(num int) error {
	config.MaxCoroutineNum = 1
	return nil
}
