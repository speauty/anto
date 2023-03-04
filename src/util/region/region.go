package region

import (
	"gui.subtitle/src/srv/mt"
)

type Region int

func (r Region) FromInt(num int) Region {
	return Region(num)
}

func (r Region) ToInt() int {
	return int(r)
}

func (r Region) GetZh(id mt.Id) string {
	switch id {
	case mt.IdTencent:
		return tencentZHMaps[r]
	default:
		return ""
	}
}

func (r Region) GetZhMaps(id mt.Id) []string {
	switch id {
	case mt.IdTencent:
		return tencentZHMaps
	default:
		return []string{}
	}
}

func (r Region) GetEn(id mt.Id) string {
	switch id {
	case mt.IdTencent:
		return tencentENMaps[r]
	default:
		return ""
	}
}
