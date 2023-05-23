package msg

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func Info(window fyne.Window, content, title string, closeFn func()) {
	if title == "" {
		title = "友情提示"
	}
	if content == "" {
		content = "试试，操作成功没？"
	}
	info := dialog.NewInformation(title, content, window)
	info.SetDismissText("关闭")
	if closeFn != nil {
		info.SetOnClosed(closeFn)
	}
	info.Show()
}

func Confirm(window fyne.Window, content, title string, callback func(bool)) {
	if title == "" {
		title = "友情提示"
	}
	if content == "" {
		content = "需要执行该操作吗"
	}
	confirm := dialog.NewConfirm(title, content, callback, window)
	confirm.SetConfirmText("确认")
	confirm.SetDismissText("取消")
	confirm.Show()
}

func Error(window fyne.Window, err error, node string) {
	info := dialog.NewError(err, window)
	info.SetDismissText("关闭")
	info.Show()
}
