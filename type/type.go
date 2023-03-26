package _type

const (
	ModeFull  TranslateMode = "全量翻译"
	ModeDelta TranslateMode = "增量翻译"
)

var modes = []TranslateMode{ModeFull, ModeDelta}

type TranslateMode string

func (customType TranslateMode) String() string {
	return string(customType)
}

func (customType TranslateMode) GetModes() []string { // ComboBox不支持typedef么?
	var strModes []string
	for _, mode := range modes {
		strModes = append(strModes, mode.String())
	}
	return strModes
}

func (customType TranslateMode) GetIdx() int {
	for idx, mode := range modes {
		if customType == mode {
			return idx
		}
	}
	return 0
}

const (
	LangDirectionFrom LangDirection = "来源语种"
	LangDirectionTo   LangDirection = "目标语种"
)

var langDirectionTypes = []LangDirection{LangDirectionFrom, LangDirectionTo}

type LangDirection string

func (customType LangDirection) String() string {
	return string(customType)
}

func (customType LangDirection) GetDirections() []string { // ComboBox不支持typedef么?
	var strLangDirections []string
	for _, langDirection := range langDirectionTypes {
		strLangDirections = append(strLangDirections, langDirection.String())
	}
	return strLangDirections
}

func (customType LangDirection) GetIdx() int {
	for idx, langDirection := range langDirectionTypes {
		if customType == langDirection {
			return idx
		}
	}
	return 0
}
