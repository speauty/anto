package niutrans

import (
	"anto/domain/service/translator"
)

// https://niutrans.com/documents/contents/trans_text#languageList
var langSupported = []translator.LangPair{
	{"zh", "中文"},
	{"en", "英语"},
	{"ko", "韩语"},
	{"ja", "日语"},
	{"fr", "法语"},
	{"de", "德语"},
	{"it", "意大利语"},
}
