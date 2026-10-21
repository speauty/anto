package deepseek

import "anto/domain/service/translator"

var languages = []translator.LangPair{
	{Key: "zh", Name: "中文"}, {Key: "en", Name: "英语"},
	{Key: "ja", Name: "日语"}, {Key: "ko", Name: "韩语"},
	{Key: "fr", Name: "法语"}, {Key: "de", Name: "德语"},
	{Key: "es", Name: "西班牙语"}, {Key: "ru", Name: "俄语"},
	{Key: "it", Name: "意大利语"}, {Key: "pt", Name: "葡萄牙语"},
}
