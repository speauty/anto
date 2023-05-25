package win

import (
	"anto/common"
	"anto/lib/log"
	page2 "anto/platform/win/page"
	"anto/platform/win/ui"
	"anto/platform/win/ui/msg"
	"errors"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"go.uber.org/zap"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
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
	mainWindow           *walk.MainWindow
	actionDownloadHandle *walk.Action
}

func (customM *TTMenu) GetMenus() []MenuItem {
	return []MenuItem{
		Menu{
			Text: "文件",
			Items: []MenuItem{
				Action{Text: "设置", OnTriggered: customM.eventSettings},
				Separator{},
				Action{AssignTo: &customM.actionDownloadHandle, Text: "下载新版", OnTriggered: customM.eventActionDownloadLatestVersion},
				Action{Text: "清除日志", OnTriggered: customM.eventActionDelLog},
				Separator{},
				Action{Text: "退出", OnTriggered: customM.eventActionQuit},
			},
		},
		Action{Text: "字幕翻译", OnTriggered: customM.eventSubtitleTranslate},
		Action{Text: "关于我们", OnTriggered: customM.eventActionAboutUS},
	}
}

func (customM *TTMenu) eventSettings() {
	currentPage := page2.GetSettings()
	customM.eventGoPage(currentPage.GetId(), currentPage.GetName())
}

func (customM *TTMenu) eventActionDownloadLatestVersion() {
	if customM.actionDownloadHandle.Enabled() {
		_ = customM.actionDownloadHandle.SetEnabled(false)
	} else {
		return
	}
	mainWindow := ui.Singleton().GetWindow()
	isOk, _ := msg.Confirm(mainWindow, fmt.Sprintf("下载应用的最新版本，是否继续？"))
	if isOk {
		go func() {
			defer func() {
				_ = customM.actionDownloadHandle.SetEnabled(true)
			}()
			resp, err := http.Get(common.DownloadLatestVersionUrl)
			if err != nil {
				msg.Err(mainWindow, fmt.Errorf("下载最新版本异常, 错误: %s", err))
				return
			}
			defer func() {
				_ = resp.Body.Close()
			}()
			appBytes, _ := io.ReadAll(resp.Body)
			if resp.StatusCode == http.StatusNotFound {
				msg.Err(mainWindow, errors.New("下载最新版本异常, 错误: 暂未找到"))
				return
			}
			fileName := filepath.Base(common.DownloadLatestVersionUrl)

			if err := os.WriteFile(fileName, appBytes, os.ModePerm); err != nil {
				msg.Err(mainWindow, fmt.Errorf("下载最新版本异常, 错误: %s", err))
				return
			}
			msg.Info(mainWindow, fmt.Sprintf("下载最新版本成功[%s], 关闭当前应用, 双击打开对应可执行文件即可", fileName))
		}()
	} else {
		_ = customM.actionDownloadHandle.SetEnabled(true)
	}
}

func (customM *TTMenu) eventActionDelLog() {
	mainWindow := ui.Singleton().GetWindow()
	isOk, _ := msg.Confirm(mainWindow, fmt.Sprintf("清除日志会删除今日之前的所有日志文件，是否继续？"))
	if isOk {
		now := time.Now()
		currentDayZero := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
		cntDel := 0
		_ = filepath.Walk("logs", func(path string, info fs.FileInfo, err error) error {
			if info.ModTime().Before(currentDayZero) {
				if err := os.Remove(path); err == nil {
					cntDel++
				}
			}
			return nil
		})
		if cntDel > 0 {
			msg.Info(ui.Singleton().GetWindow(), fmt.Sprintf("删除历史日志文件(数量: %d)", cntDel))
		} else {
			msg.Info(ui.Singleton().GetWindow(), "暂无历史日志文件")
		}
	}
}

func (customM *TTMenu) eventActionQuit() {
	mainWindow := ui.Singleton().GetWindow()
	isOk, _ := msg.Confirm(mainWindow, fmt.Sprintf("即将退出当前应用，是否确认？"))
	if isOk {
		_ = mainWindow.Close()
	}
}

func (customM *TTMenu) eventActionAboutUS() {
	currentPage := page2.GetAboutUs()
	customM.eventGoPage(currentPage.GetId(), currentPage.GetName())
}

func (customM *TTMenu) eventSubtitleTranslate() {
	currentPage := page2.GetSubripTranslate()
	customM.eventGoPage(currentPage.GetId(), currentPage.GetName())
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
