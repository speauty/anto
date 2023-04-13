package page

import (
	"anto/lib/util"
	"anto/platform/win/ui/pack"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"sync"
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
		pack.UIGroupBox(
			pack.NewUIGroupBoxArgs(nil).
				SetTitle("关于我们").SetVisible(true).SetLayoutVBox(false).
				SetWidgets(
					pack.NewWidgetGroup().Append(
						pack.UITextLabel(pack.NewUITextLabelArgs(nil).
							SetCustomSize(Size{Width: 100, Height: 300}).
							SetText(
								`    你好，欢迎使用字幕翻译工具（以下简称：Anto），我是Anto的研发人员Speauty，现主要研究字幕方面，比如提取、翻译、合成等。所以如果有什么好的想法，可以直接给我发邮件。
    
交流邮箱：speauty@163.com
Github项目地址：https://github.com/speauty/tools.subrip.translator
`)),
					).AppendZeroVSpacer().GetWidgets(),
				),
		),
	)
}

func (customPage *AboutUs) Reset() {

}
