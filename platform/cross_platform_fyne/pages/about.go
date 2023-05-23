package pages

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"net/url"
	"sync"
)

var (
	apiPageAbout  *PageAbout
	oncePageAbout sync.Once
)

func APIPageAbout() *PageAbout {
	oncePageAbout.Do(func() {
		apiPageAbout = &PageAbout{
			id:   "page.about",
			name: "关于",
		}
	})
	return apiPageAbout
}

type PageAbout struct {
	window    fyne.Window
	id        string
	name      string
	isDefault bool
}

func (page *PageAbout) GetID() string { return page.id }

func (page *PageAbout) GetName() string { return page.name }

func (page *PageAbout) GetWindow() fyne.Window { return page.window }

func (page *PageAbout) SetWindow(win fyne.Window) { page.window = win }

func (page *PageAbout) IsDefault() bool {
	return page.isDefault
}

func (page *PageAbout) OnClose() {

}

func (page *PageAbout) OnReset() {

}

func (page *PageAbout) OnRender() fyne.CanvasObject {
	tmpLabel := widget.NewLabel(`您好，欢迎使用Anto，一款不专业的开源字幕翻译应用，主要支持当前主流翻译服务(主要包括阿里云、腾讯云、有道云和百度等)

和异步翻译，间接有效地提高了您的外语视频观看体验。如果有什么意见建议之类的，欢迎反馈~联系方式就在最下面哦~
`)
	wxImage := canvas.NewImageFromURI(storage.NewURI("https://speauty.oss-cn-chengdu.aliyuncs.com/anto/wxpay.jpg"))
	wxImage.SetMinSize(fyne.NewSize(360, 300))
	wxImage.FillMode = canvas.ImageFillContain

	aLiImage := canvas.NewImageFromURI(storage.NewURI("https://speauty.oss-cn-chengdu.aliyuncs.com/anto/alipay.jpg"))
	aLiImage.SetMinSize(fyne.NewSize(360, 300))
	aLiImage.FillMode = canvas.ImageFillContain

	repositoryUrl, _ := url.Parse("https://github.com/speauty/anto")

	return container.NewVBox(
		pageTitle(page.GetName()),
		tmpLabel,
		container.NewHBox(wxImage, widget.NewSeparator(), aLiImage),
		container.NewCenter(widget.NewLabel("~如果喜欢的话，那就给作者加个鸡腿吧~")),
		widget.NewHyperlink("项目地址: https://github.com/speauty/anto", repositoryUrl),
		widget.NewLabel("邮箱地址: speauty@163.com"),
	)
}
