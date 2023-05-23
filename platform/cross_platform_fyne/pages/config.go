package pages

import (
	"anto/domain/repository"
	serviceTranslator "anto/domain/service/translator"
	"anto/platform/cross_platform_fyne/msg"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"strconv"
	"sync"
)

var (
	apiPageConfig  *PageConfig
	oncePageConfig sync.Once
)

func APIPageConfig() *PageConfig {
	oncePageConfig.Do(func() {
		apiPageConfig = &PageConfig{
			id:              "page.config",
			name:            "设置",
			translatorNames: repository.GetTranslators().GetAllNames(),
		}
	})
	return apiPageConfig
}

type PageConfig struct {
	window          fyne.Window
	id              string
	name            string
	isDefault       bool
	translatorNames []string
}

func (page *PageConfig) GetID() string { return page.id }

func (page *PageConfig) GetName() string { return page.name }

func (page *PageConfig) GetWindow() fyne.Window { return page.window }

func (page *PageConfig) SetWindow(win fyne.Window) { page.window = win }

func (page *PageConfig) IsDefault() bool { return page.isDefault }

func (page *PageConfig) OnClose() {}

func (page *PageConfig) OnReset() {}

func (page *PageConfig) OnRender() fyne.CanvasObject {
	var currentEngine serviceTranslator.ImplTranslator

	currentEngineName := binding.NewString()

	if len(page.translatorNames) > 0 {
		_ = currentEngineName.Set(page.translatorNames[0])
	}
	currentEngineId := binding.NewString()
	currentEngineAppKey := binding.NewString()
	currentEngineAppSecret := binding.NewString()
	currentEngineProjectKey := binding.NewString()
	currentEngineMaxLength := binding.NewString()
	currentEngineQPS := binding.NewString()
	currentEngineProcMax := binding.NewString()

	entryAppKey := widget.NewEntryWithData(currentEngineAppKey)
	entryAppSecret := widget.NewEntryWithData(currentEngineAppSecret)
	entryProjectKey := widget.NewEntryWithData(currentEngineProjectKey)

	selectorEngine := widget.NewSelectEntry(page.translatorNames)
	selectorEngine.PlaceHolder = "请选择翻译引擎进行配置"
	selectorEngine.Bind(currentEngineName)
	currentEngineName.AddListener(binding.NewDataListener(func() {
		if err := selectorEngine.Validate(); err != nil {
			currentEngine = nil
			return
		}
		engineName, _ := currentEngineName.Get()
		currentEngine = repository.GetTranslators().GetByName(engineName)
		_ = currentEngineId.Set(currentEngine.GetId())
		_ = currentEngineMaxLength.Set(fmt.Sprintf("%d", currentEngine.GetCfg().GetMaxSingleTextLength()))
		_ = currentEngineQPS.Set(fmt.Sprintf("%d", currentEngine.GetCfg().GetQPS()))
		_ = currentEngineProcMax.Set(fmt.Sprintf("%d", currentEngine.GetCfg().GetMaxCoroutineNum()))
		_ = currentEngineAppKey.Set(currentEngine.GetCfg().GetAK())
		_ = currentEngineAppSecret.Set(currentEngine.GetCfg().GetSK())
		_ = currentEngineProjectKey.Set(currentEngine.GetCfg().GetPK())
	}))
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

	entryMaxLength := widget.NewEntryWithData(currentEngineMaxLength)
	entryMaxLength.Validator = func(currentVal string) error {
		parseInt, err := strconv.Atoi(currentVal)
		if err != nil {
			return fmt.Errorf("无效数字: %s", currentVal)
		}
		if parseInt <= 0 {
			return fmt.Errorf("最大长度必须大于0: %s", currentVal)
		}
		return nil
	}

	entryQPS := widget.NewEntryWithData(currentEngineQPS)
	entryQPS.Validator = func(currentVal string) error {
		parseInt, err := strconv.Atoi(currentVal)
		if err != nil {
			return fmt.Errorf("无效数字: %s", currentVal)
		}
		if parseInt < 0 {
			return fmt.Errorf("QPS不能小于0: %s", currentVal)
		}
		return nil
	}

	entryProcMax := widget.NewEntryWithData(currentEngineProcMax)
	entryProcMax.Validator = func(currentVal string) error {
		parseInt, err := strconv.Atoi(currentVal)
		if err != nil {
			return fmt.Errorf("无效数字: %s", currentVal)
		}
		if parseInt < 0 {
			return fmt.Errorf("协程数量不能小于0: %s", currentVal)
		}
		return nil
	}

	configForm := widget.NewForm(
		&widget.FormItem{Text: "当前引擎", Widget: selectorEngine, HintText: "当前需要修改配置的引擎"},
		&widget.FormItem{Text: "应用标识", Widget: entryAppKey, HintText: "部分引擎无该参数"},
		&widget.FormItem{Text: "应用密钥", Widget: entryAppSecret, HintText: "部分引擎无该参数"},
		&widget.FormItem{Text: "项目标识", Widget: entryProjectKey, HintText: "部分引擎无该参数"},
		&widget.FormItem{Text: "最大长度", Widget: entryMaxLength, HintText: "单次请求最大字符数量(含标点符号)"},
		&widget.FormItem{Text: "请求QPS", Widget: entryQPS, HintText: "每秒最大请求数量, 该值有一定误差, 建议小于标准值"},
		&widget.FormItem{Text: "协程数量", Widget: entryProcMax, HintText: "翻译并行程序数量(合理即可, 该值并不是越大越好)"},
	)

	configForm.SubmitText = "保存"
	configForm.OnSubmit = func() {
		appKey, _ := currentEngineAppKey.Get()
		appSecret, _ := currentEngineAppSecret.Get()
		projectKey, _ := currentEngineProjectKey.Get()
		maxLengthStr, _ := currentEngineMaxLength.Get()
		qpsStr, _ := currentEngineQPS.Get()
		coroutineNumStr, _ := currentEngineProcMax.Get()
		maxLength, _ := strconv.Atoi(maxLengthStr)
		qps, _ := strconv.Atoi(qpsStr)
		coroutineNum, _ := strconv.Atoi(coroutineNumStr)

		_ = currentEngine.GetCfg().SetAK(appKey)
		_ = currentEngine.GetCfg().SetSK(appSecret)
		_ = currentEngine.GetCfg().SetPK(projectKey)
		_ = currentEngine.GetCfg().SetMaxSingleTextLength(maxLength)
		_ = currentEngine.GetCfg().SetQPS(qps)
		_ = currentEngine.GetCfg().SetMaxCoroutineNum(coroutineNum)
		if err := currentEngine.GetCfg().Sync(); err != nil {
			msg.Error(page.GetWindow(), err)
			return
		}

		msg.Info(page.GetWindow(), "更新配置成功, 建议重启应用(全局配置暂未同步, 可能会造成干扰)", "", nil)
	}

	return container.NewVBox(
		pageTitle(page.GetName()),
		configForm,
	)
}
