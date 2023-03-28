package page

import (
	"errors"
	"fmt"
	"github.com/golang-module/carbon"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"sync"
	"translator/domain"
	"translator/task"
	"translator/tst/tt_translator/ling_va"
	"translator/tst/tt_ui/handle"
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
		apiSubripTranslate.chanLog = make(chan string, 128)
		apiSubripTranslate.cronSyncLog()

	})
	return apiSubripTranslate
}

type SubripTranslate struct {
	id          string
	name        string
	chanLog     chan string
	engines     []*_type.StdComboBoxModel
	mainWindow  *walk.MainWindow
	rootWidget  *walk.Composite
	ptrEngine   *walk.ComboBox
	ptrFromLang *walk.ComboBox
	ptrToLang   *walk.ComboBox

	ptrMode       *walk.ComboBox
	ptrMainExport *walk.ComboBox

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
		pack.TTComposite(pack.NewTTCompositeArgs(nil).SetLayoutHBox(true).SetWidgets(
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

		pack.TTComposite(pack.NewTTCompositeArgs(nil).SetLayoutHBox(true).SetWidgets(
			pack.NewWidgetGroup().Append(
				pack.TTLabel(pack.NewTTLabelArgs(nil).SetText("翻译模式")),
				pack.TTComboBox(pack.NewTTComboBoxArgs(&customPage.ptrMode).SetModel(_type.ModeDelta.GetModes()).SetCurrentIdx(_type.ModeDelta.GetIdx())),
				pack.TTLabel(pack.NewTTLabelArgs(nil).SetText("导出主轨道")),
				pack.TTComboBox(pack.NewTTComboBoxArgs(&customPage.ptrMainExport).SetModel(_type.LangDirectionFrom.GetDirections()).SetCurrentIdx(0)),
			).AppendZeroHSpacer().GetWidgets(),
		)),

		StdBrowserSelectorWidget("字幕文件", customPage.eventSrtFileOnClicked, &customPage.ptrSrtFile),
		StdBrowserSelectorWidget("字幕目录", customPage.eventSrtDirOnClicked, &customPage.ptrSrtDir),
		pack.TTComposite(pack.NewTTCompositeArgs(nil).SetLayoutHBox(true).SetWidgets(
			pack.NewWidgetGroup().Append(
				pack.TTPushBtn(pack.NewTTPushBtnArgs(nil).SetText("翻译").SetOnClicked(customPage.eventBtnTranslate)),
				pack.TTPushBtn(pack.NewTTPushBtnArgs(nil).SetText("清空日志").SetOnClicked(customPage.flushLog)),
			).AppendZeroHSpacer().GetWidgets(),
		)),

		pack.TTComposite(pack.NewTTCompositeArgs(nil).SetLayoutHBox(true).SetWidgets(
			pack.NewWidgetGroup().Append(
				pack.TTGroupBox(pack.NewTTGroupBoxArgs(nil).SetVisible(true).SetTitle("日志").SetWidgets(
					pack.NewWidgetGroup().Append(
						pack.TTTextEdit(pack.NewTextEditWrapperArgs(&customPage.ptrLog).SetReadOnly(true).SetVScroll(true)),
					).AppendZeroHSpacer().AppendZeroVSpacer().GetWidgets(),
				)),
			).AppendZeroHSpacer().AppendZeroVSpacer().GetWidgets(),
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
		_ = customPage.ptrMode.SetCurrentIndex(_type.ModeDelta.GetIdx())
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

func (customPage *SubripTranslate) eventSrtFileOnClicked() {
	handle.FileDialogHandle(handle.NewFileDialogHandleArgs(customPage.mainWindow, &customPage.ptrSrtFile).SetTitle("选择字幕文件").SetFilter("*.*|*.srt"))
}

func (customPage *SubripTranslate) eventSrtDirOnClicked() {
	handle.FileDialogHandle(handle.NewFileDialogHandleArgs(customPage.mainWindow, &customPage.ptrSrtDir).SetTitle("选择字幕目录").Folder())
}

func (customPage *SubripTranslate) eventBtnTranslate() {
	currentEngine := domain.GetTranslators().GetById(customPage.engines[customPage.ptrEngine.CurrentIndex()].Key)
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
	no := util.Uid()
	tTranslate := new(task.Translate).
		SetTaskNo(no).
		SetTranslator(currentEngine).
		SetFromLang(fromLang).SetToLang(toLang).
		SetTranslateMode(_type.TranslateMode(mode)).SetMainTrackReport(_type.LangDirection(mainTrackExport)).
		SetSrtFile(strFile).SetSrtDir(strDir).
		SetChanLog(customPage.chanLog)
	if customPage.ptrSrtFile != nil {
		_ = customPage.ptrSrtFile.SetText("")
	}

	if customPage.ptrSrtDir != nil {
		_ = customPage.ptrSrtDir.SetText("")
	}
	msg.Info(customPage.mainWindow, "投递任务成功")
	customPage.appendToLog(fmt.Sprintf(
		"[%s] 投递任务[编号: %s]成功[引擎: %s, 来源语种: %s, 目标语种: %s, 翻译模式: %s, 导出主轨道: %s, 字幕文件: %s, 字幕目录: %s]",
		carbon.Now().Layout(carbon.DateTimeLayout), no, currentEngine.GetName(),
		currentEngine.GetLangSupported()[customPage.ptrFromLang.CurrentIndex()].Name,
		currentEngine.GetLangSupported()[customPage.ptrToLang.CurrentIndex()].Name,
		mode, mainTrackExport, strFile, strDir,
	))

	go func() {
		tTranslate.Run()
	}()
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
