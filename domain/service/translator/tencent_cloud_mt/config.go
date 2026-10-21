package tencent_cloud_mt

import (
	"anto/domain/service/translator"

	"github.com/spf13/viper"
)

type Config struct {
	*translator.DefaultConfig
	SecretID        string `mapstructure:"secret_id"`
	SecretKey       string `mapstructure:"secret_key"`
	Region          string `mapstructure:"region"`
	QPS             int    `mapstructure:"qps"`
	MaxCharNum      int    `mapstructure:"max_single_text_length"`
	MaxCoroutineNum int    `mapstructure:"max_coroutine_num"`
}

func (config *Config) Default() translator.ImplConfig {
	return &Config{SecretID: "", SecretKey: "", Region: "ap-guangzhou", QPS: 5, MaxCharNum: 2000, MaxCoroutineNum: 5}
}

func (config *Config) SyncDisk(currentViper *viper.Viper) error {
	for tag, value := range config.JoinAllTagAndValue(API(), config, "mapstructure") {
		currentViper.Set(tag, value)
	}
	return nil
}

func (config *Config) GetAK() string           { return config.SecretID }
func (config *Config) GetSK() string           { return config.SecretKey }
func (config *Config) GetRegion() string       { return config.Region }
func (config *Config) GetQPS() int             { return config.QPS }
func (config *Config) GetMaxCharNum() int      { return config.MaxCharNum }
func (config *Config) GetMaxCoroutineNum() int { return config.MaxCoroutineNum }
func (config *Config) SetAK(value string) error {
	if err := config.ValidatorStr(value); err != nil {
		return err
	}
	config.SecretID = value
	return nil
}
func (config *Config) SetSK(value string) error {
	if err := config.ValidatorStr(value); err != nil {
		return err
	}
	config.SecretKey = value
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
func (config *Config) SetRegion(value string) error {
	if !validRegion(value) {
		return translator.ErrInvalidConnectionRegion
	}
	config.Region = value
	return nil
}
