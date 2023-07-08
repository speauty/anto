package tencent_cloud_mt

import (
	"anto/domain/service/translator"
	"github.com/spf13/viper"
	"strconv"
)

type Config struct {
	*translator.DefaultConfig
	SecretId        string `mapstructure:"secret_id"`  // 用于标识接口调用者身份
	SecretKey       string `mapstructure:"secret_key"` // 用于验证接口调用者的身份
	Region          string `mapstructure:"region"`     // 地域参数
	ProjectId       int64  `mapstructure:"project_id"` // 项目ID, 默认值为0
	QPS             int    `mapstructure:"qps"`
	MaxCharNum      int    `mapstructure:"max_single_text_length"`
	MaxCoroutineNum int    `mapstructure:"max_coroutine_num"`
}

func (config *Config) Default() translator.ImplConfig {
	return &Config{
		SecretId: "", SecretKey: "", Region: "ap-chengdu", ProjectId: 0,
		MaxCharNum: 2000, QPS: 5, MaxCoroutineNum: 4,
	}
}

func (config *Config) SyncDisk(currentViper *viper.Viper) error {
	tagAndVal := config.JoinAllTagAndValue(API(), config, "mapstructure")

	for tag, val := range tagAndVal {
		currentViper.Set(tag, val)
	}
	return nil
}

func (config *Config) GetAK() string            { return config.SecretId }
func (config *Config) GetSK() string            { return config.SecretKey }
func (config *Config) GetRegion() string        { return config.Region }
func (config *Config) GetProjectKey() string    { return strconv.Itoa(int(config.ProjectId)) }
func (config *Config) GetProjectKeyPtr() *int64 { return &config.ProjectId }

func (config *Config) GetQPS() int             { return config.QPS }
func (config *Config) GetMaxCharNum() int      { return config.MaxCharNum }
func (config *Config) GetMaxCoroutineNum() int { return config.MaxCoroutineNum }

func (config *Config) SetAK(str string) error {
	if err := config.ValidatorStr(str); err != nil {
		return err
	}
	config.SecretId = str
	return nil
}

func (config *Config) SetSK(str string) error {
	if err := config.ValidatorStr(str); err != nil {
		return err
	}
	config.SecretKey = str
	return nil
}

func (config *Config) SetProjectKey(projectKey string) error {
	if err := config.ValidatorStr(projectKey); err != nil {
		return err
	}
	tmpPK, err := strconv.Atoi(projectKey)
	if err != nil {
		return err
	}
	config.ProjectId = int64(tmpPK)
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
