package page

import (
	"anto/lib/util"
	"anto/platform/win/ui/pack"
	"anto/resource"
	"bytes"
	"image/jpeg"
	"sync"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var (
	apiAboutUs  *AboutUs
	onceAboutUs sync.Once
	wxPayImg    *walk.Bitmap
	aliPayImg   *walk.Bitmap
)

func GetAboutUs() *AboutUs {
	onceAboutUs.Do(func() {
		apiAboutUs = new(AboutUs)
		apiAboutUs.id = util.Uid()
		apiAboutUs.name = "关于我们"

		wxPayData, err := jpeg.Decode(bytes.NewReader(resource.WxPay))
		if err == nil && wxPayData != nil {
			wxPayImg, _ = walk.NewBitmapFromImage(wxPayData)
		}

		aliPayData, err := jpeg.Decode(bytes.NewReader(resource.ALiPay))
		if err == nil && aliPayData != nil {
			aliPayImg, _ = walk.NewBitmapFromImage(aliPayData)
		}

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
		customPage.getAboutUsWidget(),
		customPage.getSponsorshipWidget(),
	)
}

func (customPage *AboutUs) getAboutUsWidget() Widget {
	return pack.UIGroupBox(
		pack.NewUIGroupBoxArgs(nil).
			SetTitle("关于我们").SetCustomSize(Size{Height: 90}).SetLayoutVBox(false).
			SetWidgets(
				pack.NewWidgetGroup().Append(
					pack.UITextLabel(pack.NewUITextLabelArgs(nil).
						SetCustomSize(Size{Width: 100}).
						SetText(
							`    你好，欢迎使用字幕翻译工具（以下简称：Anto），我是Anto的研发人员Speauty，现主要研究字幕方面，比如提取、翻译、合成等。所以如果有什么好的想法，可以直接给我发邮件。

交流邮箱：speauty@163.com
项目地址：https://github.com/speauty/anto
`)),
				).AppendZeroVSpacer().GetWidgets(),
			),
	)
}

func (customPage *AboutUs) getSponsorshipWidget() Widget {
	return pack.UIGroupBox(
		pack.NewUIGroupBoxArgs(nil).
			SetTitle("加个鸡腿").SetLayoutHBox(false).
			SetWidgets(
				pack.NewWidgetGroup().Append(
					pack.UIImageView(pack.NewUIImageViewArgs(nil).SetImage(wxPayImg)),
					pack.UIImageView(pack.NewUIImageViewArgs(nil).SetImage(aliPayImg)),
				).AppendZeroHSpacer().AppendZeroVSpacer().GetWidgets(),
			),
	)
}

func (customPage *AboutUs) Reset() {}
