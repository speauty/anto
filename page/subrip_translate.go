package page

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"sync"
	"translator/tst/tt_ui/pack"
	"translator/util"
)

var (
	apiSubripTranslate  *SubripTranslate
	onceSubripTranslate sync.Once
)

func GetSubripTranslate() *SubripTranslate {
	onceSubripTranslate.Do(func() {
		apiSubripTranslate = new(SubripTranslate)
		apiSubripTranslate.id = util.Uid()
		apiSubripTranslate.name = "字幕翻译"
	})
	return apiSubripTranslate
}

type SubripTranslate struct {
	id         string
	name       string
	mainWindow *walk.MainWindow
	rootWidget *walk.Composite
}

func (customPage *SubripTranslate) GetId() string {
	return customPage.id
}

func (customPage *SubripTranslate) GetName() string {
	return customPage.name
}

func (customPage *SubripTranslate) BindWindow(win *walk.MainWindow) {
	customPage.mainWindow = win
}

func (customPage *SubripTranslate) SetVisible(isVisible bool) {
	if customPage.rootWidget != nil {
		customPage.rootWidget.SetVisible(isVisible)
	}
}

func (customPage *SubripTranslate) GetWidget() Widget {
	return StdRootWidget(&customPage.rootWidget,
		pack.TTLabel(pack.NewTTLabelArgs(nil).SetText("神秘代码：字幕翻译")),
	)
}

func (customPage *SubripTranslate) Reset() {

}
