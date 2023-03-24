package page

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"sync"
	"translator/tst/tt_ui/pack"
	"translator/util"
)

var (
	apiSettings  *Settings
	onceSettings sync.Once
)

func GetSettings() *Settings {
	onceSettings.Do(func() {
		apiSettings = new(Settings)
		apiSettings.id = util.Uid()
		apiSettings.name = "模板页"
	})
	return apiSettings
}

type Settings struct {
	id         string
	name       string
	mainWindow *walk.MainWindow
	rootWidget *walk.Composite
}

func (customPage *Settings) GetId() string {
	return customPage.id
}

func (customPage *Settings) GetName() string {
	return customPage.name
}

func (customPage *Settings) BindWindow(win *walk.MainWindow) {
	customPage.mainWindow = win
}

func (customPage *Settings) SetVisible(isVisible bool) {
	if customPage.rootWidget != nil {
		customPage.rootWidget.SetVisible(isVisible)
	}
}

func (customPage *Settings) GetWidget() Widget {
	return StdRootWidget(&customPage.rootWidget,
		pack.TTLabel(pack.NewTTLabelArgs(nil).SetText("神秘代码：设置")),
	)
}

func (customPage *Settings) Reset() {

}
