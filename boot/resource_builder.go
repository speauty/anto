package boot

import (
	_const "anto/const"
	_ "embed"
	"fmt"
	"os"
)

//go:embed favicon.ico
var bytesFav []byte

type ResourceBuilder struct {
}

func (customRB *ResourceBuilder) Install() {
	//customRB.genICO()
	customRB.genDirLogs()
	customRB.genCfgYml()
}

func (customRB *ResourceBuilder) genDirLogs() {
	dirname := "logs"
	fd, err := os.Stat(dirname)
	if err == nil && (fd != nil && fd.IsDir()) {
		return
	}
	if err := os.Mkdir(dirname, os.ModePerm); err != nil {
		panic(fmt.Errorf("创建日志目录失败(如果存在相应logs文件, 请手动处理), 错误: %s", err))
	}
}

func (customRB *ResourceBuilder) genCfgYml() {
	filename := "cfg.yml"
	_, err := os.Stat(filename)
	if err == nil {
		return
	}
	fd, err := os.Create(filename)
	if err != nil {
		panic(fmt.Errorf("创建YML配置文件失败, 错误: %s", err))
	}

	if _, err = fd.Write([]byte(_const.CfgYmlDefaultContent)); err != nil {
		panic(fmt.Errorf("写入YML默认配置失败, 错误: %s", err))
	}
}

func (customRB *ResourceBuilder) genICO() {
	filename := "favicon.ico"
	_, err := os.Stat(filename)
	if err == nil {
		return
	}
	fd, err := os.Create(filename)
	if err != nil {
		panic(fmt.Errorf("创建图标文件失败, 错误: %s", err))
	}

	if _, err = fd.Write(bytesFav); err != nil {
		panic(fmt.Errorf("写入图标文件失败, 错误: %s", err))
	}
}
