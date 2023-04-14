package ui

import (
	"github.com/lxn/walk"
	"github.com/lxn/win"
)

func (customG *Gui) setWindowFlag() {
	win.SetWindowLong(customG.GetWindow().Handle(), win.GWL_STYLE,
		win.GetWindowLong(customG.GetWindow().Handle(), win.GWL_STYLE) & ^win.WS_MAXIMIZEBOX & ^win.WS_THICKFRAME)
}

func (customG *Gui) setWindowCenter() {
	scrWidth := win.GetSystemMetrics(win.SM_CXSCREEN)
	scrHeight := win.GetSystemMetrics(win.SM_CYSCREEN)
	_ = customG.GetWindow().SetBounds(walk.Rectangle{
		X: int((scrWidth - width) / 2), Y: int((scrHeight - height) / 2),
		Width: width, Height: height,
	})
}

func (customG *Gui) setWindowMinAndMaxSize() {
	minMaxSize := walk.Size{Width: width, Height: height}
	_ = customG.GetWindow().SetMinMaxSize(minMaxSize, minMaxSize)
}
