package caiyunai

import "anto/domain/service/translator"

var langSupported = []translator.LangPair{
	{Key: "zh", Name: "中文"},
	{Key: "zh-Hant", Name: "繁体中文"},
	{Key: "en", Name: "英语"},
	{Key: "ja", Name: "日语"},
	{Key: "ko", Name: "韩语"},
	{Key: "de", Name: "德语"},
	{Key: "es", Name: "西班牙语"},
	{Key: "fr", Name: "法语"},
	{Key: "it", Name: "意大利语"},
	{Key: "pt", Name: "葡萄牙语"},
	{Key: "ru", Name: "俄语"},
	{Key: "tr", Name: "土耳其语"},
	{Key: "vi", Name: "越南语"},
}

// SupportsDirection follows the translation matrix in Caiyun's API docs.
func SupportsDirection(from, to string) bool {
	if from == to {
		return false
	}
	switch from {
	case "zh", "zh-Hant":
		return to == "zh" || to == "zh-Hant" || to == "en" || to == "ja" || to == "ko"
	case "en", "ja", "ko", "de", "es", "fr", "it", "pt", "ru", "tr", "vi":
		return to == "zh" || to == "zh-Hant"
	default:
		return false
	}
}
