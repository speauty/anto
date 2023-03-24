package tt_ui

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"sync"
	"translator/tst/tt_log"
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
	win     *walk.MainWindow
	cfg     *Cfg
	menus   []MenuItem
	pageCtl *PageCtl
}

func (customG *Gui) Init(cfg *Cfg) error {
	customG.cfg = cfg

	if err := customG.genMainWindow(); err != nil {
		return err
	}

	customG.setWindowFlag()
	customG.setWindowCenter()

	customG.pageCtl.Bind(customG.GetWindow())

	customG.GetWindow().SetVisible(true)

	_ = customG.GetWindow().SetFocus()

	return nil
}

func (customG *Gui) Run() {
	customG.GetWindow().Run()
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
		AssignTo:       &customG.win,
		Title:          customG.cfg.Title,
		Icon:           customG.cfg.Icon,
		Visible:        false,
		Layout:         VBox{MarginsZero: true},
		MenuItems:      customG.menus,
		StatusBarItems: customG.defaultStatusBars(),
		Children:       customG.pageCtl.GetWidgets(),
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

func (customG *Gui) log() *tt_log.TTLog {
	return tt_log.GetInstance()
}
