package e_fyne

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func (ag *AppGui) ChanConsole() chan string {
	return ag.chanConsole
}

func (ag *AppGui) PushToConsole(msg string) {
	ag.chanConsole <- msg
}

func (ag *AppGui) initConsoleWindow() {
	if ag.consoleWindow != nil {
		return
	}

	consoleWindow := ag.app.NewWindow(fmt.Sprintf("%s-控制台", ag.config.AppName))
	consoleWindow.Resize(fyne.NewSize(400, ag.config.AppMainWindowDefaultHeight))
	//consoleWindow.SetFixedSize(true)
	consoleWindow.SetCloseIntercept(consoleWindow.Hide)

	ag.chanConsole = make(chan string, appChanConsoleCnt)
	ag.consoleWindow = consoleWindow
	ag.consoleWindow.Hide()
	ag.refresherConsoleWindow()
}

func (ag *AppGui) refresherConsoleWindow() {
	go func() {
		for msgConsole := range ag.chanConsole {
			ag.listConsole = append(ag.listConsole, msgConsole)
			bl := widget.NewList(func() int {
				return len(ag.listConsole)
			}, func() fyne.CanvasObject {
				return widget.NewLabel("")
			}, func(listItemID widget.ListItemID, renderObj fyne.CanvasObject) {
				renderObj.(*widget.Label).Wrapping = fyne.TextWrapBreak
				renderObj.(*widget.Label).SetText(ag.listConsole[listItemID])
			})
			ag.consoleWindow.SetContent(bl)
		}
	}()
}

func (ag *AppGui) ClearConsole() {
	ag.listConsole = []string{}
	ag.consoleWindow.SetContent(widget.NewLabel(""))
}
