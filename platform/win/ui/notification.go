package ui

import (
	"anto/platform/win/ui/msg"
	"errors"
	"github.com/lxn/walk"
)

func (customG *Gui) initNotification() {
	if customG.notification == nil {
		tmpNotification, err := walk.NewNotifyIcon(customG.GetWindow())
		if err != nil {
			msg.Err(customG.GetWindow(), errors.New("初始化系统通知失败"))
			return
		}
		customG.notification = tmpNotification
		//defer tmpNotification.Dispose()
	}

	_ = customG.notification.SetVisible(true)
}

func (customG *Gui) Notification() *walk.NotifyIcon {
	return customG.notification
}
