package page

import (
	"anto/common"
	"anto/cron/detector"
	"anto/cron/translate"
	"anto/domain/repository"
	"anto/domain/service/translator"
	"anto/domain/service/translator/ling_va"
	"anto/lib/util"
	"anto/platform/win/ui/handle"
	"anto/platform/win/ui/msg"
	"anto/platform/win/ui/pack"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var (
	apiSubtitleTranslate  *SubtitleTranslate
	onceSubtitleTranslate sync.Once
	comboBoxModel         = new(common.StdComboBoxModel)
	trackExportedOptions  = []string{"全部轨道", "主轨", "副轨"}
	modeTranslate         = common.ModeDelta.GetModes()
	languageDirections    = common.LangDirectionFrom.GetDirections()
)

func GetSubripTranslate() *SubtitleTranslate {
	onceSubtitleTranslate.Do(func() {
		apiSubtitleTranslate = new(SubtitleTranslate)
		apiSubtitleTranslate.id = util.Uid()
		apiSubtitleTranslate.name = "字幕翻译"
		apiSubtitleTranslate.engines = repository.GetTranslators().GetNames()
		apiSubtitleTranslate.chanLog = make(chan string, 128)
		detector.Singleton().SetMsgRedirect(apiSubtitleTranslate.chanLog)
		translate.Singleton().SetMsgRedirect(apiSubtitleTranslate.chanLog)

		apiSubtitleTranslate.cronSyncLog()

	})
	return apiSubtitleTranslate
}

type SubtitleTranslate struct {
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

	ptrConfigView *walk.TextLabel
	ptrLog        *walk.TextEdit

	dropFilesEventId int
}

func (customPage *SubtitleTranslate) GetId() string {
	return customPage.id
}

func (customPage *SubtitleTranslate) GetName() string {
	return customPage.name
}

func (customPage *SubtitleTranslate) BindWindow(win *walk.MainWindow) {
	customPage.mainWindow = win
}

func (customPage *SubtitleTranslate) SetVisible(isVisible bool) {
	if customPage.rootWidget != nil {
		customPage.rootWidget.SetVisible(isVisible)
	}
}

func (customPage *SubtitleTranslate) GetWidget() Widget {
	widgets := []Widget{}
	widgets = append(widgets, pack.UIComposite(pack.NewUICompositeArgs(nil).SetLayoutHBox(false).SetWidgets(
		pack.NewWidgetGroup().Append(
			pack.UIComposite(pack.NewUICompositeArgs(nil).SetLayoutVBox(true).SetWidgets(
				pack.NewWidgetGroup().Append(customPage.getFormWidget()...).AppendZeroVSpacer().AppendZeroHSpacer().GetWidgets(),
			)),
			pack.UIComposite(pack.NewUICompositeArgs(nil).
				SetLayoutVBox(true).SetAlignment(AlignHNearVNear).
				SetCustomSize(Size{Width: 200}).SetWidgets(
				pack.NewWidgetGroup().Append(customPage.getConfigWidget()...).AppendZeroHSpacer().AppendZeroVSpacer().GetWidgets(),
			)),
		).AppendZeroHSpacer().GetWidgets(),
	)))

	widgets = append(widgets, customPage.getConsoleWidget(), VSpacer{})
	return StdRootWidget(&customPage.rootWidget, widgets...)
}

func (customPage *SubtitleTranslate) Reset() {
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
		_ = customPage.ptrMainExport.SetCurrentIndex(1)
	}

	if customPage.ptrSrtFile != nil {
		_ = customPage.ptrSrtFile.SetText("")
	}

	if customPage.ptrSrtDir != nil {
		_ = customPage.ptrSrtDir.SetText("")
	}

	if customPage.dropFilesEventId == 0 {
		dfEvent := customPage.rootWidget.DropFiles()
		customPage.dropFilesEventId = dfEvent.Attach(func(files []string) {
			if len(files) < 1 {
				return
			}
			currentFile := files[0]
			fsInfo, err := os.Stat(currentFile)
			if err != nil {
				return
			}
			if fsInfo.IsDir() {
				_ = customPage.ptrSrtDir.SetText(currentFile)
			} else {
				if util.IsSrtFile(currentFile) {
					_ = customPage.ptrSrtFile.SetText(currentFile)
				} else {
					// @todo 触发两次弹窗提示，暂时注释
					// msg.Err(customPage.mainWindow, fmt.Errorf("无效字幕文件(%s)", filepath.Base(currentFile)))
				}
			}
		})
	}
}

func (customPage *SubtitleTranslate) getConfigWidget() []Widget {
	return []Widget{
		pack.UILabel(pack.NewUILabelArgs(nil).SetText("友情提示: 支持文件和目录拖放哦")),
		pack.UILabel(pack.NewUILabelArgs(nil).SetText("引擎信息")),
		pack.UITextLabel(pack.NewUITextLabelArgs(&customPage.ptrConfigView).SetText("暂无")),
	}
}

func (customPage *SubtitleTranslate) getFormWidget() []Widget {
	return []Widget{
		pack.UIComposite(pack.NewUICompositeArgs(nil).SetLayoutHBox(true).SetWidgets(
			pack.NewWidgetGroup().Append(
				pack.UILabel(pack.NewUILabelArgs(nil).SetText("翻译引擎")),
				pack.UIComboBox(pack.NewUIComboBoxArgs(&customPage.ptrEngine).
					SetModel(apiSubtitleTranslate.engines).
					SetBindingMember(comboBoxModel.BindKey()).
					SetDisplayMember(comboBoxModel.DisplayKey()).
					SetCurrentIdx(-1).
					SetOnCurrentIdxChangedFn(customPage.eventEngineOnChange),
				),

				pack.UILabel(pack.NewUILabelArgs(nil).SetText(common.LangDirectionFrom.String())),
				pack.UIComboBox(pack.NewUIComboBoxArgs(&customPage.ptrFromLang).
					SetModel(ling_va.API().GetLangSupported()).SetBindingMember(comboBoxModel.BindKey()).SetDisplayMember(comboBoxModel.DisplayKey())),

				pack.UILabel(pack.NewUILabelArgs(nil).SetText(common.LangDirectionTo.String())),
				pack.UIComboBox(pack.NewUIComboBoxArgs(&customPage.ptrToLang).
					SetModel(ling_va.API().GetLangSupported()).SetBindingMember(comboBoxModel.BindKey()).SetDisplayMember(comboBoxModel.DisplayKey()).SetCurrentIdx(1)),
			).AppendZeroHSpacer().GetWidgets(),
		)),

		pack.UIComposite(pack.NewUICompositeArgs(nil).SetLayoutHBox(true).SetWidgets(
			pack.NewWidgetGroup().Append(
				pack.UILabel(pack.NewUILabelArgs(nil).SetText("翻译模式")),
				pack.UIComboBox(pack.NewUIComboBoxArgs(&customPage.ptrMode).SetModel(modeTranslate).SetCurrentIdx(common.ModeDelta.GetIdx())),
				pack.UILabel(pack.NewUILabelArgs(nil).SetText("导出轨道")),
				pack.UIComboBox(pack.NewUIComboBoxArgs(&customPage.ptrFlagTrackExport).SetModel(trackExportedOptions).SetCurrentIdx(0)),
				pack.UILabel(pack.NewUILabelArgs(nil).SetText("导出主轨")),
				pack.UIComboBox(pack.NewUIComboBoxArgs(&customPage.ptrMainExport).SetModel(languageDirections).SetCurrentIdx(1)),
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
	}
}

func (customPage *SubtitleTranslate) getConsoleWidget() Widget {
	return pack.UIComposite(pack.NewUICompositeArgs(nil).SetLayoutVBox(false).SetWidgets(
		pack.NewWidgetGroup().Append(
			pack.UITextEdit(pack.NewUITextEditArgs(&customPage.ptrLog).SetCustomSize(Size{Height: 344}).SetReadOnly(true).SetVScroll(true)),
		).AppendZeroHSpacer().GetWidgets(),
	))
}

func (customPage *SubtitleTranslate) eventEngineOnChange() {
	_ = customPage.ptrConfigView.SetText("")
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
	configViewText := fmt.Sprintf(
		"名称: %s(%s), 缩写: %s"+
			"\n引擎配置\n",
		currentEngine.GetName(), currentEngine.GetId(), currentEngine.GetShortId(),
	)
	if currentEngine.GetCfg().GetAK() != translator.ConfigInvalidStr {
		configViewText += fmt.Sprintf("应用: %s\n", currentEngine.GetCfg().GetAK())
	}
	if currentEngine.GetCfg().GetSK() != translator.ConfigInvalidStr {
		configViewText += fmt.Sprintf("密钥: %s\n", currentEngine.GetCfg().GetSK())
	}
	if currentEngine.GetCfg().GetProjectKey() != translator.ConfigInvalidStr {
		configViewText += fmt.Sprintf("项目: %s, ", currentEngine.GetCfg().GetProjectKey())
	}
	if currentEngine.GetCfg().GetRegion() != translator.ConfigInvalidStr {
		configViewText += fmt.Sprintf("地区: %s\n", currentEngine.GetCfg().GetRegion())
	}
	configViewText += fmt.Sprintf(
		"最长字符: %d, QPS: %d, 协程数量: %d\n支持语种数量: %d",
		currentEngine.GetCfg().GetMaxCharNum(), currentEngine.GetCfg().GetQPS(),
		currentEngine.GetCfg().GetMaxCoroutineNum(),
		len(currentEngine.GetLangSupported()),
	)
	_ = customPage.ptrConfigView.SetText(configViewText)
	return
}

func (customPage *SubtitleTranslate) eventSrtFileOnClicked() {
	handle.FileDialogHandle(handle.NewFileDialogHandleArgs(customPage.mainWindow, &customPage.ptrSrtFile).SetTitle("选择字幕文件").SetFilter("*.*|*.srt"))
}

func (customPage *SubtitleTranslate) eventSrtDirOnClicked() {
	handle.FileDialogHandle(handle.NewFileDialogHandleArgs(customPage.mainWindow, &customPage.ptrSrtDir).SetTitle("选择字幕目录").Folder())
}

func (customPage *SubtitleTranslate) eventBtnTranslate() {
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

func (customPage *SubtitleTranslate) setLangComboBox(ptr *walk.ComboBox, model interface{}, idx int) {
	_ = ptr.SetModel(model)
	_ = ptr.SetCurrentIndex(idx)
}

func (customPage *SubtitleTranslate) appendToLog(msg string) {
	if customPage.ptrLog.Text() == "" {
		_ = customPage.ptrLog.SetText(fmt.Sprintf("%s", msg))
	} else {
		_ = customPage.ptrLog.SetText(fmt.Sprintf("%s\r\n%s", msg, customPage.ptrLog.Text()))
	}
}

func (customPage *SubtitleTranslate) flushLog() {
	_ = customPage.ptrLog.SetText("")
}

func (customPage *SubtitleTranslate) cronSyncLog() {
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
