package ui

import (
	"anto/lib/log"
	"context"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"sync"
)

const width = 800
const height = 600

var (
	apiGui  *Gui
	onceGui sync.Once
)

func GetInstance() *Gui {
	onceGui.Do(func() {
		apiGui = new(Gui)
		apiGui.pageCtl = new(PageCtl)
	})
	return apiGui
}

type Gui struct {
	ctx         context.Context
	ctxCancelFn context.CancelFunc
	win         *walk.MainWindow
	cfg         *Cfg
	menus       []MenuItem
	pageCtl     *PageCtl
}

func (customG *Gui) Init(cfg *Cfg) error {
	customG.cfg = cfg

	if err := customG.genMainWindow(); err != nil {
		return err
	}

	customG.setWindowFlag()
	customG.setWindowCenter()
	customG.setWindowMinAndMaxSize()

	customG.pageCtl.Bind(customG.GetWindow())

	customG.GetWindow().SetVisible(true)

	_ = customG.GetWindow().SetFocus()

	return nil
}

func (customG *Gui) Run(ctx context.Context, fnCancel context.CancelFunc) {
	customG.ctx = ctx
	customG.ctxCancelFn = fnCancel

	customG.win.Closing().Attach(customG.eventCustomClose)

	customG.GetWindow().Run()
}

func (customG *Gui) Close() {
	_ = customG.win.Close()
}

func (customG *Gui) eventCustomClose(canceled *bool, reason walk.CloseReason) {
	reason = walk.CloseReasonUser
	*canceled = false
	customG.ctxCancelFn()
}

func (customG *Gui) GetWindow() *walk.MainWindow {
	return customG.win
}

func (customG *Gui) RegisterPages(pages ...IPage) {
	customG.pageCtl.PushPages(pages...)
	return
}

func (customG *Gui) GoPage(pageId string) error {
	return customG.pageCtl.SetCurrent(pageId)
}

func (customG *Gui) RegisterMenus(menus []MenuItem) {
	customG.menus = menus
	return
}

func (customG *Gui) genMainWindow() error {
	return MainWindow{
		AssignTo:   &customG.win,
		Title:      customG.cfg.Title,
		Icon:       "./favicon.ico",
		Visible:    false,
		Persistent: true,
		Layout:     VBox{MarginsZero: true},
		MenuItems:  customG.menus,
		Children:   customG.pageCtl.GetWidgets(),
	}.Create()
}

func (customG *Gui) setWindowFlag() {
	win.SetWindowLong(customG.GetWindow().Handle(), win.GWL_STYLE,
		win.GetWindowLong(customG.GetWindow().Handle(), win.GWL_STYLE) & ^win.WS_MAXIMIZEBOX & ^win.WS_THICKFRAME)
}

func (customG *Gui) setWindowCenter() {
	scrWidth := win.GetSystemMetrics(win.SM_CXSCREEN)
	scrHeight := win.GetSystemMetrics(win.SM_CYSCREEN)
	_ = customG.GetWindow().SetBounds(walk.Rectangle{
		X:      int((scrWidth - width) / 2),
		Y:      int((scrHeight - height) / 2),
		Width:  width,
		Height: height,
	})
}

func (customG *Gui) setWindowMinAndMaxSize() {
	minMaxSize := walk.Size{Width: width, Height: height}
	_ = customG.GetWindow().SetMinMaxSize(minMaxSize, minMaxSize)
}

func (customG *Gui) log() *log.Log {
	return log.Singleton()
}
