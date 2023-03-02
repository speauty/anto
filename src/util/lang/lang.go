package lang

import "sort"

type StrLang string

const (
	ZH StrLang = "zh"
	EN StrLang = "en"
)

func (sl StrLang) ToString() string {
	return string(sl)
}

func (sl StrLang) GetCH() string {
	return maps[sl]
}

func (sl StrLang) GetMaps() []string {
	var chMaps, keys []string
	for key, _ := range maps {
		keys = append(keys, key.ToString())
	}
	sort.Strings(keys)
	for _, key := range keys {
		chMaps = append(chMaps, maps[StrLang(key)])
	}
	return chMaps
}

func (sl StrLang) GetLangByIdx(idx int) StrLang {
	chMaps := sl.GetMaps()
	for lang, val := range maps {
		if val == chMaps[idx] {
			return lang
		}
	}
	return ""
}

var (
	maps = map[StrLang]string{
		ZH: "中文",
		EN: "英文",
	}
)
