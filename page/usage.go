package page

import (
	"anto/tst/tt_ui/pack"
	"anto/util"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"sync"
)

var (
	apiUsage  *Usage
	onceUsage sync.Once
)

func GetUsage() *Usage {
	onceUsage.Do(func() {
		apiUsage = new(Usage)
		apiUsage.id = util.Uid()
		apiUsage.name = "使用手册"
	})
	return apiUsage
}

type Usage struct {
	id         string
	name       string
	mainWindow *walk.MainWindow
	rootWidget *walk.Composite
}

func (customPage *Usage) GetId() string {
	return customPage.id
}

func (customPage *Usage) GetName() string {
	return customPage.name
}

func (customPage *Usage) BindWindow(win *walk.MainWindow) {
	customPage.mainWindow = win
}

func (customPage *Usage) SetVisible(isVisible bool) {
	if customPage.rootWidget != nil {
		customPage.rootWidget.SetVisible(isVisible)
	}
}

func (customPage *Usage) GetWidget() Widget {
	return StdRootWidget(&customPage.rootWidget,
		pack.TTLabel(pack.NewTTLabelArgs(nil).SetText("神秘代码：使用手册")),
	)
}

func (customPage *Usage) Reset() {

}
