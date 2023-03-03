package mt

import "context"

type TextTranslateResp struct {
	Idx           string // 文本序号
	StrTranslated string // 翻译文本
}

type MT interface {
	GetId() Id
	GetName() string
	GetCfg() interface{}
	Init(context.Context, interface{}) error
	TextTranslate(context.Context, interface{}) ([]TextTranslateResp, error)
	TextBatchTranslate(context.Context, interface{}) ([]TextTranslateResp, error)
}

type Id string

const (
	IdALiYun  Id = "EngineALi"
	IdBaiDu   Id = "EngineBaiDu"
	IdYouDao  Id = "EngineYouDao"
	IdTencent Id = "EngineTencent"
)

type Engine int

// FromInt 从int转化
func (engine Engine) FromInt(num int) Engine {
	return Engine(num)
}

// ToInt 转化成int输出
func (engine Engine) ToInt() int {
	return int(engine)
}

func (engine Engine) GetZHArrays() []string {
	return engineZHMaps
}

func (engine Engine) GetZH() string {
	return engineZHMaps[engine]
}

const (
	// EngineALiYun 阿里云机器翻译
	EngineALiYun Engine = iota
	// EngineBaiDu 百度翻译
	EngineBaiDu
	// EngineYouDao 有道翻译
	EngineYouDao
	// EngineTencent 腾讯翻译
	EngineTencent
)

var engineZHMaps = []string{"阿里云", "百度", "有道", "腾讯"}

const BlockSep string = "\n"
const BlockIdxContentSep string = "@<"

const MaxCoroutine int = 100
