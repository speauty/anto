package pages

import (
	"anto/domain/repository"
	serviceTranslator "anto/domain/service/translator"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"reflect"
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
			isDefault:       true,
			translatorNames: repository.GetTranslators().GetNames(),
		}
	})
	return apiPageConfig
}

type PageConfig struct {
	id              string
	name            string
	isDefault       bool
	translatorNames []string
}

func (page *PageConfig) GetID() string {
	return page.id
}

func (page *PageConfig) GetName() string {
	return page.name
}

func (page *PageConfig) IsDefault() bool {
	return page.isDefault
}

func (page *PageConfig) OnClose() {}

func (page *PageConfig) OnReset() {}

func (page *PageConfig) OnRender() fyne.CanvasObject {
	//globalConfig := cfg.Singleton()
	var currentEngine serviceTranslator.ImplTranslator
	configForm := widget.NewForm()

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

	fiAppKey := &widget.FormItem{Text: "应用标识", Widget: entryAppKey}
	fiAppSecret := &widget.FormItem{Text: "应用密钥", Widget: entryAppSecret}
	fiProjectKey := &widget.FormItem{Text: "项目标识", Widget: entryProjectKey}

	selectorEngine := widget.NewSelectEntry(page.translatorNames)
	selectorEngine.PlaceHolder = "请选择翻译引擎进行配置"
	configForm.AppendItem(&widget.FormItem{Text: "当前引擎", Widget: selectorEngine, HintText: "当前需要修改配置的引擎"})
	selectorEngine.Bind(currentEngineName)
	currentEngineName.AddListener(binding.NewDataListener(func() {
		defer configForm.Refresh()
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

		_ = currentEngineAppKey.Set("")
		_ = currentEngineAppSecret.Set("")
		_ = currentEngineProjectKey.Set("")
		entryAppKey.Disable()
		entryAppSecret.Disable()
		entryProjectKey.Disable()
		fmt.Println(reflect.TypeOf(currentEngine.GetCfg()))
		//switch currentEngine.GetCfg().(type) {
		//case *ling_va.Cfg:
		//	_ = currentEngineAppKey.Set(currentEngine.GetCfg().(*ling_va.Cfg).DataId)
		//	entryAppKey.Enable()
		//case *huawei_cloud_nlp.Cfg:
		//	_ = currentEngineAppKey.Set(currentEngine.GetCfg().(*huawei_cloud_nlp.Cfg).AKId)
		//	_ = currentEngineAppSecret.Set(currentEngine.GetCfg().(*huawei_cloud_nlp.Cfg).SkKey)
		//	_ = currentEngineProjectKey.Set(currentEngine.GetCfg().(*huawei_cloud_nlp.Cfg).ProjectId)
		//	entryAppKey.Enable()
		//	entryAppSecret.Enable()
		//	entryProjectKey.Enable()
		//case *niutrans.Cfg:
		//	_ = currentEngineAppKey.Set(currentEngine.GetCfg().(*niutrans.Cfg).AppKey)
		//	entryAppKey.Enable()
		//case *caiyunai.Cfg:
		//	_ = currentEngineAppSecret.Set(currentEngine.GetCfg().(*caiyunai.Cfg).Token)
		//	entryAppSecret.Enable()
		//case *openapi_youdao.Cfg:
		//	_ = currentEngineAppKey.Set(currentEngine.GetCfg().(*openapi_youdao.Cfg).AppKey)
		//	_ = currentEngineAppSecret.Set(currentEngine.GetCfg().(*openapi_youdao.Cfg).AppSecret)
		//	entryAppKey.Enable()
		//	entryAppSecret.Enable()
		//case *volcengine.Cfg:
		//	_ = currentEngineAppKey.Set(currentEngine.GetCfg().(*volcengine.Cfg).AccessKey)
		//	_ = currentEngineAppSecret.Set(currentEngine.GetCfg().(*volcengine.Cfg).SecretKey)
		//	entryAppKey.Enable()
		//	entryAppSecret.Enable()
		//case *baidu.Cfg:
		//	_ = currentEngineAppKey.Set(currentEngine.GetCfg().(*baidu.Cfg).AppId)
		//	_ = currentEngineAppSecret.Set(currentEngine.GetCfg().(*baidu.Cfg).AppKey)
		//	entryAppKey.Enable()
		//	entryAppSecret.Enable()
		//case *tencent_cloud_mt.Cfg:
		//	_ = currentEngineAppKey.Set(currentEngine.GetCfg().(*tencent_cloud_mt.Cfg).SecretId)
		//	_ = currentEngineAppSecret.Set(currentEngine.GetCfg().(*tencent_cloud_mt.Cfg).SecretKey)
		//	entryAppKey.Enable()
		//	entryAppSecret.Enable()
		//case *ali_cloud_mt.Cfg:
		//	_ = currentEngineAppKey.Set(currentEngine.GetCfg().(*ali_cloud_mt.Cfg).AKId)
		//	_ = currentEngineAppSecret.Set(currentEngine.GetCfg().(*ali_cloud_mt.Cfg).AKSecret)
		//	entryAppKey.Enable()
		//	entryAppSecret.Enable()
		//}
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
		parseInt, err := strconv.ParseInt(currentVal, 10, 32)
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
		parseInt, err := strconv.ParseInt(currentVal, 10, 32)
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
		parseInt, err := strconv.ParseInt(currentVal, 10, 32)
		if err != nil {
			return fmt.Errorf("无效数字: %s", currentVal)
		}
		if parseInt < 0 {
			return fmt.Errorf("协程数量不能小于0: %s", currentVal)
		}
		return nil
	}

	configForm.AppendItem(fiAppKey)
	configForm.AppendItem(fiAppSecret)
	configForm.AppendItem(fiProjectKey)

	configForm.AppendItem(&widget.FormItem{Text: "最大长度", Widget: entryMaxLength, HintText: "单次请求最大字符数量(含标点符号)"})
	configForm.AppendItem(&widget.FormItem{Text: "请求QPS", Widget: entryQPS, HintText: "每秒最大请求数量, 该值有一定误差, 建议小于标准值"})
	configForm.AppendItem(&widget.FormItem{Text: "协程数量", Widget: entryProcMax, HintText: "翻译并行程序数量(合理即可, 该值并不是越大越好)"})

	configForm.SubmitText = "保存"
	configForm.OnSubmit = func() {
		fmt.Println("提交")
	}

	return container.NewVBox(
		pageTitle(page.GetName()),
		configForm,
	)
}
