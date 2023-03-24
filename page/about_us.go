package page

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"sync"
	"translator/tst/tt_ui/pack"
	"translator/util"
)

var (
	apiAboutUs  *AboutUs
	onceAboutUs sync.Once
)

func GetAboutUs() *AboutUs {
	onceAboutUs.Do(func() {
		apiAboutUs = new(AboutUs)
		apiAboutUs.id = util.Uid()
		apiAboutUs.name = "关于我们"
	})
	return apiAboutUs
}

type AboutUs struct {
	id         string
	name       string
	mainWindow *walk.MainWindow
	rootWidget *walk.Composite
}

func (customPage *AboutUs) GetId() string {
	return customPage.id
}

func (customPage *AboutUs) GetName() string {
	return customPage.name
}

func (customPage *AboutUs) BindWindow(win *walk.MainWindow) {
	customPage.mainWindow = win
}

func (customPage *AboutUs) SetVisible(isVisible bool) {
	if customPage.rootWidget != nil {
		customPage.rootWidget.SetVisible(isVisible)
	}
}

func (customPage *AboutUs) GetWidget() Widget {
	return StdRootWidget(&customPage.rootWidget,
		pack.TTLabel(pack.NewTTLabelArgs(nil).SetText("神秘代码：关于我们")),
	)
}

func (customPage *AboutUs) Reset() {

}
