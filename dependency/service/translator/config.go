package translator

import "errors"

const (
	ConfigInvalidStr string = "nil"
	ConfigInvalidInt int    = -1
)

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

	SyncDisk() error // 同步到磁盘
}

type DefaultConfig struct{}

func (defaultConfig *DefaultConfig) AK() string { return ConfigInvalidStr }

func (defaultConfig *DefaultConfig) SK() string { return ConfigInvalidStr }

func (defaultConfig *DefaultConfig) ProjectKey() string { return ConfigInvalidStr }

func (defaultConfig *DefaultConfig) Region() string { return ConfigInvalidStr }

func (defaultConfig *DefaultConfig) QPS() int { return ConfigInvalidInt }

func (defaultConfig *DefaultConfig) MaxCharNum() int { return ConfigInvalidInt }

func (defaultConfig *DefaultConfig) MaxCoroutineNum() int { return ConfigInvalidInt }

func (defaultConfig *DefaultConfig) SetAK(ak string) error { return nil }

func (defaultConfig *DefaultConfig) SetSK(sk string) error { return nil }

func (defaultConfig *DefaultConfig) SetProjectKey(projectKey string) error { return nil }

func (defaultConfig *DefaultConfig) SetRegion(region string) error { return nil }

func (defaultConfig *DefaultConfig) SetQPS(num int) error { return nil }

func (defaultConfig *DefaultConfig) SetMaxCharNum(num int) error { return nil }

func (defaultConfig *DefaultConfig) SetMaxCoroutineNum(num int) error { return nil }

func (defaultConfig *DefaultConfig) SyncDisk() error {
	return errors.New("当前配置暂未实现磁盘同步方法")
}
