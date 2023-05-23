package pages

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"sync"
)

var (
	apiPageSubtitleTranslate  *PageSubtitleTranslate
	oncePageSubtitleTranslate sync.Once
)

func APIPageSubtitleTranslate() *PageSubtitleTranslate {
	oncePageSubtitleTranslate.Do(func() {
		apiPageSubtitleTranslate = &PageSubtitleTranslate{
			id:   "page.subtitle.translate",
			name: "字幕翻译",
			//isDefault: true,
		}
	})
	return apiPageSubtitleTranslate
}

type PageSubtitleTranslate struct {
	window    fyne.Window
	id        string
	name      string
	isDefault bool
}

func (page *PageSubtitleTranslate) GetID() string { return page.id }

func (page *PageSubtitleTranslate) GetName() string { return page.name }

func (page *PageSubtitleTranslate) GetWindow() fyne.Window { return page.window }

func (page *PageSubtitleTranslate) SetWindow(win fyne.Window) { page.window = win }

func (page *PageSubtitleTranslate) IsDefault() bool { return page.isDefault }

func (page *PageSubtitleTranslate) OnClose() {}

func (page *PageSubtitleTranslate) OnReset() {}

func (page *PageSubtitleTranslate) OnRender() fyne.CanvasObject {
	return container.NewCenter(
		widget.NewLabel(page.GetName()),
	)
}
