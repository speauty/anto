package page

import (
	"anto/cfg"
	"anto/common"
	"anto/domain/repository"
	"anto/domain/service/translator"
	"anto/lib/log"
	"anto/lib/util"
	"anto/platform/win/ui/msg"
	"anto/platform/win/ui/pack"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"go.uber.org/zap"
)

var (
	apiSettings        *Settings
	onceSettings       sync.Once
	stdLineEditSize    = Size{Width: 100}
	stdNumLineEditSize = Size{Width: 40}
)

func GetSettings() *Settings {
	onceSettings.Do(func() {
		apiSettings = new(Settings)
		apiSettings.id = util.Uid()
		apiSettings.name = "设置"
		apiSettings.engines = repository.GetTranslators().GetNamesAll()
	})
	return apiSettings
}

type Settings struct {
	id         string
	name       string
	mainWindow *walk.MainWindow
	rootWidget *walk.Composite

	engines               []*common.StdComboBoxModel
	currentEngineId       string
	ptrCurrentEngineCombo *walk.ComboBox

	ptrAKWrapper *walk.Composite
	ptrAKInput   *walk.LineEdit

	ptrSKWrapper *walk.Composite
	ptrSKInput   *walk.LineEdit

	ptrPKWrapper *walk.Composite
	ptrPKInput   *walk.LineEdit

	ptrRegionWrapper *walk.Composite
	ptrRegionInput   *walk.LineEdit

	ptrMaxCharNumWrapper *walk.Composite
	ptrMaxCharNumInput   *walk.LineEdit

	ptrQPSWrapper *walk.Composite
	ptrQPSInput   *walk.LineEdit

	ptrMaxCoroutineNumWrapper *walk.Composite
	ptrMaxCoroutineNumInput   *walk.LineEdit
}

func (customPage *Settings) GetId() string { return customPage.id }

func (customPage *Settings) GetName() string { return customPage.name }

func (customPage *Settings) BindWindow(win *walk.MainWindow) { customPage.mainWindow = win }

func (customPage *Settings) SetVisible(isVisible bool) {
	if customPage.rootWidget != nil {
		customPage.rootWidget.SetVisible(isVisible)
	}
}

func (customPage *Settings) GetWidget() Widget {
	defer customPage.Reset()
	return StdRootWidget(&customPage.rootWidget,
		pack.UIScrollView(pack.NewUIScrollViewArgs(nil).SetChildren(
			pack.NewWidgetGroup().Append(
				customPage.getTranslationConfigWidget(),
			).AppendZeroVSpacer().GetWidgets())),
	)
}

func (customPage *Settings) Reset() {}

func (customPage *Settings) getTranslationConfigWidget() Widget {
	return pack.UIGroupBox(pack.NewUIGroupBoxArgs(nil).SetTitle("翻译引擎").SetLayoutVBox(false).SetCustomSize(Size{Width: 180}).SetWidgets(
		pack.NewWidgetGroup().Append(
			pack.UILabel(pack.NewUILabelArgs(nil).SetText("~不熟悉的参数保持默认即可, 协程数量必须小于QPS哦~")),
			pack.UIComposite(pack.NewUICompositeArgs(nil).SetLayoutHBox(true).SetWidgets(
				pack.NewWidgetGroup().Append(
					pack.UILabel(pack.NewUILabelArgs(nil).SetText("设置").SetCustomSize(Size{Width: 50})),
					pack.UIComboBox(pack.NewUIComboBoxArgs(&customPage.ptrCurrentEngineCombo).SetCustomSize(stdLineEditSize).SetModel(apiSettings.engines).
						SetBindingMember(comboBoxModel.BindKey()).SetDisplayMember(comboBoxModel.DisplayKey()).
						SetCurrentIdx(-1).
						SetOnCurrentIdxChangedFn(customPage.eventCurrentEngineComboOnChanged),
					),
				).AppendZeroHSpacer().GetWidgets(),
			)),
			pack.UIComposite(pack.NewUICompositeArgs(nil).SetLayoutHBox(true).SetWidgets(
				pack.NewWidgetGroup().Append(
					customPage.builderFormItemInputWidget(&customPage.ptrAKWrapper, &customPage.ptrAKInput, "服务应用", stdLineEditSize),
					customPage.builderFormItemInputWidget(&customPage.ptrSKWrapper, &customPage.ptrSKInput, "服务密钥", stdLineEditSize),
					customPage.builderFormItemInputWidget(&customPage.ptrPKWrapper, &customPage.ptrPKInput, "服务项目", stdLineEditSize),
				).AppendZeroHSpacer().GetWidgets(),
			)),
			pack.UIComposite(pack.NewUICompositeArgs(nil).SetLayoutHBox(true).SetWidgets(
				pack.NewWidgetGroup().Append(
					customPage.builderFormItemInputWidget(&customPage.ptrRegionWrapper, &customPage.ptrRegionInput, "服务地区", stdLineEditSize),
				).AppendZeroHSpacer().GetWidgets(),
			)),
			pack.UIComposite(pack.NewUICompositeArgs(nil).SetLayoutHBox(true).SetWidgets(
				pack.NewWidgetGroup().Append(
					customPage.builderFormItemInputWidget(&customPage.ptrMaxCharNumWrapper, &customPage.ptrMaxCharNumInput, "最大长度", stdLineEditSize),
					customPage.builderFormItemInputWidget(&customPage.ptrQPSWrapper, &customPage.ptrQPSInput, "QPS", stdLineEditSize),
					customPage.builderFormItemInputWidget(&customPage.ptrMaxCoroutineNumWrapper, &customPage.ptrMaxCoroutineNumInput, "协程数量", stdLineEditSize),
				).AppendZeroHSpacer().GetWidgets(),
			)),
			pack.UIComposite(pack.NewUICompositeArgs(nil).SetLayoutHBox(true).SetWidgets(
				pack.NewWidgetGroup().Append(
					pack.UIPushBtn(pack.NewUIPushBtnArgs(nil).SetText("保存").SetOnClicked(customPage.eventSubmit)),
					pack.UIPushBtn(pack.NewUIPushBtnArgs(nil).SetText("一键还原").SetOnClicked(customPage.eventRestore)),
				).AppendZeroHSpacer().GetWidgets(),
			)),
		).AppendZeroHSpacer().GetWidgets(),
	))
}

func (customPage *Settings) builderFormItemInputWidget(ptrWrapper **walk.Composite, ptrInput **walk.LineEdit, label string, inputSize Size) Widget {
	return pack.UIComposite(pack.NewUICompositeArgs(ptrWrapper).SetLayoutHBox(true).SetCustomSize(Size{Width: 180}).SetAlignment(AlignHCenterVFar).SetWidgets(
		pack.NewWidgetGroup().Append(
			pack.UILabel(pack.NewUILabelArgs(nil).SetCustomSize(Size{Width: 50}).SetText(label)),
			pack.UILineEdit(pack.NewUILineEditArgs(ptrInput).SetCustomSize(inputSize)),
		).AppendZeroHSpacer().GetWidgets(),
	))
}

func (customPage *Settings) eventCurrentEngineComboOnChanged() {
	currentId := customPage.engines[customPage.ptrCurrentEngineCombo.CurrentIndex()].Key
	if currentId == "" {
		msg.Err(customPage.mainWindow, fmt.Errorf("当前翻译引擎无效，请重新选择"))
		_ = customPage.ptrCurrentEngineCombo.SetCurrentIndex(-1)
		return
	}
	if customPage.currentEngineId == currentId {
		return
	}
	currentEngine := repository.GetTranslators().GetById(currentId)
	if customPage.mainWindow != nil && currentEngine == nil {
		msg.Err(customPage.mainWindow, fmt.Errorf("当前翻译引擎未注册，请重新选择"))
		_ = customPage.ptrCurrentEngineCombo.SetCurrentIndex(-1)
		return
	}

	customPage.resetFormInputs()

	if currentEngine.GetCfg().GetAK() != translator.ConfigInvalidStr {
		customPage.ptrAKWrapper.SetVisible(true)
		_ = customPage.ptrAKInput.SetText(currentEngine.GetCfg().GetAK())
	}

	if currentEngine.GetCfg().GetSK() != translator.ConfigInvalidStr {
		customPage.ptrSKWrapper.SetVisible(true)
		_ = customPage.ptrSKInput.SetText(currentEngine.GetCfg().GetSK())
	}

	if currentEngine.GetCfg().GetProjectKey() != translator.ConfigInvalidStr {
		customPage.ptrPKWrapper.SetVisible(true)
		_ = customPage.ptrPKInput.SetText(currentEngine.GetCfg().GetProjectKey())
	}

	if currentEngine.GetCfg().GetRegion() != translator.ConfigInvalidStr {
		customPage.ptrRegionWrapper.SetVisible(true)
		_ = customPage.ptrRegionInput.SetText(currentEngine.GetCfg().GetRegion())
	}

	if currentEngine.GetCfg().GetMaxCharNum() != translator.ConfigInvalidInt {
		customPage.ptrMaxCharNumWrapper.SetVisible(true)
		_ = customPage.ptrMaxCharNumInput.SetText(strconv.Itoa(currentEngine.GetCfg().GetMaxCharNum()))
	}

	if currentEngine.GetCfg().GetQPS() != translator.ConfigInvalidInt {
		customPage.ptrQPSWrapper.SetVisible(true)
		_ = customPage.ptrQPSInput.SetText(strconv.Itoa(currentEngine.GetCfg().GetQPS()))
	}

	if currentEngine.GetCfg().GetMaxCoroutineNum() != translator.ConfigInvalidInt {
		customPage.ptrMaxCoroutineNumWrapper.SetVisible(true)
		_ = customPage.ptrMaxCoroutineNumInput.SetText(strconv.Itoa(currentEngine.GetCfg().GetMaxCoroutineNum()))
	}

	customPage.currentEngineId = currentId
}

func (customPage *Settings) eventSubmit() {

	currentId := customPage.engines[customPage.ptrCurrentEngineCombo.CurrentIndex()].Key
	if currentId == "" {
		msg.Err(customPage.mainWindow, fmt.Errorf("当前翻译引擎无效，请重新选择"))
		_ = customPage.ptrCurrentEngineCombo.SetCurrentIndex(-1)
		return
	}
	currentEngine := repository.GetTranslators().GetById(currentId)
	if customPage.mainWindow != nil && currentEngine == nil {
		msg.Err(customPage.mainWindow, fmt.Errorf("当前翻译引擎未注册，请重新选择"))
		_ = customPage.ptrCurrentEngineCombo.SetCurrentIndex(-1)
		return
	}

	if customPage.ptrAKWrapper.Visible() {
		if err := currentEngine.GetCfg().SetAK(customPage.ptrAKInput.Text()); err != nil {
			_ = customPage.ptrAKInput.SetFocus()
			msg.Err(customPage.mainWindow, err)
			return
		}
	}
	if customPage.ptrSKWrapper.Visible() {
		if err := currentEngine.GetCfg().SetSK(customPage.ptrSKInput.Text()); err != nil {
			_ = customPage.ptrSKInput.SetFocus()
			msg.Err(customPage.mainWindow, err)
			return
		}
	}
	if customPage.ptrPKWrapper.Visible() {
		if err := currentEngine.GetCfg().SetProjectKey(customPage.ptrPKInput.Text()); err != nil {
			_ = customPage.ptrPKInput.SetFocus()
			msg.Err(customPage.mainWindow, err)
			return
		}
	}
	if customPage.ptrRegionWrapper.Visible() {
		if err := currentEngine.GetCfg().SetRegion(customPage.ptrRegionInput.Text()); err != nil {
			_ = customPage.ptrPKInput.SetFocus()
			msg.Err(customPage.mainWindow, err)
			return
		}
	}
	if customPage.ptrMaxCharNumWrapper.Visible() {
		maxCharNumVal, err := strconv.Atoi(customPage.ptrMaxCharNumInput.Text())
		if err != nil {
			_ = customPage.ptrMaxCharNumInput.SetFocus()
			msg.Err(customPage.mainWindow, errors.New("无效数字"))
			return
		}
		if err := currentEngine.GetCfg().SetMaxCharNum(maxCharNumVal); err != nil {
			_ = customPage.ptrMaxCharNumInput.SetFocus()
			msg.Err(customPage.mainWindow, err)
			return
		}
	}
	if customPage.ptrQPSWrapper.Visible() {
		qpsVal, err := strconv.Atoi(customPage.ptrQPSInput.Text())
		if err != nil {
			_ = customPage.ptrQPSInput.SetFocus()
			msg.Err(customPage.mainWindow, errors.New("无效数字"))
			return
		}
		if err := currentEngine.GetCfg().SetQPS(qpsVal); err != nil {
			_ = customPage.ptrQPSInput.SetFocus()
			msg.Err(customPage.mainWindow, err)
			return
		}
	}
	if customPage.ptrMaxCoroutineNumWrapper.Visible() {
		maxCoroutineNumVal, err := strconv.Atoi(customPage.ptrMaxCoroutineNumInput.Text())
		if err != nil {
			_ = customPage.ptrMaxCoroutineNumInput.SetFocus()
			msg.Err(customPage.mainWindow, errors.New("无效数字"))
			return
		}
		if err := currentEngine.GetCfg().SetMaxCoroutineNum(maxCoroutineNumVal); err != nil {
			_ = customPage.ptrMaxCoroutineNumInput.SetFocus()
			msg.Err(customPage.mainWindow, err)
			return
		}
	}

	if err := currentEngine.GetCfg().SyncDisk(cfg.Singleton().GetViper()); err != nil {
		log.Singleton().Error("保存配置到文件失败", zap.Error(err))
		msg.Err(customPage.mainWindow, errors.New("保存配置失败"))
		return
	}

	if err := cfg.Singleton().Sync(); err != nil {
		log.Singleton().Error("写入配置失败", zap.Error(err))
		msg.Err(customPage.mainWindow, errors.New("写入配置失败"))
		return
	}

	msg.Info(customPage.mainWindow, "保存配置成功, 建议重启一下当前应用哦~如果没有生效的话")
}

func (customPage *Settings) eventRestore() {
	msg.Info(customPage.mainWindow, "手动删除当前目录下的[cfg.yml]文件, 重启应用即可, 配置会还原到初始默认, 请谨慎操作!!!")
}

func (customPage *Settings) resetFormInputs() {
	customPage.ptrAKWrapper.SetVisible(false)
	_ = customPage.ptrAKInput.SetText("")

	customPage.ptrSKWrapper.SetVisible(false)
	_ = customPage.ptrSKInput.SetText("")

	customPage.ptrPKWrapper.SetVisible(false)
	_ = customPage.ptrPKInput.SetText("")

	customPage.ptrRegionWrapper.SetVisible(false)
	_ = customPage.ptrRegionInput.SetText("")

	customPage.ptrQPSWrapper.SetVisible(false)
	_ = customPage.ptrQPSInput.SetText("")

	customPage.ptrMaxCharNumWrapper.SetVisible(false)
	_ = customPage.ptrMaxCharNumInput.SetText("")

	customPage.ptrMaxCoroutineNumWrapper.SetVisible(false)
	_ = customPage.ptrMaxCoroutineNumInput.SetText("")
}
