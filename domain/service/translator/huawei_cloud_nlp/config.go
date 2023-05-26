package huawei_cloud_nlp

import (
	"anto/domain/service/translator"
	"github.com/spf13/viper"
)

type Config struct {
	*translator.DefaultConfig
	AKId            string `mapstructure:"ak_id"`      // Access Key
	SkKey           string `mapstructure:"sk_key"`     // Secret Access Key
	Region          string `mapstructure:"region"`     // 当前接口开发的区域, 目前仅支持华北-北京四终端节点
	ProjectId       string `mapstructure:"project_id"` // 项目ID
	QPS             int    `mapstructure:"qps"`
	MaxCharNum      int    `mapstructure:"max_single_text_length"`
	MaxCoroutineNum int    `mapstructure:"max_coroutine_num"`
}

func (config *Config) Default() translator.ImplConfig {
	return &Config{
		AKId: "", SkKey: "", Region: "cn-north-4", ProjectId: "",
		MaxCharNum: 2000, QPS: 20, MaxCoroutineNum: 10,
	}
}

func (config *Config) SyncDisk(currentViper *viper.Viper) error {
	tagAndVal := config.JoinAllTagAndValue(API(), config, "mapstructure")

	for tag, val := range tagAndVal {
		currentViper.Set(tag, val)
	}
	return nil
}

func (config *Config) GetAK() string { return config.AKId }
func (config *Config) GetSK() string { return config.SkKey }
func (config *Config) GetRegion() string {
	if config.Region != "" {
		return config.Region
	}
	return "cn-north-4"
}
func (config *Config) GetProjectKey() string { return config.ProjectId }

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
	config.SkKey = str
	return nil
}

func (config *Config) SetRegion(str string) error {
	if err := config.ValidatorStr(str); err != nil {
		return err
	}
	config.Region = str
	return nil
}

func (config *Config) SetProjectKey(str string) error {
	if str != "" {
		if err := config.ValidatorStr(str); err != nil {
			return err
		}
	}

	config.ProjectId = str
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
