package cross_platform_fyne

import "fyne.io/fyne/v2"

type ImplPage interface {
	GetID() string
	GetName() string
	GetWindow() fyne.Window
	SetWindow(win fyne.Window)
	IsDefault() bool
	OnClose()
	OnReset()
	OnRender() fyne.CanvasObject
}
