package translator

type InterfaceTranslator interface {
	Init(cfg interface{})
	GetId() string
	GetName() string
	GetCfg() interface{}
	GetQPS() int
	GetProcMax() int
	GetTextMaxLen() int
	GetLangSupported() []LangPair
	GetSep() string
	IsValid() bool
	Translate(*TranslateArgs) (*TranslateRes, error)
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
