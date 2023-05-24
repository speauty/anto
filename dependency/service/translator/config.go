package translator

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"reflect"
)

const (
	ConfigInvalidStr string = "nil"
	ConfigInvalidInt int    = -1
)

// ImplConfig 引擎配置接口
type ImplConfig interface {
	AK() string // access-key or app-key or data-id and so on
	SK() string
	ProjectKey() string
	Region() string
	QPS() int
	MaxCharNum() int
	MaxCoroutineNum() int

	SetAK(ak string) error
	SetSK(sk string) error
	SetProjectKey(projectKey string) error
	SetRegion(region string) error
	SetQPS(num int) error
	SetMaxCharNum(num int) error
	SetMaxCoroutineNum(num int) error

	Default() ImplConfig               // 输出默认配置
	SyncDisk(viper *viper.Viper) error // 同步到磁盘
}

// DefaultConfig 默认配置结构体, 供具体引擎配置嵌入使用
type DefaultConfig struct{}

func (defaultConfig *DefaultConfig) AK() string { return ConfigInvalidStr }

func (defaultConfig *DefaultConfig) SK() string { return ConfigInvalidStr }

func (defaultConfig *DefaultConfig) ProjectKey() string { return ConfigInvalidStr }

func (defaultConfig *DefaultConfig) Region() string { return ConfigInvalidStr }

func (defaultConfig *DefaultConfig) QPS() int { return ConfigInvalidInt }

func (defaultConfig *DefaultConfig) MaxCharNum() int { return ConfigInvalidInt }

func (defaultConfig *DefaultConfig) MaxCoroutineNum() int { return ConfigInvalidInt }

func (defaultConfig *DefaultConfig) SetAK(_ string) error { return nil }

func (defaultConfig *DefaultConfig) SetSK(_ string) error { return nil }

func (defaultConfig *DefaultConfig) SetProjectKey(_ string) error { return nil }

func (defaultConfig *DefaultConfig) SetRegion(_ string) error { return nil }

func (defaultConfig *DefaultConfig) SetQPS(_ int) error { return nil }

func (defaultConfig *DefaultConfig) SetMaxCharNum(_ int) error { return nil }

func (defaultConfig *DefaultConfig) SetMaxCoroutineNum(_ int) error { return nil }

func (defaultConfig *DefaultConfig) Default() ImplConfig { return nil }

func (defaultConfig *DefaultConfig) SyncDisk(_ *viper.Viper) error {
	return errors.New("当前配置暂未实现磁盘同步方法")
}

func (defaultConfig *DefaultConfig) JoinAllTagAndValue(engine ImplTranslator, config ImplConfig, tagName string) map[string]interface{} {
	engineId := engine.GetId()
	configType := reflect.TypeOf(config)
	configVal := reflect.ValueOf(config)
	if configType.Kind() == reflect.Ptr { // 指针不支持
		configType = configType.Elem()
	}
	if configVal.Kind() == reflect.Ptr { // 指针不支持
		configVal = configVal.Elem()
	}

	result := make(map[string]interface{})

	for i := 0; i < configType.NumField(); i++ {
		currentField := configType.Field(i)
		// 当前仅支持整型和字符串
		if currentField.Type.Kind() != reflect.Int &&
			currentField.Type.Kind() != reflect.String {
			continue
		}
		tagVal := currentField.Tag.Get(tagName)
		if tagVal == "" || tagVal == "-" {
			continue
		}
		// @todo 可以直接在这里IO的, 但是想了一下还是交给具体配置处理, 毕竟功能还是要分开, 该方法只负责联合标签和具体值
		result[fmt.Sprintf("%s.%s", engineId, tagVal)] = configVal.Field(i).Interface()
	}
	return result
}
