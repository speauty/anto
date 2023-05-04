package ui

import (
	"anto/lib/log"
	"context"
	"sync"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

const width = 800
const height = 600

var (
	apiSingleton  *Gui
	onceSingleton sync.Once
)

func Singleton() *Gui {
	onceSingleton.Do(func() {
		apiSingleton = new(Gui)
		apiSingleton.pageCtl = new(PageCtl)
	})
	return apiSingleton
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

	if err := customG.initMainWindow(); err != nil {
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

func (customG *Gui) initMainWindow() error {
	return MainWindow{
		AssignTo: &customG.win,
		Title:    customG.cfg.Title, Icon: "./favicon.ico",
		Visible: false, Persistent: true,
		Layout:    VBox{MarginsZero: true},
		MenuItems: customG.menus,
		Children:  customG.pageCtl.GetWidgets(),
	}.Create()
}

func (customG *Gui) log() *log.Log {
	return log.Singleton()
}
