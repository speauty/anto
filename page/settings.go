package page

import (
	"errors"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"go.uber.org/zap"
	"strings"
	"sync"
	"translator/cfg"
	_const "translator/const"
	"translator/tst/tt_log"
	"translator/tst/tt_ui/msg"
	"translator/tst/tt_ui/pack"
	"translator/util"
)

var (
	apiSettings  *Settings
	onceSettings sync.Once
)

func GetSettings() *Settings {
	onceSettings.Do(func() {
		apiSettings = new(Settings)
		apiSettings.id = util.Uid()
		apiSettings.name = "设置"
	})
	return apiSettings
}

type Settings struct {
	id         string
	name       string
	mainWindow *walk.MainWindow
	rootWidget *walk.Composite

	ptrEnv *walk.ComboBox
}

func (customPage *Settings) GetId() string {
	return customPage.id
}

func (customPage *Settings) GetName() string {
	return customPage.name
}

func (customPage *Settings) BindWindow(win *walk.MainWindow) {
	customPage.mainWindow = win
}

func (customPage *Settings) SetVisible(isVisible bool) {
	if customPage.rootWidget != nil {
		customPage.rootWidget.SetVisible(isVisible)
	}
}

func (customPage *Settings) GetWidget() Widget {
	defer customPage.Reset()
	return StdRootWidget(&customPage.rootWidget,
		pack.TTComposite(pack.NewTTCompositeArgs(nil).SetLayoutVBox(true).SetWidgets(
			pack.NewWidgetGroup().Append(
				pack.TTGroupBox(pack.NewTTGroupBoxArgs(nil).SetTitle("应用").SetLayoutVBox(false).SetWidgets(
					pack.NewWidgetGroup().Append(
						pack.TTComposite(pack.NewTTCompositeArgs(nil).SetLayoutHBox(true).SetWidgets(
							pack.NewWidgetGroup().Append(
								pack.TTLabel(pack.NewTTLabelArgs(nil).SetText("环境")),
								pack.TTComboBox(pack.NewTTComboBoxArgs(&customPage.ptrEnv).SetModel([]string{"Debug", "Release"}).SetCurrentIdx(0)),
							).AppendZeroHSpacer().GetWidgets(),
						)),
						pack.TTComposite(pack.NewTTCompositeArgs(nil).SetLayoutHBox(true).SetWidgets(
							pack.NewWidgetGroup().Append(
								pack.TTPushBtn(pack.NewTTPushBtnArgs(nil).SetText("同步").SetOnClicked(customPage.eventSync)),
							).AppendZeroHSpacer().GetWidgets(),
						)),
					).AppendZeroHSpacer().GetWidgets(),
				)),
			).AppendZeroVSpacer().GetWidgets(),
		)),
	)
}

func (customPage *Settings) Reset() {
	if customPage.ptrEnv != nil {
		currentIdx := 0
		if cfg.GetInstance().App.Env == _const.EnvRelease {
			currentIdx = 1
		}
		_ = customPage.ptrEnv.SetCurrentIndex(currentIdx)
	}
}

func (customPage *Settings) eventSync() {
	currentEnv := strings.ToLower(customPage.ptrEnv.Text())
	cntModified := 0
	if currentEnv != cfg.GetInstance().App.Env {
		cfg.GetInstance().App.Env = currentEnv
		cntModified++
	}
	if cntModified == 0 {
		return
	}
	if err := cfg.GetInstance().Sync(); err != nil {
		tt_log.GetInstance().Error("同步配置到文件失败", zap.Error(err))
		msg.Err(customPage.mainWindow, errors.New("同步配置到文件失败"))
		return
	}
	msg.Info(customPage.mainWindow, "同步配置成功")
}
