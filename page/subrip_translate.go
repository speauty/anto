package page

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"sync"
	"translator/domain"
	"translator/tst/tt_translator/ling_va"
	"translator/tst/tt_ui/msg"
	"translator/tst/tt_ui/pack"
	_type "translator/type"
	"translator/util"
)

var (
	apiSubripTranslate  *SubripTranslate
	onceSubripTranslate sync.Once
	comboBoxModel       = new(_type.StdComboBoxModel)
)

func GetSubripTranslate() *SubripTranslate {
	onceSubripTranslate.Do(func() {
		apiSubripTranslate = new(SubripTranslate)
		apiSubripTranslate.id = util.Uid()
		apiSubripTranslate.name = "字幕翻译"
		apiSubripTranslate.engines = domain.GetTranslators().GetNames()

	})
	return apiSubripTranslate
}

type SubripTranslate struct {
	id          string
	name        string
	engines     []*_type.StdComboBoxModel
	mainWindow  *walk.MainWindow
	rootWidget  *walk.Composite
	ptrEngine   *walk.ComboBox
	ptrFromLang *walk.ComboBox
	ptrToLang   *walk.ComboBox
}

func (customPage *SubripTranslate) GetId() string {
	return customPage.id
}

func (customPage *SubripTranslate) GetName() string {
	return customPage.name
}

func (customPage *SubripTranslate) BindWindow(win *walk.MainWindow) {
	customPage.mainWindow = win
}

func (customPage *SubripTranslate) SetVisible(isVisible bool) {
	if customPage.rootWidget != nil {
		customPage.rootWidget.SetVisible(isVisible)
	}
}

func (customPage *SubripTranslate) GetWidget() Widget {
	return StdRootWidget(&customPage.rootWidget,
		pack.TTComposite(pack.NewTTCompositeArgs(nil).
			SetLayoutHBox(true).SetWidgets(
			pack.NewWidgetGroup().Append(
				pack.TTLabel(pack.NewTTLabelArgs(nil).SetText("翻译引擎")),
				pack.TTComboBox(pack.NewTTComboBoxArgs(&customPage.ptrEngine).
					SetModel(apiSubripTranslate.engines).SetBindingMember(comboBoxModel.BindKey()).SetDisplayMember(comboBoxModel.DisplayKey()).SetCurrentIdx(0).SetOnCurrentIdxChangedFn(customPage.eventEngineOnChange)),
				pack.TTLabel(pack.NewTTLabelArgs(nil).SetText(_type.LangDirectionFrom.String())),
				pack.TTComboBox(pack.NewTTComboBoxArgs(&customPage.ptrFromLang).
					SetModel(ling_va.GetInstance().GetLangSupported()).SetBindingMember(comboBoxModel.BindKey()).SetDisplayMember(comboBoxModel.DisplayKey())),
				pack.TTLabel(pack.NewTTLabelArgs(nil).SetText(_type.LangDirectionTo.String())),
				pack.TTComboBox(pack.NewTTComboBoxArgs(&customPage.ptrToLang).
					SetModel(ling_va.GetInstance().GetLangSupported()).SetBindingMember(comboBoxModel.BindKey()).SetDisplayMember(comboBoxModel.DisplayKey()).SetCurrentIdx(1)),
			).AppendZeroHSpacer().GetWidgets(),
		)),
		pack.TTComposite(pack.NewTTCompositeArgs(nil).
			SetLayoutHBox(true).SetWidgets(
			pack.NewWidgetGroup().Append(
				pack.TTLabel(pack.NewTTLabelArgs(nil).SetText("翻译模式")),
				pack.TTComboBox(pack.NewTTComboBoxArgs(nil).
					SetModel(_type.ModeFull.GetModes()).SetCurrentIdx(_type.ModeFull.GetIdx())),
				pack.TTLabel(pack.NewTTLabelArgs(nil).SetText("导出主轨道")),
				pack.TTComboBox(pack.NewTTComboBoxArgs(nil).
					SetModel(_type.LangDirectionFrom.GetDirections()).SetCurrentIdx(0)),
			).AppendZeroHSpacer().GetWidgets(),
		)),
		VSpacer{},
	)
}

func (customPage *SubripTranslate) Reset() {

}

func (customPage *SubripTranslate) eventEngineOnChange() {
	currentId := ""
	for _, engine := range customPage.engines {
		if engine.Name == customPage.ptrEngine.Text() {
			currentId = engine.Key
			break
		}
	}
	if currentId == "" {
		msg.Err(customPage.mainWindow, fmt.Errorf("当前翻译引擎无效，请重新选择"))
		_ = customPage.ptrEngine.SetCurrentIndex(-1)
		customPage.setLangComboBox(customPage.ptrFromLang, nil, -1)
		customPage.setLangComboBox(customPage.ptrToLang, nil, -1)
		return
	}
	currentEngine := domain.GetTranslators().GetById(currentId)
	if customPage.mainWindow != nil && currentEngine == nil {
		msg.Err(customPage.mainWindow, fmt.Errorf("当前翻译引擎未注册，请重新选择"))
		_ = customPage.ptrEngine.SetCurrentIndex(-1)
		customPage.setLangComboBox(customPage.ptrFromLang, nil, -1)
		customPage.setLangComboBox(customPage.ptrToLang, nil, -1)
		return
	}
	customPage.setLangComboBox(customPage.ptrFromLang, currentEngine.GetLangSupported(), 1)
	customPage.setLangComboBox(customPage.ptrToLang, currentEngine.GetLangSupported(), 0)
	return
}

func (customPage *SubripTranslate) setLangComboBox(ptr *walk.ComboBox, model interface{}, idx int) {
	_ = ptr.SetModel(model)
	_ = ptr.SetCurrentIndex(idx)
}
