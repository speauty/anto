package cross_platform_fyne

import (
	"anto/platform/cross_platform_fyne/pages"
	"fyne.io/fyne/v2"
)

func menus() *fyne.MainMenu {
	appGui := API()

	{
		pages.APIPageEnv().SetWindow(appGui.MainWindow())
		pages.APIPageConfig().SetWindow(appGui.MainWindow())
		pages.APIPageAbout().SetWindow(appGui.MainWindow())
		pages.APIPageSubtitleTranslate().SetWindow(appGui.MainWindow())

		appGui.RegisterPages(
			pages.APIPageEnv(), pages.APIPageConfig(), pages.APIPageAbout(),
			pages.APIPageSubtitleTranslate(),
		)
	}

	return fyne.NewMainMenu(
		&fyne.Menu{
			Label: "文件", Items: []*fyne.MenuItem{
				{Label: "环境", Action: triggerMenuEnv},
				{Label: "设置", Action: triggerMenuConfig},
				{Label: "控制台", Action: triggerMenuConsole},
			},
		},
		&fyne.Menu{
			Label: "翻译助手", Items: []*fyne.MenuItem{
				{Label: "字幕翻译", Action: triggerMenuSubtitleTranslate},
			},
		},
		&fyne.Menu{
			Label: "帮助", Items: []*fyne.MenuItem{
				{Label: "清空控制台", Action: appGui.ClearConsole}, // 临时菜单, 后期换成右键
				{Label: "关于", Action: triggerMenuAbout},
			},
		},
	)
}

func triggerMenuEnv() {
	API().toPage(pages.APIPageEnv().GetID())
}

func triggerMenuConfig() {
	API().toPage(pages.APIPageConfig().GetID())
}

func triggerMenuConsole() {
	API().ConsoleWindow().Show()
}

func triggerMenuSubtitleTranslate() {
	API().toPage(pages.APIPageSubtitleTranslate().GetID())
}

func triggerMenuAbout() {
	API().toPage(pages.APIPageAbout().GetID())
}
