package huawei_cloud_nlp

import (
	"anto/tst/tt_translator"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/nlp/v2/model"
)

var langSupported = []tt_translator.LangK{
	{"zh", "中文(简体)"},
	{"en", "英语"},
	{"ru", "俄语"},
	{"fr", "法语"},
	{"de", "德语"},
	{"ko", "韩语"},
	{"ja", "日语"},
	{"th", "泰语"},
	{"ar", "阿拉伯语"},
	{"pt", "葡萄牙语"},
	{"tr", "土耳其语"},
	{"es", "西班牙语"},
	{"vi", "越南语"},
}

var langTo = map[string]model.TextTranslationReqTo{
	"zh": model.GetTextTranslationReqToEnum().ZH,
	"en": model.GetTextTranslationReqToEnum().EN,
	"ru": model.GetTextTranslationReqToEnum().RU,
	"fr": model.GetTextTranslationReqToEnum().FR,
	"de": model.GetTextTranslationReqToEnum().DE,
	"ko": model.GetTextTranslationReqToEnum().KO,
	"ja": model.GetTextTranslationReqToEnum().JA,
	"th": model.GetTextTranslationReqToEnum().TH,
	"ar": model.GetTextTranslationReqToEnum().AR,
	"pt": model.GetTextTranslationReqToEnum().PT,
	"tr": model.GetTextTranslationReqToEnum().TR,
	"es": model.GetTextTranslationReqToEnum().ES,
	"vi": model.GetTextTranslationReqToEnum().VI,
}

var langFrom = map[string]model.TextTranslationReqFrom{
	"zh": model.GetTextTranslationReqFromEnum().ZH,
	"en": model.GetTextTranslationReqFromEnum().EN,
	"ru": model.GetTextTranslationReqFromEnum().RU,
	"fr": model.GetTextTranslationReqFromEnum().FR,
	"de": model.GetTextTranslationReqFromEnum().DE,
	"ko": model.GetTextTranslationReqFromEnum().KO,
	"ja": model.GetTextTranslationReqFromEnum().JA,
	"th": model.GetTextTranslationReqFromEnum().TH,
	"ar": model.GetTextTranslationReqFromEnum().AR,
	"pt": model.GetTextTranslationReqFromEnum().PT,
	"tr": model.GetTextTranslationReqFromEnum().TR,
	"es": model.GetTextTranslationReqFromEnum().ES,
	"vi": model.GetTextTranslationReqFromEnum().VI,
}
