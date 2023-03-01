package mt

type TextTranslateResp struct {
	Idx           string // 文本序号
	StrTranslated string // 翻译文本
}

type MT interface {
	GetName() string
	Init(cfg interface{}) error
	TextBatchTranslate(args interface{}) ([]TextTranslateResp, error)
}
