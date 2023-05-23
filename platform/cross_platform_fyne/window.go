package cross_platform_fyne

import (
	"context"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/golang-module/carbon"
	"sync"
)

var (
	instance *AppGui
	once     sync.Once
)

func API() *AppGui {
	once.Do(func() {
		instance = new(AppGui)
	})
	return instance
}

type AppGui struct {
	ctx            context.Context
	ctxCancelFn    context.CancelFunc
	config         *Config
	app            fyne.App
	mainWindow     fyne.Window
	consoleWindow  fyne.Window
	pages          map[string]ImplPage
	currentPageId  string
	chanNextPageId chan string
	chanConsole    chan string
	listConsole    []string
}

func (ag *AppGui) Config() *Config {
	return ag.config
}

func (ag *AppGui) App() fyne.App {
	return ag.app
}

func (ag *AppGui) MainWindow() fyne.Window {
	return ag.mainWindow
}

func (ag *AppGui) ConsoleWindow() fyne.Window {
	return ag.consoleWindow
}

func (ag *AppGui) Init(config *Config) {
	if ag.app != nil { // 禁止重复初始化
		return
	}
	if config == nil {
		config = config.Default()
	}
	ag.config = config

	ag.initApp()

	ag.initMainWindow()
	ag.initConsoleWindow()
	ag.initPages()

	ag.RegisterMainMenus(menus())

	t := &customTheme{}
	t.SetFonts("simkai.ttf", "")
	ag.app.Settings().SetTheme(t)
}

func (ag *AppGui) RegisterMainMenus(fyneMenu *fyne.MainMenu) {
	ag.mainWindow.SetMainMenu(fyneMenu)
}

func (ag *AppGui) Run(ctx context.Context, fnCancel context.CancelFunc) {
	ag.ctx = ctx
	ag.ctxCancelFn = fnCancel

	ag.mainWindow.SetCloseIntercept(ag.eventClose)

	defer func() {
		if err := recover(); err != nil {
			ag.PushToConsole(fmt.Sprintf("时间: %s, 致命错误: %s", carbon.Now(), err))
		}
	}()
	go func() {
		for {
			appSize := ag.mainWindow.Content().Size()
			if appSize.Width > appMainWindowDefaultWidth || appSize.Height > appMainWindowDefaultHeight {
				fmt.Println(fmt.Sprintf(
					"主窗口大小更新: (%.2f, %.2f), 标准大小: (%.2f, %.2f)",
					appSize.Width, appSize.Height,
					appMainWindowDefaultWidth, appMainWindowDefaultHeight,
				))
			}
			break
		}
	}()
	ag.app.Run()
}

func (ag *AppGui) Close() {
	ag.app.Quit()
}

func (ag *AppGui) eventClose() {
	ag.ctxCancelFn()
	ag.app.Quit()
}

func (ag *AppGui) initApp() {
	ag.app = app.New()
}

func (ag *AppGui) initMainWindow() {
	mainWindow := ag.app.NewWindow(ag.config.AppName)
	mainWindow.SetPadded(true)
	mainWindow.Resize(ag.config.MainWindowSize())
	mainWindow.CenterOnScreen()
	mainWindow.SetFixedSize(true)
	mainWindow.Show()
	mainWindow.SetMaster()
	ag.mainWindow = mainWindow
}
