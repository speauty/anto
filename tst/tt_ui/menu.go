package tt_ui

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var actionStatusBarHandle *walk.Action

func (customG *Gui) defaultMenu() []MenuItem {
	return []MenuItem{
		Menu{
			Text: "文件",
			Items: []MenuItem{
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
	_ = customG.GetWindow().Close()
}
