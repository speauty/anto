package baidu_cloud_mt

import (
	"anto/domain/service/translator"

	"github.com/spf13/viper"
)

type Config struct {
	*translator.DefaultConfig
	AppID           string `mapstructure:"app_id"`
	AppSecret       string `mapstructure:"app_secret"`
	ServiceTier     string `mapstructure:"service_tier"`
	MaxCoroutineNum int    `mapstructure:"max_coroutine_num"`
}

func (config *Config) Default() translator.ImplConfig {
	return &Config{AppID: "", AppSecret: "", ServiceTier: TierStandard, MaxCoroutineNum: 1}
}

const (
	TierStandard = "standard"
	TierAdvanced = "advanced"
	TierPremium  = "premium"
)

// TierLimits returns the request character ceiling and QPS for each plan.
// The 6000-byte UTF-8 ceiling is enforced separately by splitText.
func TierLimits(tier string) (maxChars, qps int, valid bool) {
	switch tier {
	case TierStandard:
		return 1000, 1, true
	case TierAdvanced:
		return 6000, 10, true
	case TierPremium:
		return 6000, 100, true
	default:
		return 0, 0, false
	}
}

func (config *Config) SyncDisk(currentViper *viper.Viper) error {
	for tag, value := range config.JoinAllTagAndValue(API(), config, "mapstructure") {
		currentViper.Set(tag, value)
	}
	return nil
}

func (config *Config) GetAK() string { return config.AppID }
func (config *Config) GetSK() string { return config.AppSecret }
func (config *Config) GetQPS() int {
	_, qps, valid := TierLimits(config.ServiceTier)
	if !valid {
		return 1
	}
	return qps
}
func (config *Config) GetMaxCharNum() int {
	maxChars, _, valid := TierLimits(config.ServiceTier)
	if !valid {
		return 1000
	}
	return maxChars
}
func (config *Config) GetMaxCoroutineNum() int { return config.MaxCoroutineNum }

func (config *Config) SetAK(value string) error {
	if err := config.ValidatorStr(value); err != nil {
		return err
	}
	config.AppID = value
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
	return &tierControlledLimitError{}
}
func (config *Config) SetMaxCharNum(value int) error {
	return &tierControlledLimitError{}
}
func (config *Config) SetMaxCoroutineNum(value int) error {
	if err := config.ValidatorNum(value); err != nil {
		return err
	}
	config.MaxCoroutineNum = value
	return nil
}

type tierControlledLimitError struct{}

func (*tierControlledLimitError) Error() string {
	return "百度云请求长度和 QPS 由服务版本决定"
}
