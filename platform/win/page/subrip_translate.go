package page

import (
	"anto/common"
	"anto/cron/detector"
	"anto/cron/translate"
	"anto/dependency/repository"
	"anto/dependency/service/translator/ling_va"
	"anto/lib/util"
	"anto/platform/win/ui/handle"
	"anto/platform/win/ui/msg"
	"anto/platform/win/ui/pack"
	"errors"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"sync"
)

var (
	apiSubripTranslate  *SubripTranslate
	onceSubripTranslate sync.Once
	comboBoxModel       = new(common.StdComboBoxModel)
)

func GetSubripTranslate() *SubripTranslate {
	onceSubripTranslate.Do(func() {
		apiSubripTranslate = new(SubripTranslate)
		apiSubripTranslate.id = util.Uid()
		apiSubripTranslate.name = "字幕翻译"
		apiSubripTranslate.engines = repository.GetTranslators().GetNames()
		apiSubripTranslate.chanLog = make(chan string, 128)
		detector.Singleton().SetMsgRedirect(apiSubripTranslate.chanLog)
		translate.Singleton().SetMsgRedirect(apiSubripTranslate.chanLog)

		apiSubripTranslate.cronSyncLog()

	})
	return apiSubripTranslate
}

type SubripTranslate struct {
	id          string
	name        string
	chanLog     chan string
	engines     []*common.StdComboBoxModel
	mainWindow  *walk.MainWindow
	rootWidget  *walk.Composite
	ptrEngine   *walk.ComboBox
	ptrFromLang *walk.ComboBox
	ptrToLang   *walk.ComboBox

	ptrMode       *walk.ComboBox
	ptrMainExport *walk.ComboBox

	ptrFlagTrackExport *walk.ComboBox

	ptrSrtFile *walk.Label
	ptrSrtDir  *walk.Label

	ptrLog *walk.TextEdit
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
		pack.UIComposite(pack.NewUICompositeArgs(nil).SetLayoutHBox(true).SetWidgets(
			pack.NewWidgetGroup().Append(
				pack.UILabel(pack.NewUILabelArgs(nil).SetText("翻译引擎")),
				pack.UIComboBox(pack.NewUIComboBoxArgs(&customPage.ptrEngine).
					SetModel(apiSubripTranslate.engines).
					SetBindingMember(comboBoxModel.BindKey()).
					SetDisplayMember(comboBoxModel.DisplayKey()).
					SetCurrentIdx(0).
					SetOnCurrentIdxChangedFn(customPage.eventEngineOnChange),
				),

				pack.UILabel(pack.NewUILabelArgs(nil).SetText(common.LangDirectionFrom.String())),
				pack.UIComboBox(pack.NewUIComboBoxArgs(&customPage.ptrFromLang).
					SetModel(ling_va.Singleton().GetLangSupported()).SetBindingMember(comboBoxModel.BindKey()).SetDisplayMember(comboBoxModel.DisplayKey())),

				pack.UILabel(pack.NewUILabelArgs(nil).SetText("    "+common.LangDirectionTo.String())),
				pack.UIComboBox(pack.NewUIComboBoxArgs(&customPage.ptrToLang).
					SetModel(ling_va.Singleton().GetLangSupported()).SetBindingMember(comboBoxModel.BindKey()).SetDisplayMember(comboBoxModel.DisplayKey()).SetCurrentIdx(1)),
			).AppendZeroHSpacer().GetWidgets(),
		)),

		pack.UIComposite(pack.NewUICompositeArgs(nil).SetLayoutHBox(true).SetWidgets(
			pack.NewWidgetGroup().Append(
				pack.UILabel(pack.NewUILabelArgs(nil).SetText("翻译模式")),
				pack.UIComboBox(pack.NewUIComboBoxArgs(&customPage.ptrMode).SetModel(common.ModeDelta.GetModes()).SetCurrentIdx(common.ModeDelta.GetIdx())),
				pack.UILabel(pack.NewUILabelArgs(nil).SetText("导出轨道")),
				pack.UIComboBox(pack.NewUIComboBoxArgs(&customPage.ptrFlagTrackExport).SetModel([]string{"全部轨道", "主轨", "副轨"}).SetCurrentIdx(0)),
				pack.UILabel(pack.NewUILabelArgs(nil).SetText("导出主轨道")),
				pack.UIComboBox(pack.NewUIComboBoxArgs(&customPage.ptrMainExport).SetModel(common.LangDirectionFrom.GetDirections()).SetCurrentIdx(1)),
			).AppendZeroHSpacer().GetWidgets(),
		)),

		StdBrowserSelectorWidget("字幕文件", customPage.eventSrtFileOnClicked, &customPage.ptrSrtFile),
		StdBrowserSelectorWidget("字幕目录", customPage.eventSrtDirOnClicked, &customPage.ptrSrtDir),
		pack.UIComposite(pack.NewUICompositeArgs(nil).SetLayoutHBox(true).SetWidgets(
			pack.NewWidgetGroup().Append(
				pack.UIPushBtn(pack.NewUIPushBtnArgs(nil).SetText("翻译").SetOnClicked(customPage.eventBtnTranslate)),
				pack.UIPushBtn(pack.NewUIPushBtnArgs(nil).SetText("清空输出").SetOnClicked(customPage.flushLog)),
			).AppendZeroHSpacer().GetWidgets(),
		)),

		pack.UIComposite(pack.NewUICompositeArgs(nil).SetLayoutVBox(false).SetWidgets(
			pack.NewWidgetGroup().Append(
				pack.UITextEdit(pack.NewUITextEditArgs(&customPage.ptrLog).SetCustomSize(Size{Height: 344}).SetReadOnly(true).SetVScroll(true)),
			).AppendZeroHSpacer().GetWidgets(),
		)),

		VSpacer{},
	)
}

func (customPage *SubripTranslate) Reset() {
	if customPage.ptrEngine != nil {
		_ = customPage.ptrEngine.SetCurrentIndex(0)
	}

	if customPage.ptrFromLang != nil {
		_ = customPage.ptrFromLang.SetCurrentIndex(1)
	}

	if customPage.ptrToLang != nil {
		_ = customPage.ptrToLang.SetCurrentIndex(0)
	}

	if customPage.ptrMode != nil {
		_ = customPage.ptrMode.SetCurrentIndex(common.ModeDelta.GetIdx())
	}

	if customPage.ptrMainExport != nil {
		_ = customPage.ptrMainExport.SetCurrentIndex(0)
	}

	if customPage.ptrSrtFile != nil {
		_ = customPage.ptrSrtFile.SetText("")
	}

	if customPage.ptrSrtDir != nil {
		_ = customPage.ptrSrtDir.SetText("")
	}
}

func (customPage *SubripTranslate) eventEngineOnChange() {
	currentId := customPage.engines[customPage.ptrEngine.CurrentIndex()].Key
	if currentId == "" {
		msg.Err(customPage.mainWindow, fmt.Errorf("当前翻译引擎无效，请重新选择"))
		_ = customPage.ptrEngine.SetCurrentIndex(-1)
		customPage.setLangComboBox(customPage.ptrFromLang, nil, -1)
		customPage.setLangComboBox(customPage.ptrToLang, nil, -1)
		return
	}
	currentEngine := repository.GetTranslators().GetById(currentId)
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

func (customPage *SubripTranslate) eventSrtFileOnClicked() {
	handle.FileDialogHandle(handle.NewFileDialogHandleArgs(customPage.mainWindow, &customPage.ptrSrtFile).SetTitle("选择字幕文件").SetFilter("*.*|*.srt"))
}

func (customPage *SubripTranslate) eventSrtDirOnClicked() {
	handle.FileDialogHandle(handle.NewFileDialogHandleArgs(customPage.mainWindow, &customPage.ptrSrtDir).SetTitle("选择字幕目录").Folder())
}

func (customPage *SubripTranslate) eventBtnTranslate() {
	currentEngine := repository.GetTranslators().GetById(customPage.engines[customPage.ptrEngine.CurrentIndex()].Key)
	fromLang := currentEngine.GetLangSupported()[customPage.ptrFromLang.CurrentIndex()].Key
	toLang := currentEngine.GetLangSupported()[customPage.ptrToLang.CurrentIndex()].Key
	mode := customPage.ptrMode.Text()
	mainTrackExport := customPage.ptrMainExport.Text()
	strFile := customPage.ptrSrtFile.Text()
	strDir := customPage.ptrSrtDir.Text()

	if fromLang == "" {
		msg.Err(customPage.mainWindow, errors.New("请选择来源语种"))
		return
	}
	if toLang == "" {
		msg.Err(customPage.mainWindow, errors.New("请选择目标语种"))
		return
	}
	if fromLang == toLang {
		msg.Err(customPage.mainWindow, errors.New("来源语种和目标语种相同，请重新选择"))
		return
	}
	if strFile == "" && strDir == "" {
		msg.Err(customPage.mainWindow, errors.New("请选择字幕文件或目录, 优先使用字幕文件"))
		return
	}
	detector.Singleton().Push(&detector.StrDetectorData{
		Translator: currentEngine, FromLang: fromLang, ToLang: toLang,
		TranslateMode: common.TranslateMode(mode), MainTrackReport: common.LangDirection(mainTrackExport),
		SrtFile: strFile, SrtDir: strDir, FlagTrackExport: customPage.ptrFlagTrackExport.CurrentIndex(),
	})

	if customPage.ptrSrtFile != nil {
		_ = customPage.ptrSrtFile.SetText("")
	}

	if customPage.ptrSrtDir != nil {
		_ = customPage.ptrSrtDir.SetText("")
	}
	msg.Info(customPage.mainWindow, "投递任务成功")

	return
}

func (customPage *SubripTranslate) setLangComboBox(ptr *walk.ComboBox, model interface{}, idx int) {
	_ = ptr.SetModel(model)
	_ = ptr.SetCurrentIndex(idx)
}

func (customPage *SubripTranslate) appendToLog(msg string) {
	if customPage.ptrLog.Text() == "" {
		_ = customPage.ptrLog.SetText(fmt.Sprintf("%s", msg))
	} else {
		_ = customPage.ptrLog.SetText(fmt.Sprintf("%s\r\n%s", msg, customPage.ptrLog.Text()))
	}
}

func (customPage *SubripTranslate) flushLog() {
	_ = customPage.ptrLog.SetText("")
}

func (customPage *SubripTranslate) cronSyncLog() {
	go func() {
		for true {
			select {
			case logSrt := <-customPage.chanLog:
				if logSrt != "" {
					customPage.appendToLog(logSrt)
				}
			}
		}
	}()
}
