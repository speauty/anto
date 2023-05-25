package ali_cloud_mt

import (
	"anto/domain/service/translator"
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	*translator.DefaultConfig
	AKId            string `mapstructure:"ak_id"`
	AKSecret        string `mapstructure:"ak_secret"`
	Region          string `mapstructure:"region"`
	QPS             int    `mapstructure:"qps"`
	MaxCharNum      int    `mapstructure:"max_single_text_length"`
	MaxCoroutineNum int    `mapstructure:"max_coroutine_num"`
}

func (config *Config) Default() translator.ImplConfig {
	return &Config{
		AKId: "", AKSecret: "", Region: "cn-hangzhou",
		MaxCharNum: 3000, QPS: 50, MaxCoroutineNum: 25,
	}
}

func (config *Config) SyncDisk(currentViper *viper.Viper) error {
	tagAndVal := config.JoinAllTagAndValue(API(), config, "mapstructure")

	for tag, val := range tagAndVal {
		fmt.Println(tag, val)
		currentViper.Set(tag, val)
	}
	return nil
}

func (config *Config) GetAK() string     { return config.AKId }
func (config *Config) GetSK() string     { return config.AKSecret }
func (config *Config) GetRegion() string { return config.Region }

func (config *Config) GetQPS() int             { return config.QPS }
func (config *Config) GetMaxCharNum() int      { return config.MaxCharNum }
func (config *Config) GetMaxCoroutineNum() int { return config.MaxCoroutineNum }

func (config *Config) SetAK(str string) error {
	if err := config.ValidatorStr(str); err != nil {
		return err
	}
	config.AKId = str
	return nil
}
func (config *Config) SetSK(str string) error {
	if err := config.ValidatorStr(str); err != nil {
		return err
	}
	config.AKSecret = str
	return nil
}
func (config *Config) SetRegion(str string) error {
	if err := config.ValidatorStr(str); err != nil {
		return err
	}
	config.Region = str
	return nil
}

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
