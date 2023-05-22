package pages

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func pageTitle(title string) fyne.CanvasObject {
	return container.NewCenter(widget.NewLabel(title))
}
