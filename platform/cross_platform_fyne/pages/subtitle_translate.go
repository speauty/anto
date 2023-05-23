package pages

import (
	"anto/common"
	"anto/cron/detector"
	"anto/domain/repository"
	serviceTranslator "anto/domain/service/translator"
	"anto/platform/cross_platform_fyne/msg"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"path/filepath"
	"strings"
	"sync"
)

var (
	apiPageSubtitleTranslate  *PageSubtitleTranslate
	oncePageSubtitleTranslate sync.Once
	trackExport               = []string{"全部轨道", "主字幕轨", "副字幕轨"}
)

func APIPageSubtitleTranslate() *PageSubtitleTranslate {
	oncePageSubtitleTranslate.Do(func() {
		apiPageSubtitleTranslate = &PageSubtitleTranslate{
			id:   "page.subtitle.translate",
			name: "字幕翻译",
			//isDefault:       true,
			translatorNames: repository.GetTranslators().GetNames(),
		}
	})
	return apiPageSubtitleTranslate
}

type PageSubtitleTranslate struct {
	window          fyne.Window
	id              string
	name            string
	isDefault       bool
	translatorNames []string
}

func (page *PageSubtitleTranslate) GetID() string { return page.id }

func (page *PageSubtitleTranslate) GetName() string { return page.name }

func (page *PageSubtitleTranslate) GetWindow() fyne.Window { return page.window }

func (page *PageSubtitleTranslate) SetWindow(win fyne.Window) { page.window = win }

func (page *PageSubtitleTranslate) IsDefault() bool { return page.isDefault }

func (page *PageSubtitleTranslate) OnClose() {}

func (page *PageSubtitleTranslate) OnReset() {}

func (page *PageSubtitleTranslate) OnRender() fyne.CanvasObject {
	var currentEngine serviceTranslator.ImplTranslator

	currentFromLang := binding.NewString()
	currentToLang := binding.NewString()
	selectorFromLang := widget.NewSelectEntry([]string{})
	selectorFromLang.Bind(currentFromLang)
	selectorToLang := widget.NewSelectEntry([]string{})
	selectorToLang.Bind(currentToLang)

	selectorFromLang.Validator = func(currentVal string) error {
		if currentVal == "" {
			return errors.New("请选择来源语种")
		}
		if currentVal == selectorToLang.Text {
			return fmt.Errorf("来源语种(%s)与目标语种(%s)不能为同一语种", currentVal, selectorToLang.Text)
		} else {
			if selectorToLang.Text != "" {
				selectorToLang.SetValidationError(nil)
			}
		}

		return nil
	}

	selectorToLang.Validator = func(currentVal string) error {
		if currentVal == "" {
			return errors.New("请选择目标语种")
		}
		if currentVal == selectorFromLang.Text {
			return fmt.Errorf("目标语种(%s)与来源语种(%s)不能为同一语种", currentVal, selectorFromLang.Text)
		} else {
			if selectorFromLang.Text != "" {
				selectorFromLang.SetValidationError(nil)
			}
		}
		return nil
	}

	currentEngineName := binding.NewString()
	if len(page.translatorNames) > 0 {
		_ = currentEngineName.Set(page.translatorNames[0])
	}

	showCurrentEngineNameLabel := widget.NewLabel("名称: 飞洒发后即可(fs fs)")
	showCurrentEngineDetailLabel := widget.NewLabel("最大长度: 5000  请求QPS: 100  协程数量: 50")

	selectorEngine := widget.NewSelectEntry(page.translatorNames)
	selectorEngine.PlaceHolder = "请选择翻译引擎进行配置"
	selectorEngine.Bind(currentEngineName)
	selectorEngine.Validator = func(currentVal string) error {
		if currentVal == "" {
			return errors.New("请选择当前引擎")
		}
		for _, name := range page.translatorNames {
			if name == currentVal {
				return nil
			}
		}
		return fmt.Errorf("暂无当前引擎: %s", currentVal)
	}

	currentEngineName.AddListener(binding.NewDataListener(func() {
		if err := selectorEngine.Validate(); err != nil {
			currentEngine = nil
			return
		}
		engineName, _ := currentEngineName.Get()
		currentEngine = repository.GetTranslators().GetByName(engineName)
		_ = currentFromLang.Set("")
		_ = currentToLang.Set("")
		var languages []string
		for _, currentLang := range currentEngine.GetLangSupported() {
			if strings.Contains(currentLang.Name, "中文") {
				_ = currentToLang.Set(currentLang.Name)
			}
			if strings.Contains(currentLang.Name, "英语") {
				_ = currentFromLang.Set(currentLang.Name)
			}
			languages = append(languages, currentLang.Name)
		}

		selectorFromLang.SetOptions(languages)
		selectorToLang.SetOptions(languages)

		showCurrentEngineNameLabel.SetText(
			fmt.Sprintf(
				"名称: %s(%s %s)  支持语种数量: %d",
				currentEngine.GetName(), currentEngine.GetId(), currentEngine.GetShortId(),
				len(currentEngine.GetLangSupported()),
			))
		showCurrentEngineDetailLabel.SetText(
			fmt.Sprintf(
				"最大长度: %d  请求QPS: %d  协程数量: %d",
				currentEngine.GetCfg().GetMaxSingleTextLength(), currentEngine.GetCfg().GetQPS(), currentEngine.GetCfg().GetMaxCoroutineNum(),
			))
	}))

	currentTranslateMode := binding.NewString()
	_ = currentTranslateMode.Set(common.ModeFull.String())

	selectorTranslateMode := widget.NewSelectEntry(common.ModeFull.GetModes())
	selectorTranslateMode.Bind(currentTranslateMode)
	selectorTranslateMode.Validator = func(currentVal string) error {
		if currentVal == "" {
			return errors.New("请选择翻译模式")
		}
		for _, name := range common.ModeFull.GetModes() {
			if name == currentVal {
				return nil
			}
		}
		return fmt.Errorf("暂无当前翻译模式: %s", currentVal)
	}

	currentTrackerExport := binding.NewString()
	_ = currentTrackerExport.Set(trackExport[0])
	selectorTrackerExport := widget.NewSelectEntry(trackExport)
	selectorTrackerExport.Bind(currentTrackerExport)
	selectorTrackerExport.Validator = func(currentVal string) error {
		if currentVal == "" {
			return errors.New("请选择导出轨道")
		}
		for _, name := range trackExport {
			if name == currentVal {
				return nil
			}
		}
		return fmt.Errorf("暂无当前导出轨道: %s", currentVal)
	}

	currentMainTrackerLang := binding.NewString()
	_ = currentMainTrackerLang.Set(common.LangDirectionTo.String())
	selectorMainTrackerLang := widget.NewSelectEntry(common.LangDirectionTo.GetDirections())
	selectorMainTrackerLang.Bind(currentMainTrackerLang)
	selectorMainTrackerLang.Validator = func(currentVal string) error {
		if currentVal == "" {
			return errors.New("请选择主轨语种")
		}
		for _, name := range common.LangDirectionTo.GetDirections() {
			if name == currentVal {
				return nil
			}
		}
		return fmt.Errorf("暂无当前主轨语种: %s", currentVal)
	}

	currentSrtFile := binding.NewString()
	showCurrentSrtFileLabel := widget.NewLabel("当前字幕文件: 暂无")
	showCurrentSrtFileLabel.Wrapping = fyne.TextWrapBreak
	currentSrtFile.AddListener(binding.NewDataListener(func() {
		currentSrtFileVal, _ := currentSrtFile.Get()
		if currentSrtFileVal == "" {
			showCurrentSrtFileLabel.SetText("当前字幕文件: 暂无")
		} else {
			showCurrentSrtFileLabel.SetText(fmt.Sprintf("当前字幕文件: %s", filepath.Base(currentSrtFileVal)))
		}
	}))
	handleFile := dialog.NewFileOpen(func(uriReadCloser fyne.URIReadCloser, err error) {
		if err != nil {
			msg.Error(page.GetWindow(), err)
			return
		}
		if uriReadCloser == nil {
			_ = currentSrtFile.Set("")
			return
		}
		_ = currentSrtFile.Set(uriReadCloser.URI().Path())
	}, page.GetWindow())
	handleFile.SetConfirmText("选择")
	handleFile.SetDismissText("取消")
	handleFile.SetFilter(storage.NewExtensionFileFilter([]string{".srt"}))

	currentSrtFolder := binding.NewString()
	showCurrentSrtFolderLabel := widget.NewLabel("当前字幕目录: 暂无")
	showCurrentSrtFolderLabel.Wrapping = fyne.TextWrapBreak
	currentSrtFolder.AddListener(binding.NewDataListener(func() {
		currentSrtFolderVal, _ := currentSrtFolder.Get()
		if currentSrtFolderVal == "" {
			showCurrentSrtFolderLabel.SetText("当前字幕目录: 暂无")
		} else {
			showCurrentSrtFolderLabel.SetText(fmt.Sprintf("当前字幕目录: %s", currentSrtFolderVal))
		}
	}))
	handleFolder := dialog.NewFolderOpen(func(listUri fyne.ListableURI, err error) {
		if err != nil {
			msg.Error(page.GetWindow(), err)
			return
		}
		if listUri == nil {
			_ = currentSrtFolder.Set("")
			return
		}
		_ = currentSrtFolder.Set(listUri.Path())
	}, page.GetWindow())
	handleFolder.SetConfirmText("选择")
	handleFolder.SetDismissText("取消")

	configForm := widget.NewForm(
		&widget.FormItem{Text: "当前引擎", Widget: selectorEngine, HintText: "当前翻译使用引擎"},
		&widget.FormItem{Text: "翻译模式", Widget: selectorTranslateMode, HintText: "增量翻译只处理未翻译的字幕行"},
		&widget.FormItem{Text: "来源语种", Widget: selectorFromLang, HintText: "当前字幕语种"},
		&widget.FormItem{Text: "目标语种", Widget: selectorToLang, HintText: "目标字幕语种"},
		&widget.FormItem{Text: "主轨语种", Widget: selectorMainTrackerLang, HintText: "将选择语种导出到主轨道上"},
		&widget.FormItem{Text: "导出轨道", Widget: selectorTrackerExport},
		&widget.FormItem{Text: "字幕文件", Widget: widget.NewButton("选择文件", func() {
			handleFile.Show()
		})},
		&widget.FormItem{Text: "字幕目录", Widget: widget.NewButton("选择目录", func() {
			handleFolder.Show()
		}), HintText: "不支持递归目录搜集字幕文件"},
	)

	configForm.SubmitText = "提交"
	configForm.OnSubmit = func() {
		fromLang := ""
		toLang := ""
		mode := ""
		mainTrackExport := ""
		srtFile := ""
		srtDir := ""
		IdxTrackExport := 0

		currentFromLangVal, _ := currentFromLang.Get()
		currentToLangVal, _ := currentToLang.Get()
		for _, currentLang := range currentEngine.GetLangSupported() {
			if fromLang == "" && currentLang.Name == currentFromLangVal {
				fromLang = currentLang.Key
			}
			if toLang == "" && currentLang.Name == currentToLangVal {
				toLang = currentLang.Key
			}
			if fromLang != "" && toLang != "" {
				break
			}
		}

		if fromLang == "" || toLang == "" {
			msg.Error(page.GetWindow(), errors.New("请重新选择来源语种或目标语种"))
			return
		}

		mode, _ = currentTranslateMode.Get()
		mainTrackExport, _ = currentMainTrackerLang.Get()

		srtFile, _ = currentSrtFile.Get()
		srtDir, _ = currentSrtFolder.Get()

		if srtFile == "" && srtDir == "" {
			msg.Error(page.GetWindow(), errors.New("请选择字幕文件或目录"))
			return
		}

		currentTrackExport, _ := currentTrackerExport.Get()
		for idx, track := range trackExport {
			if currentTrackExport == track {
				IdxTrackExport = idx
				break
			}
		}

		detector.Singleton().Push(&detector.StrDetectorData{
			Translator: currentEngine, FromLang: fromLang, ToLang: toLang,
			TranslateMode: common.TranslateMode(mode), MainTrackReport: common.LangDirection(mainTrackExport),
			SrtFile: srtFile, SrtDir: srtDir, FlagTrackExport: IdxTrackExport,
		})

		msg.Info(page.GetWindow(), "提交翻译任务成功, 打开控制台(文件 => 控制台)可查看详细信息", "", func() {
			_ = currentSrtFile.Set("")
			_ = currentSrtFolder.Set("")
		})
	}

	return container.NewVBox(
		pageTitle(page.GetName()),
		widget.NewSeparator(),
		container.NewHBox(
			container.NewGridWrap(
				fyne.NewSize(400, 500),
				configForm,
			),
			widget.NewSeparator(),
			container.NewGridWrap(
				fyne.NewSize(300, 30),
				widget.NewLabel("当前引擎信息"),
				showCurrentEngineNameLabel,
				showCurrentEngineDetailLabel,
				showCurrentSrtFileLabel,
				showCurrentSrtFolderLabel,
			),
		),
	)
}
