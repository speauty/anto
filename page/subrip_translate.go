package page

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"sync"
	"translator/tst/tt_translator/ling_va"
	"translator/tst/tt_ui/pack"
	"translator/util"
)

type LangItem struct {
	Key  string
	Name string
}

var (
	apiSubripTranslate  *SubripTranslate
	onceSubripTranslate sync.Once
	languages           = []LangItem{
		{"en", "英语"},
		{"zh_cn", "中文"},
	}
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
	id         string
	name       string
	mainWindow *walk.MainWindow
	rootWidget *walk.Composite
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
	var tmpCombo *walk.ComboBox
	return StdRootWidget(&customPage.rootWidget,
		pack.TTComposite(pack.NewTTCompositeArgs(nil).
			SetLayoutHBox(true).SetWidgets(
			pack.NewWidgetGroup().Append(
				pack.TTLabel(pack.NewTTLabelArgs(nil).SetText("翻译引擎")),
				pack.TTComboBox(pack.NewTTComboBoxArgs(nil).
					SetModel([]string{"阿里云", "百度云", "腾讯云", "ChatGPT"}).SetCurrentIdx(0)),
				pack.TTLabel(pack.NewTTLabelArgs(nil).SetText("来源语种")),
				pack.TTComboBox(pack.NewTTComboBoxArgs(nil).
					SetModel(ling_va.GetInstance().GetLangSupported()).SetBindingMember("Key").SetDisplayMember("Name").SetCurrentIdx(0)),
				pack.TTLabel(pack.NewTTLabelArgs(nil).SetText("目标语种")),
				pack.TTComboBox(pack.NewTTComboBoxArgs(&tmpCombo).
					SetModel(ling_va.GetInstance().GetLangSupported()).SetBindingMember("Key").SetDisplayMember("Name").SetCurrentIdx(1)),
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
