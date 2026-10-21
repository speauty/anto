package core

import (
	"anto/cfg"
	"anto/domain/repository"
	"anto/domain/service/translator/baidu_cloud_mt"
	"anto/domain/service/translator/deepseek"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

const redactedSecret = "********"

// ConfigValidationError marks user supplied values that fail range or shape
// checks, allowing the HTTP layer to return a client error instead of a
// generic persistence failure.
type ConfigValidationError struct{ Message string }

func (e *ConfigValidationError) Error() string { return e.Message }

// AppConfig contains client preferences. Provider credentials are kept in the
// Providers section so the Flutter client does not depend on Viper's map shape.
type AppConfig struct {
	DefaultTranslator string `json:"default_translator"`
	DefaultFromLang   string `json:"default_from_lang"`
	DefaultToLang     string `json:"default_to_lang"`
	OutputDirectory   string `json:"output_directory"`
	TrackExport       int    `json:"track_export"`
	KeepOriginal      bool   `json:"keep_original"`
	NotifyOnComplete  bool   `json:"notify_on_complete"`
	AutoConnect       bool   `json:"auto_connect"`
	MinimizeToTray    bool   `json:"minimize_to_tray"`
	LogLevel          string `json:"log_level"`
	Concurrency       int    `json:"concurrency"`
}

type ProviderConfig struct {
	Configured bool                   `json:"configured"`
	Values     map[string]interface{} `json:"values"`
}

type ConfigDTO struct {
	App       AppConfig                 `json:"app"`
	Providers map[string]ProviderConfig `json:"providers"`
}

// ConfigUpdate is deliberately a stable envelope. The values inside one
// provider remain extensible because every vendor has different credentials.
type ConfigUpdate struct {
	App       json.RawMessage                   `json:"app,omitempty"`
	Providers map[string]map[string]interface{} `json:"providers,omitempty"`
	// UI is accepted for one release so older Flutter prototypes can migrate
	// without silently losing their settings.
	UI map[string]interface{} `json:"ui,omitempty"`
}

func appConfigFromSettings(values map[string]interface{}) AppConfig {
	app := AppConfig{TrackExport: 2, KeepOriginal: true, NotifyOnComplete: true, AutoConnect: true, LogLevel: "信息", Concurrency: 5}
	ui, _ := values["ui"].(map[string]interface{})
	app.DefaultTranslator = stringValue(ui["default_translator"])
	app.DefaultFromLang = stringValue(ui["default_from_lang"])
	app.DefaultToLang = stringValue(ui["default_to_lang"])
	app.OutputDirectory = stringValue(ui["output_directory"])
	app.TrackExport = intValue(ui["track_export"], app.TrackExport)
	app.KeepOriginal = boolValue(ui["keep_original"], app.KeepOriginal)
	app.NotifyOnComplete = boolValue(ui["notify_on_complete"], app.NotifyOnComplete)
	app.AutoConnect = boolValue(ui["auto_connect"], app.AutoConnect)
	app.MinimizeToTray = boolValue(ui["minimize_to_tray"], app.MinimizeToTray)
	app.LogLevel = stringValueWithDefault(ui["log_level"], app.LogLevel)
	app.Concurrency = intValue(ui["concurrency"], app.Concurrency)
	return app
}

func appConfigValues(app AppConfig) map[string]interface{} {
	return map[string]interface{}{
		"default_translator": app.DefaultTranslator,
		"default_from_lang":  app.DefaultFromLang,
		"default_to_lang":    app.DefaultToLang,
		"output_directory":   app.OutputDirectory,
		"track_export":       app.TrackExport,
		"keep_original":      app.KeepOriginal,
		"notify_on_complete": app.NotifyOnComplete,
		"auto_connect":       app.AutoConnect,
		"minimize_to_tray":   app.MinimizeToTray,
		"log_level":          app.LogLevel,
		"concurrency":        app.Concurrency,
	}
}

func buildConfigDTO() ConfigDTO {
	settings := cfg.Singleton().GetViper().AllSettings()
	dto := ConfigDTO{App: appConfigFromSettings(settings), Providers: make(map[string]ProviderConfig)}
	for key, raw := range settings {
		if key == "ui" || repository.GetTranslators().GetById(key) == nil {
			continue
		}
		values, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		dto.Providers[key] = ProviderConfig{Values: redactValues(values), Configured: providerConfigured(key)}
	}
	for _, item := range repository.GetTranslators().GetNamesAll() {
		provider := dto.Providers[item.Key]
		if provider.Values == nil {
			provider.Values = make(map[string]interface{})
		}
		if engine := repository.GetTranslators().GetById(item.Key); engine != nil {
			for key, value := range configValues(engine.GetCfg()) {
				if _, exists := provider.Values[key]; !exists {
					provider.Values[key] = value
				}
			}
		}
		provider.Values = redactValues(provider.Values)
		provider.Configured = providerConfigured(item.Key)
		dto.Providers[item.Key] = provider
	}
	return dto
}

func configValues(config interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	value := reflect.ValueOf(config)
	if !value.IsValid() {
		return result
	}
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return result
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return result
	}
	typeOf := value.Type()
	for index := 0; index < value.NumField(); index++ {
		field := typeOf.Field(index)
		tag := strings.Split(field.Tag.Get("mapstructure"), ",")[0]
		if tag == "" || tag == "-" {
			continue
		}
		// These values are runtime token cache fields, not user editable
		// provider settings and must never be sent back by the client.
		if tag == "access_token" || tag == "expired_at" {
			continue
		}
		current := value.Field(index)
		if current.CanInterface() && (current.Kind() == reflect.String || current.Kind() >= reflect.Int && current.Kind() <= reflect.Int64) {
			result[tag] = current.Interface()
		}
	}
	return result
}

func validateAppConfig(app AppConfig) error {
	if app.DefaultTranslator != "" && repository.GetTranslators().GetById(app.DefaultTranslator) == nil {
		return &ConfigValidationError{Message: fmt.Sprintf("默认翻译引擎不存在: %s", app.DefaultTranslator)}
	}
	if app.TrackExport < 0 || app.TrackExport > 2 {
		return &ConfigValidationError{Message: "track_export 必须是 0、1 或 2"}
	}
	if app.Concurrency < 1 || app.Concurrency > 256 {
		return &ConfigValidationError{Message: "concurrency 必须在 1 到 256 之间"}
	}
	return nil
}

func validateProviderValues(provider string, values map[string]interface{}) error {
	if provider == "deepseek" {
		if raw, ok := values["base_url"]; ok {
			base, valid := raw.(string)
			if !valid {
				return &ConfigValidationError{Message: "DeepSeek 服务地址无效"}
			}
			if err := deepseek.ValidateBaseURL(base); err != nil {
				return &ConfigValidationError{Message: err.Error()}
			}
		}
		if raw, ok := values["model"]; ok {
			model, valid := raw.(string)
			if !valid || strings.TrimSpace(model) == "" || len(model) > 128 {
				return &ConfigValidationError{Message: "DeepSeek 模型名无效"}
			}
		}
		if raw, ok := values["prompt"]; ok {
			prompt, valid := raw.(string)
			if !valid || len(prompt) > 8000 {
				return &ConfigValidationError{Message: "DeepSeek 提示词不能超过 8000 字节"}
			}
		}
	}
	if provider == "baidu_cloud_mt" {
		if _, exists := values["qps"]; exists {
			return &ConfigValidationError{Message: "百度云 QPS 由服务版本决定"}
		}
		if _, exists := values["max_single_text_length"]; exists {
			return &ConfigValidationError{Message: "百度云单次最大长度由服务版本决定"}
		}
		if raw, exists := values["service_tier"]; exists {
			tier, ok := raw.(string)
			if !ok {
				return &ConfigValidationError{Message: "百度云服务版本无效"}
			}
			if _, _, valid := baidu_cloud_mt.TierLimits(tier); !valid {
				return &ConfigValidationError{Message: "百度云服务版本无效"}
			}
		}
	}
	for key, value := range values {
		if key != "qps" && key != "max_single_text_length" && key != "max_coroutine_num" {
			continue
		}
		number, ok := numberValue(value)
		if !ok || number < 1 {
			return &ConfigValidationError{Message: fmt.Sprintf("翻译引擎 %s 的 %s 必须是大于 0 的整数", provider, key)}
		}
		if key == "max_single_text_length" && number > 1000000 {
			return &ConfigValidationError{Message: fmt.Sprintf("翻译引擎 %s 的 %s 不能超过 1000000", provider, key)}
		}
		if provider == "deepseek" && key == "max_single_text_length" && number > 6000 {
			return &ConfigValidationError{Message: "DeepSeek 单次最大长度不能超过 6000 个字符"}
		}
		if provider == "tencent_cloud_mt" && key == "max_single_text_length" && number > 2000 {
			return &ConfigValidationError{Message: "腾讯云单次最大长度不能超过 2000 个字符"}
		}
		if provider == "youdao_batch" && key == "max_single_text_length" && number > 5000 {
			return &ConfigValidationError{Message: "有道智云单次最大长度不能超过 5000 个字符"}
		}
		if key == "qps" && number > 100000 {
			return &ConfigValidationError{Message: fmt.Sprintf("翻译引擎 %s 的 qps 不能超过 100000", provider)}
		}
		if key == "max_coroutine_num" && number > 1024 {
			return &ConfigValidationError{Message: fmt.Sprintf("翻译引擎 %s 的 max_coroutine_num 不能超过 1024", provider)}
		}
	}
	return nil
}

func numberValue(value interface{}) (int, bool) {
	switch result := value.(type) {
	case int:
		return result, true
	case int8:
		return int(result), true
	case int16:
		return int(result), true
	case int32:
		return int(result), true
	case int64:
		return int(result), true
	case uint:
		return int(result), true
	case uint8:
		return int(result), true
	case uint16:
		return int(result), true
	case uint32:
		return int(result), true
	case uint64:
		if uint64(int(result)) != result {
			return 0, false
		}
		return int(result), true
	case float64:
		if result != float64(int(result)) {
			return 0, false
		}
		return int(result), true
	case float32:
		if result != float32(int(result)) {
			return 0, false
		}
		return int(result), true
	default:
		return 0, false
	}
}

func providerConfigured(id string) bool {
	engine := repository.GetTranslators().GetById(id)
	return engine != nil && engine.IsValid()
}

func redactValues(values map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(values))
	for key, value := range values {
		if key == "access_token" || key == "expired_at" || key == "configured" || key == "values" {
			continue
		}
		if isSecretKey(key) && stringValue(value) != "" {
			result[key] = redactedSecret
		} else {
			result[key] = value
		}
	}
	return result
}

func isSecretKey(key string) bool {
	key = strings.ToLower(key)
	if key == "project_key" {
		return false
	}
	for _, marker := range []string{"key", "secret", "token", "password"} {
		if strings.Contains(key, marker) {
			return true
		}
	}
	return false
}

func stringValue(value interface{}) string {
	if value == nil {
		return ""
	}
	if result, ok := value.(string); ok {
		return result
	}
	return ""
}

func stringValueWithDefault(value interface{}, fallback string) string {
	if result := stringValue(value); result != "" {
		return result
	}
	return fallback
}

func intValue(value interface{}, fallback int) int {
	switch result := value.(type) {
	case int:
		return result
	case int64:
		return int(result)
	case float64:
		return int(result)
	}
	return fallback
}

func boolValue(value interface{}, fallback bool) bool {
	if result, ok := value.(bool); ok {
		return result
	}
	return fallback
}
