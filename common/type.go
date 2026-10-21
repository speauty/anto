package common

const (
	ModeFull  TranslateMode = "全量翻译"
	ModeDelta TranslateMode = "增量翻译"
)

type TranslateMode string

func (customType TranslateMode) String() string {
	return string(customType)
}

const (
	LangDirectionFrom LangDirection = "来源语种"
	LangDirectionTo   LangDirection = "目标语种"
)

type LangDirection string

func (customType LangDirection) String() string {
	return string(customType)
}
