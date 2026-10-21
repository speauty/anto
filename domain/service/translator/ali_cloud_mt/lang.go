package ali_cloud_mt

import (
	"anto/domain/service/translator"
)

var langSupported = []translator.LangPair{
	{Key: "zh", Name: "中文"},
	{Key: "en", Name: "英语"},
	{Key: "ja", Name: "日语"},
	{Key: "ko", Name: "韩语"},
	{Key: "ru", Name: "俄语"},
	{Key: "fr", Name: "法语"},
	{Key: "de", Name: "德语"},
	{Key: "ar", Name: "阿拉伯语"},
}
