package ai_baidu

import (
	"anto/domain/service/translator"
	"encoding/json"
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/spf13/viper"
	"io"
	"time"
)

type Config struct {
	*translator.DefaultConfig
	ApiKey          string `mapstructure:"api_key"`    // 应用关键字
	SecretKey       string `mapstructure:"secret_key"` // 应用密钥
	QPS             int    `mapstructure:"qps"`
	MaxCharNum      int    `mapstructure:"max_single_text_length"`
	MaxCoroutineNum int    `mapstructure:"max_coroutine_num"`
	AccessToken     string `mapstructure:"access_token"`
	ExpiredAt       int64  `mapstructure:"expired_at"`
}

func (config *Config) Default() translator.ImplConfig {
	return &Config{
		ApiKey: "", SecretKey: "",
		MaxCharNum: 6000, QPS: 10, MaxCoroutineNum: 5,
	}
}

func (config *Config) SyncDisk(currentViper *viper.Viper) error {
	tagAndVal := config.JoinAllTagAndValue(API(), config, "mapstructure")

	for tag, val := range tagAndVal {
		currentViper.Set(tag, val)
	}
	return nil
}

func (config *Config) GetAK() string { return config.ApiKey }
func (config *Config) GetSK() string { return config.SecretKey }

func (config *Config) GetQPS() int             { return config.QPS }
func (config *Config) GetMaxCharNum() int      { return config.MaxCharNum }
func (config *Config) GetMaxCoroutineNum() int { return config.MaxCoroutineNum }

func (config *Config) SetAK(str string) error {
	if err := config.ValidatorStr(str); err != nil {
		return err
	}
	config.ApiKey = str
	return nil
}
func (config *Config) SetSK(str string) error {
	if err := config.ValidatorStr(str); err != nil {
		return err
	}
	config.SecretKey = str
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

func (config *Config) GetAccessToken() (string, error) {
	if config.AccessToken == "" || config.ExpiredAt < (time.Now().Unix()-3600) {
		tmpUrl := fmt.Sprintf(
			"https://aip.baidubce.com/oauth/2.0/token?client_id=%s&client_secret=%s&grant_type=client_credentials",
			config.GetAK(), config.GetSK(),
		)
		resp, err := req.R().Post(tmpUrl)
		if err != nil {
			return "", err
		}
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("读取报文异常, 错误: %s", err.Error())
		}
		auth := &authResp{}
		if err = json.Unmarshal(respBytes, auth); err != nil {
			return "", fmt.Errorf("解析报文异常, 错误: %s", err.Error())
		}
		config.AccessToken = auth.AccessToken
		config.ExpiredAt = time.Now().Unix() + auth.ExpiresIn
	}

	return config.AccessToken, nil
}

type authResp struct {
	RefreshToken  string `json:"refresh_token"`
	ExpiresIn     int64  `json:"expires_in"`
	SessionKey    string `json:"session_key"`
	AccessToken   string `json:"access_token"`
	Scope         string `json:"scope"`
	SessionSecret string `json:"session_secret"`
}
