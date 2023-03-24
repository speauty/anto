package tt_ui

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"translator/tst/tt_ui/msg"
)

var actionStatusBarHandle *walk.Action

func (customG *Gui) defaultMenu() []MenuItem {
	return []MenuItem{
		Menu{
			Text: "文件",
			Items: []MenuItem{
				Action{
					Text: "首选项",
				},
				Separator{},
				Action{
					AssignTo:    &actionStatusBarHandle,
					Text:        "状态栏",
					Checked:     true,
					OnTriggered: customG.eventActionStatusBar,
				},
				Separator{},
				Action{
					Text:        "退出",
					OnTriggered: customG.eventActionQuit,
				},
			},
		},
		Action{
			Text: "字幕翻译",
		},
		Menu{
			Text: "帮助",
			Items: []MenuItem{
				Action{
					Text: "使用手册",
				},
				Action{
					Text: "关于我们",
				},
			},
		},
	}
}

func (customG *Gui) eventActionStatusBar() {
	customG.GetWindow().StatusBar().SetVisible(!customG.GetWindow().StatusBar().Visible())
	if actionStatusBarHandle != nil {
		_ = actionStatusBarHandle.SetChecked(customG.GetWindow().StatusBar().Visible())
	}
}

func (customG *Gui) eventActionQuit() {
	isOk, _ := msg.Confirm(customG.GetWindow(), fmt.Sprintf("即将退出当前应用，是否确认？"))
	if isOk {
		_ = customG.GetWindow().Close()
	}
}
