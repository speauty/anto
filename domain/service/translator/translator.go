package translator

import (
	"context"
	"errors"
)

var ErrInvalidConnectionRegion = errors.New("翻译引擎地区无效")

type ConnectionResult struct {
	Endpoint  string `json:"endpoint"`
	LatencyMs int64  `json:"latency_ms"`
}

type ConnectionTester interface {
	TestConnection(context.Context, string) (*ConnectionResult, error)
}

type ImplTranslator interface {
	Init(cfg ImplConfig)
	GetId() string
	GetShortId() string
	GetName() string
	GetCfg() ImplConfig
	GetLangSupported() []LangPair
	GetSep() string
	IsValid() bool
	Translate(context.Context, *TranslateArgs) (*TranslateRes, error)
}

// BatchTranslator preserves the position of each subtitle block, including
// repeated or multiline text, when a provider accepts several texts at once.
type BatchTranslator interface {
	TranslateBatch(context.Context, string, string, []string) ([]string, error)
}

type TranslateArgs struct {
	FromLang    string
	ToLang      string
	TextContent string
}

type TranslateRes struct {
	TimeUsed int
	Msg      []string
	Results  []*TranslateResBlock
}

type TranslateResBlock struct {
	Id             string // 采用原文做ID
	TextTranslated string
}

type LangPair struct {
	Key  string // 语种编码
	Name string // 语种名称
}
