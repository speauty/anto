package menu

import (
	"anto/lib/log"
	"anto/lib/ui"
	"anto/lib/ui/msg"
	"anto/page"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"go.uber.org/zap"
	"sync"
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
				currentPage := page.GetSettings()
				customM.eventGoPage(currentPage.GetId(), currentPage.GetName())
			},
		},
		Action{
			Text: "字幕翻译",
			OnTriggered: func() {
				currentPage := page.GetSubripTranslate()
				customM.eventGoPage(currentPage.GetId(), currentPage.GetName())
			},
		},
		Menu{
			Text: "帮助",
			Items: []MenuItem{
				Action{
					Text: "使用手册",
					OnTriggered: func() {
						currentPage := page.GetUsage()
						customM.eventGoPage(currentPage.GetId(), currentPage.GetName())
					},
				},
				Action{
					Text: "关于我们",
					OnTriggered: func() {
						currentPage := page.GetAboutUs()
						customM.eventGoPage(currentPage.GetId(), currentPage.GetName())
					},
				},
			},
		},
	}
}

func (customM *TTMenu) eventActionStatusBar() {
	mainWindow := ui.GetInstance().GetWindow()
	mainWindow.StatusBar().SetVisible(!mainWindow.StatusBar().Visible())
	if customM.actionStatusBarHandle != nil {
		_ = customM.actionStatusBarHandle.SetChecked(mainWindow.StatusBar().Visible())
	}
}

func (customM *TTMenu) eventActionQuit() {
	mainWindow := ui.GetInstance().GetWindow()
	isOk, _ := msg.Confirm(mainWindow, fmt.Sprintf("即将退出当前应用，是否确认？"))
	if isOk {
		_ = mainWindow.Close()
	}
}

func (customM *TTMenu) eventGoPage(pageId string, name string) {
	if pageId == "" {
		return
	}
	if err := ui.GetInstance().GoPage(pageId); err != nil {
		log.Singleton().Error("跳转页面异常", zap.String("page", name), zap.String("id", pageId), zap.Error(err))
		msg.Err(ui.GetInstance().GetWindow(), fmt.Errorf("跳转页面[%s]异常", name))
	}
}
