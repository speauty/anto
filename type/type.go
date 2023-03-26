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
