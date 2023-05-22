package huawei_cloud_nlp

import "anto/domain/service/translator"

type Cfg struct {
	*translator.DefaultConfig
	AKId                string `mapstructure:"ak_id"`      // Access Key
	SkKey               string `mapstructure:"sk_key"`     // Secret Access Key
	Region              string `mapstructure:"region"`     // 当前接口开发的区域, 目前仅支持华北-北京四终端节点
	ProjectId           string `mapstructure:"project_id"` // 项目ID
	MaxSingleTextLength int    `mapstructure:"max_single_text_length"`
}

func (customC *Cfg) GetAK() string     { return customC.AKId }
func (customC *Cfg) GetSK() string     { return customC.SkKey }
func (customC *Cfg) GetPK() string     { return customC.ProjectId }
func (customC *Cfg) GetTML() int       { return customC.MaxSingleTextLength }
func (customC *Cfg) GetRegion() string { return customC.Region }

func (customC *Cfg) Default() translator.ImplConfig {
	return &Cfg{AKId: "", SkKey: "", Region: "cn-north-4", ProjectId: "", MaxSingleTextLength: 2000}
}

func (customC *Cfg) SetAK(ak string) error {
	customC.AKId = ak
	return nil
}

func (customC *Cfg) SetSK(sk string) error {
	customC.SkKey = sk
	return nil
}

func (customC *Cfg) SetPK(pk string) error {
	customC.ProjectId = pk
	return nil
}

func (customC *Cfg) SetRegion(region string) error {
	customC.Region = region
	return nil
}

func (customC *Cfg) SetTML(tml int) error {
	customC.MaxSingleTextLength = tml
	return nil
}
