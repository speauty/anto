package tencent_cloud_mt

import "anto/domain/service/translator"

// Tencent's classic text API uses jp/kr rather than the ja/ko codes used by
// some other providers.
var langSupported = []translator.LangPair{
	{Key: "zh", Name: "中文"},
	{Key: "zh-TW", Name: "繁体中文"},
	{Key: "en", Name: "英语"},
	{Key: "jp", Name: "日语"},
	{Key: "kr", Name: "韩语"},
	{Key: "de", Name: "德语"},
	{Key: "fr", Name: "法语"},
	{Key: "es", Name: "西班牙语"},
	{Key: "it", Name: "意大利语"},
	{Key: "tr", Name: "土耳其语"},
	{Key: "ru", Name: "俄语"},
	{Key: "pt", Name: "葡萄牙语"},
	{Key: "vi", Name: "越南语"},
	{Key: "id", Name: "印度尼西亚语"},
	{Key: "ms", Name: "马来语"},
	{Key: "th", Name: "泰语"},
}
