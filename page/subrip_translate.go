package page

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"sync"
	"translator/domain"
	"translator/tst/tt_translator/ling_va"
	"translator/tst/tt_ui/pack"
	_type "translator/type"
	"translator/util"
)

type LangItem struct {
	Key  string
	Name string
}

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
	})
	return apiSubripTranslate
}

type SubripTranslate struct {
	id          string
	name        string
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
	engines := domain.GetTranslators().GetNames()
	fmt.Println(engines)

	return StdRootWidget(&customPage.rootWidget,
		pack.TTComposite(pack.NewTTCompositeArgs(nil).
			SetLayoutHBox(true).SetWidgets(
			pack.NewWidgetGroup().Append(
				pack.TTLabel(pack.NewTTLabelArgs(nil).SetText("翻译引擎")),
				pack.TTComboBox(pack.NewTTComboBoxArgs(&customPage.ptrEngine).
					SetModel(engines).SetBindingMember(comboBoxModel.BindKey()).SetDisplayMember(comboBoxModel.DisplayKey()).SetCurrentIdx(0)),
				pack.TTLabel(pack.NewTTLabelArgs(nil).SetText("来源语种")),
				pack.TTComboBox(pack.NewTTComboBoxArgs(&customPage.ptrFromLang).
					SetModel(ling_va.GetInstance().GetLangSupported()).SetBindingMember(comboBoxModel.BindKey()).SetDisplayMember(comboBoxModel.DisplayKey()).SetCurrentIdx(0)),
				pack.TTLabel(pack.NewTTLabelArgs(nil).SetText("目标语种")),
				pack.TTComboBox(pack.NewTTComboBoxArgs(&customPage.ptrToLang).
					SetModel(ling_va.GetInstance().GetLangSupported()).SetBindingMember(comboBoxModel.BindKey()).SetDisplayMember(comboBoxModel.DisplayKey()).SetCurrentIdx(1)),
			).AppendZeroHSpacer().GetWidgets(),
		)),
		pack.TTComposite(pack.NewTTCompositeArgs(nil).
			SetLayoutHBox(true).SetWidgets(
			pack.NewWidgetGroup().Append(
				pack.TTLabel(pack.NewTTLabelArgs(nil).SetText("翻译模式")),
				pack.TTComboBox(pack.NewTTComboBoxArgs(nil).
					SetModel([]string{"全量翻译", "增量翻译"}).SetCurrentIdx(0)),
				pack.TTLabel(pack.NewTTLabelArgs(nil).SetText("导出主轨道")),
				pack.TTComboBox(pack.NewTTComboBoxArgs(nil).
					SetModel([]string{"来源语种", "目标语种"}).SetCurrentIdx(0)),
			).AppendZeroHSpacer().GetWidgets(),
		)),
		VSpacer{},
	)
}

func (customPage *SubripTranslate) Reset() {

}
