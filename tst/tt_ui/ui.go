package tt_ui

import (
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
	})
	return apiGui
}

type Gui struct {
	win *walk.MainWindow
	cfg *Cfg
}

func (customG *Gui) Init(cfg *Cfg) error {
	customG.cfg = cfg

	if err := customG.genMainWindow(); err != nil {
		return err
	}

	customG.setWindowFlag()
	customG.setWindowCenter()
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

func (customG *Gui) genMainWindow() error {
	return MainWindow{
		AssignTo:       &customG.win,
		Title:          customG.cfg.Title,
		Icon:           customG.cfg.Icon,
		Visible:        false,
		Layout:         VBox{MarginsZero: true},
		MenuItems:      customG.defaultMenu(),
		StatusBarItems: customG.defaultStatusBars(),
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
