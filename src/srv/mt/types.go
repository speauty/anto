package mt

import "context"

type TextTranslateResp struct {
	Idx           string // 文本序号
	StrTranslated string // 翻译文本
}

type MT interface {
	GetId() Id
	GetName() string
	Init(context.Context, interface{}) error
	TextTranslate(context.Context, interface{}) (*TextTranslateResp, error)
	TextBatchTranslate(context.Context, interface{}) ([]TextTranslateResp, error)
}

type Id string

const (
	ALI   Id = "ALi"
	BAIDU Id = "BaiDu"
)
