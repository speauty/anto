package youdao_batch

import "anto/domain/service/translator"

var langSupported = []translator.LangPair{
	{Key: "zh", Name: "中文"}, {Key: "zh-Hant", Name: "繁体中文"},
	{Key: "en", Name: "英语"}, {Key: "ja", Name: "日语"},
	{Key: "ko", Name: "韩语"}, {Key: "fr", Name: "法语"},
	{Key: "es", Name: "西班牙语"}, {Key: "pt", Name: "葡萄牙语"},
	{Key: "it", Name: "意大利语"}, {Key: "ru", Name: "俄语"},
	{Key: "vi", Name: "越南语"}, {Key: "de", Name: "德语"},
	{Key: "ar", Name: "阿拉伯语"}, {Key: "id", Name: "印尼语"},
	{Key: "th", Name: "泰语"}, {Key: "tr", Name: "土耳其语"},
	{Key: "uk", Name: "乌克兰语"}, {Key: "hi", Name: "印地语"},
	{Key: "nl", Name: "荷兰语"}, {Key: "pl", Name: "波兰语"},
}

func languageCode(code string) (string, bool) {
	switch code {
	case "zh":
		return "zh-CHS", true
	case "zh-Hant":
		return "zh-CHT", true
	}
	for _, lang := range langSupported {
		if lang.Key == code {
			return code, true
		}
	}
	return "", false
}

func SupportsDirection(from, to string) bool {
	if from == to {
		return false
	}
	_, sourceOK := languageCode(from)
	_, targetOK := languageCode(to)
	return sourceOK && targetOK
}
