package deepl_pro

import (
	"anto/domain/service/translator"
	"strings"
)

// This stable, common subset works for text translation on the Pro endpoint.
// Language IDs stay lowercase for consistency with the other ANTO engines.
var langSupported = []translator.LangPair{
	{Key: "zh", Name: "中文"}, {Key: "en", Name: "英语"},
	{Key: "de", Name: "德语"}, {Key: "es", Name: "西班牙语"},
	{Key: "fr", Name: "法语"}, {Key: "it", Name: "意大利语"},
	{Key: "ja", Name: "日语"}, {Key: "ko", Name: "韩语"},
	{Key: "nl", Name: "荷兰语"}, {Key: "pl", Name: "波兰语"},
	{Key: "pt", Name: "葡萄牙语"}, {Key: "ru", Name: "俄语"},
	{Key: "bg", Name: "保加利亚语"}, {Key: "cs", Name: "捷克语"},
	{Key: "da", Name: "丹麦语"}, {Key: "el", Name: "希腊语"},
	{Key: "et", Name: "爱沙尼亚语"}, {Key: "fi", Name: "芬兰语"},
	{Key: "hu", Name: "匈牙利语"}, {Key: "id", Name: "印度尼西亚语"},
	{Key: "lt", Name: "立陶宛语"}, {Key: "lv", Name: "拉脱维亚语"},
	{Key: "nb", Name: "挪威语"}, {Key: "ro", Name: "罗马尼亚语"},
	{Key: "sk", Name: "斯洛伐克语"}, {Key: "sl", Name: "斯洛文尼亚语"},
	{Key: "sv", Name: "瑞典语"}, {Key: "tr", Name: "土耳其语"},
	{Key: "uk", Name: "乌克兰语"},
}

func languageCode(value string, target bool) (string, bool) {
	for _, language := range langSupported {
		if language.Key != value {
			continue
		}
		if target {
			switch value {
			case "en":
				return "EN-US", true
			case "pt":
				return "PT-PT", true
			}
		}
		return strings.ToUpper(value), true
	}
	return "", false
}

func SupportsDirection(from, to string) bool {
	if from == to {
		return false
	}
	_, sourceOK := languageCode(from, false)
	_, targetOK := languageCode(to, true)
	return sourceOK && targetOK
}
