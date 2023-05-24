package bootstrap

import (
	_const "anto/common"
	"anto/resource"
	"fmt"
	"os"
)

type ResourceBuilder struct {
}

func (customRB *ResourceBuilder) Install() {
	customRB.installICO()
	customRB.installDirLogs()
	customRB.installCfgYml()
}

func (customRB *ResourceBuilder) installDirLogs() {
	dirname := "logs"
	fd, err := os.Stat(dirname)
	if err == nil && (fd != nil && fd.IsDir()) {
		return
	}
	if err := os.Mkdir(dirname, os.ModePerm); err != nil {
		panic(fmt.Errorf("创建日志目录失败(如果存在相应logs文件, 请手动处理), 错误: %s", err))
	}
}

func (customRB *ResourceBuilder) installCfgYml() {
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

func (customRB *ResourceBuilder) installICO() {
	filename := "favicon.ico"
	_, err := os.Stat(filename)
	if err == nil {
		return
	}
	fd, err := os.Create(filename)
	if err != nil {
		panic(fmt.Errorf("创建图标文件失败, 错误: %s", err))
	}

	if _, err = fd.Write(resource.Favicon); err != nil {
		panic(fmt.Errorf("写入图标文件失败, 错误: %s", err))
	}
}
