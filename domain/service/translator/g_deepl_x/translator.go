package g_deepl_x

import (
	"anto/domain/service/translator"
	"context"
	"github.com/OwO-Network/gdeeplx"
	"github.com/golang-module/carbon"
	"strings"
	"sync"
)

var (
	apiTranslator  *Translator
	onceTranslator sync.Once
)

func API() *Translator {
	onceTranslator.Do(func() {
		apiTranslator = New()
	})
	return apiTranslator
}

func New() *Translator {
	return &Translator{
		id:            "g_deepl_x",
		name:          "GDeeplX",
		sep:           "\n",
		langSupported: langSupported,
	}
}

type Translator struct {
	id            string
	name          string
	cfg           translator.ImplConfig
	langSupported []translator.LangPair
	sep           string
}

func (customT *Translator) Init(cfg translator.ImplConfig) { customT.cfg = cfg }

func (customT *Translator) GetId() string                           { return customT.id }
func (customT *Translator) GetShortId() string                      { return "gdx" }
func (customT *Translator) GetName() string                         { return customT.name }
func (customT *Translator) GetCfg() translator.ImplConfig           { return customT.cfg }
func (customT *Translator) GetLangSupported() []translator.LangPair { return customT.langSupported }
func (customT *Translator) GetSep() string                          { return customT.sep }
func (customT *Translator) IsValid() bool                           { return true }

func (customT *Translator) Translate(ctx context.Context, args *translator.TranslateArgs) (*translator.TranslateRes, error) {
	timeStart := carbon.Now()
	result, err := gdeeplx.Translate(args.TextContent, args.FromLang, args.ToLang, 0)
	if err != nil {
		return nil, err
	}
	resultParsed := result.(map[string]interface{})
	textTranslatedList := strings.Split(resultParsed["data"].(string), customT.sep)
	textSourceList := strings.Split(args.TextContent, customT.sep)
	if len(textSourceList) != len(textTranslatedList) {
		return nil, translator.ErrSrcAndTgtNotMatched
	}

	ret := new(translator.TranslateRes)
	for textIdx, textSource := range textSourceList {
		ret.Results = append(ret.Results, &translator.TranslateResBlock{
			Id:             textSource,
			TextTranslated: textTranslatedList[textIdx],
		})
	}

	ret.TimeUsed = int(carbon.Now().DiffAbsInSeconds(timeStart))
	return ret, nil

}
