package handle

import (
	"anto/platform/win/ui/msg"
	"fmt"
	"github.com/lxn/walk"
)

func FileDialogHandle(args *FileDialogHandleArgs) {
	dlg := new(walk.FileDialog)
	dlg.Title = args.title
	dlg.Filter = args.filter
	topic := ""
	var isAccepted bool
	var err error
	if args.isFileSelector {
		topic = "文件"
		isAccepted, err = dlg.ShowOpen(args.owner)
	} else {
		topic = "目录"
		isAccepted, err = dlg.ShowBrowseFolder(args.owner)
	}
	if err != nil {
		msg.Err(args.owner, fmt.Errorf("选择%s异常, 错误: %s", topic, err))
		return
	} else if !isAccepted {
		if (*args.pathEchoHandle).Text() == "" {
			msg.Err(args.owner, fmt.Errorf("选择%s失败, 请重新选择", topic))
		}
		return
	}
	_ = (*args.pathEchoHandle).SetText(dlg.FilePath)
}

type FileDialogHandleArgs struct {
	owner          walk.Form
	pathEchoHandle **walk.Label
	title          string
	filter         string
	isFileSelector bool
}

func NewFileDialogHandleArgs(owner walk.Form, pathEchoHandle **walk.Label) *FileDialogHandleArgs {
	return &FileDialogHandleArgs{owner: owner, pathEchoHandle: pathEchoHandle, isFileSelector: true}
}

func (customT *FileDialogHandleArgs) File() *FileDialogHandleArgs {
	customT.isFileSelector = true
	return customT
}

func (customT *FileDialogHandleArgs) Folder() *FileDialogHandleArgs {
	customT.isFileSelector = false
	return customT
}

func (customT *FileDialogHandleArgs) SetOwner(owner walk.Form) *FileDialogHandleArgs {
	customT.owner = owner
	return customT
}

func (customT *FileDialogHandleArgs) SetTitle(title string) *FileDialogHandleArgs {
	customT.title = title
	return customT
}

func (customT *FileDialogHandleArgs) SetFilter(filter string) *FileDialogHandleArgs {
	customT.filter = filter
	return customT
}

func (customT *FileDialogHandleArgs) SetPathEchoHandle(pathEchoHandle **walk.Label) *FileDialogHandleArgs {
	customT.pathEchoHandle = pathEchoHandle
	return customT
}
