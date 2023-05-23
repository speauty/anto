package pages

import (
	"anto/common"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"os"
	"runtime"
	"sync"
)

var (
	apiPageEnv  *PageEnv
	oncePageEnv sync.Once
)

func APIPageEnv() *PageEnv {
	oncePageEnv.Do(func() {
		apiPageEnv = &PageEnv{
			id:        "page.env",
			name:      "环境",
			isDefault: true,
		}
	})
	return apiPageEnv
}

type PageEnv struct {
	window    fyne.Window
	id        string
	name      string
	isDefault bool
}

func (page *PageEnv) GetID() string { return page.id }

func (page *PageEnv) GetName() string { return page.name }

func (page *PageEnv) GetWindow() fyne.Window { return page.window }

func (page *PageEnv) SetWindow(win fyne.Window) { page.window = win }

func (page *PageEnv) IsDefault() bool { return page.isDefault }

func (page *PageEnv) OnClose() {}

func (page *PageEnv) OnReset() {}

func (page *PageEnv) OnRender() fyne.CanvasObject {
	hostname, _ := os.Hostname()

	return container.NewVBox(
		pageTitle(page.GetName()),
		widget.NewCard("", "", container.NewVBox(
			widget.NewLabel(fmt.Sprintf(
				"系统  操作系统: %s-%s 主机名称: %s CPU核数: %d",
				runtime.GOOS, runtime.GOARCH,
				hostname, runtime.GOMAXPROCS(0),
			)),
			widget.NewSeparator(),
			widget.NewLabel(fmt.Sprintf(
				"应用  名称: %s 版本: %s",
				common.AppName, common.Version,
			)),
			widget.NewSeparator(),
			//widget.NewLabel("统计(待实现)"),
			//widget.NewLabel("字幕翻译: 总计(0) 成功(0) 成功率(0) 总耗时(0) 平均耗时(0) 最长耗时(0) 最短耗时(0)"),
			//widget.NewSeparator(),
		)),
	)
}
