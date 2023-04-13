package page

import (
	"anto/tst/tt_ui/pack"
	"anto/util"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"sync"
)

var (
	apiTpl  *Tpl
	onceTpl sync.Once
)

func GetTpl() *Tpl {
	onceTpl.Do(func() {
		apiTpl = new(Tpl)
		apiTpl.id = util.Uid()
		apiTpl.name = "模板页"
	})
	return apiTpl
}

type Tpl struct {
	id         string
	name       string
	mainWindow *walk.MainWindow
	rootWidget *walk.Composite
}

func (customPage *Tpl) GetId() string {
	return customPage.id
}

func (customPage *Tpl) GetName() string {
	return customPage.name
}

func (customPage *Tpl) BindWindow(win *walk.MainWindow) {
	customPage.mainWindow = win
}

func (customPage *Tpl) SetVisible(isVisible bool) {
	if customPage.rootWidget != nil {
		customPage.rootWidget.SetVisible(isVisible)
	}
}

func (customPage *Tpl) GetWidget() Widget {
	return StdRootWidget(&customPage.rootWidget,
		pack.TTLabel(pack.NewTTLabelArgs(nil).SetText("神秘代码：MLGJB")),
	)
}

func (customPage *Tpl) Reset() {

}
