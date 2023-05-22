package e_fyne

import "fyne.io/fyne/v2"

type ImplPage interface {
	GetID() string
	GetName() string
	IsDefault() bool
	OnClose()
	OnReset()
	OnRender() fyne.CanvasObject
}
