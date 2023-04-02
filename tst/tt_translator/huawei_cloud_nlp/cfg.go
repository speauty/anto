package huawei_cloud_nlp

type Cfg struct {
	AKId      string `mapstructure:"ak_id"`      // Access Key
	SkKey     string `mapstructure:"sk_key"`     // Secret Access Key
	Region    string `mapstructure:"region"`     // 当前接口开发的区域, 目前仅支持华北-北京四终端节点
	ProjectId string `mapstructure:"project_id"` // 项目ID
}

func (customC Cfg) Default() *Cfg {
	return &Cfg{AKId: "", SkKey: "", Region: "cn-north-4", ProjectId: ""}
}
