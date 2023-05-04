package win

import (
	"anto/lib/log"
	page2 "anto/platform/win/page"
	"anto/platform/win/ui"
	"anto/platform/win/ui/msg"
	"fmt"
	"sync"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"go.uber.org/zap"
)

var (
	apiTTMenu  *TTMenu
	onceTTMenu sync.Once
)

func GetInstance() *TTMenu {
	onceTTMenu.Do(func() {
		apiTTMenu = new(TTMenu)
	})
	return apiTTMenu
}

type TTMenu struct {
	mainWindow            *walk.MainWindow
	actionStatusBarHandle *walk.Action
}

func (customM *TTMenu) GetMenus() []MenuItem {
	return []MenuItem{
		/*Menu{
			Text: "文件",
			Items: []MenuItem{
				Action{
					Text: "设置",
					OnTriggered: func() {
						currentPage := page.GetSettings()
						customM.eventGoPage(currentPage.GetId(), currentPage.GetName())
					},
				},
				Separator{},
				Action{
					AssignTo:    &customM.actionStatusBarHandle,
					Text:        "状态栏",
					Checked:     true,
					OnTriggered: customM.eventActionStatusBar,
				},
				Separator{},
				Action{
					Text:        "退出",
					OnTriggered: customM.eventActionQuit,
				},
			},
		},*/
		Action{
			Text: "设置",
			OnTriggered: func() {
				currentPage := page2.GetSettings()
				customM.eventGoPage(currentPage.GetId(), currentPage.GetName())
			},
		},
		Action{
			Text: "字幕翻译",
			OnTriggered: func() {
				currentPage := page2.GetSubripTranslate()
				customM.eventGoPage(currentPage.GetId(), currentPage.GetName())
			},
		},
		Action{
			Text: "关于我们",
			OnTriggered: func() {
				currentPage := page2.GetAboutUs()
				customM.eventGoPage(currentPage.GetId(), currentPage.GetName())
			},
		},
		Menu{
			Text:    "帮助",
			Visible: false,
			Items: []MenuItem{
				Action{
					Text: "使用手册",
					OnTriggered: func() {
						currentPage := page2.GetUsage()
						customM.eventGoPage(currentPage.GetId(), currentPage.GetName())
					},
				},
				Action{
					Text: "关于我们",
					OnTriggered: func() {
						currentPage := page2.GetAboutUs()
						customM.eventGoPage(currentPage.GetId(), currentPage.GetName())
					},
				},
			},
		},
	}
}

func (customM *TTMenu) eventActionStatusBar() {
	mainWindow := ui.Singleton().GetWindow()
	mainWindow.StatusBar().SetVisible(!mainWindow.StatusBar().Visible())
	if customM.actionStatusBarHandle != nil {
		_ = customM.actionStatusBarHandle.SetChecked(mainWindow.StatusBar().Visible())
	}
}

func (customM *TTMenu) eventActionQuit() {
	mainWindow := ui.Singleton().GetWindow()
	isOk, _ := msg.Confirm(mainWindow, fmt.Sprintf("即将退出当前应用，是否确认？"))
	if isOk {
		_ = mainWindow.Close()
	}
}

func (customM *TTMenu) eventGoPage(pageId string, name string) {
	if pageId == "" {
		return
	}
	if err := ui.Singleton().GoPage(pageId); err != nil {
		log.Singleton().Error("跳转页面异常", zap.String("page", name), zap.String("id", pageId), zap.Error(err))
		msg.Err(ui.Singleton().GetWindow(), fmt.Errorf("跳转页面[%s]异常", name))
	}
}
