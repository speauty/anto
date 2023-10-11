package google_cloud

import (
	"anto/domain/service/translator"
	"cloud.google.com/go/translate"
	"context"
	"errors"
	"fmt"
	"github.com/golang-module/carbon"
	"google.golang.org/api/option"
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
		id:            "google_cloud",
		name:          "谷歌云",
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
func (customT *Translator) GetShortId() string                      { return "gc" }
func (customT *Translator) GetName() string                         { return customT.name }
func (customT *Translator) GetCfg() translator.ImplConfig           { return customT.cfg }
func (customT *Translator) GetLangSupported() []translator.LangPair { return customT.langSupported }
func (customT *Translator) GetSep() string                          { return customT.sep }
func (customT *Translator) IsValid() bool                           { return customT.cfg.GetAK() != "" }

func (customT *Translator) Translate(ctx context.Context, args *translator.TranslateArgs) (*translator.TranslateRes, error) {
	timeStart := carbon.Now()
	client, err := translate.NewClient(ctx, option.WithAPIKey(customT.cfg.GetAK()))
	if err != nil {
		return nil, err
	}
	fromLangTag := convertLangToTag(args.FromLang)
	toLangTag := convertLangToTag(args.ToLang)
	if fromLangTag == nil || toLangTag == nil {
		return nil, errors.New(fmt.Sprintf("来源语种或目标语种暂不支持"))
	}
	textRaw := strings.Split(args.TextContent, customT.GetSep())
	results, err := client.Translate(ctx, textRaw, *fromLangTag, &translate.Options{
		Source: *toLangTag, Format: translate.Text,
	})
	if err != nil {
		return nil, err
	}
	if len(results) != len(textRaw) {
		return nil, translator.ErrSrcAndTgtNotMatched
	}

	ret := new(translator.TranslateRes)
	for textIdx, textSource := range textRaw {
		ret.Results = append(ret.Results, &translator.TranslateResBlock{Id: textSource, TextTranslated: results[textIdx].Text})
	}

	ret.TimeUsed = int(carbon.Now().DiffAbsInSeconds(timeStart))
	return ret, nil

}

type lingVaMTResp struct {
	State bool `json:"__N_SSG"`
	Props struct {
		Type           int    `json:"type"`
		TextTranslated string `json:"translation"`
		Params         struct {
			FromLanguage string `json:"source"`
			ToLanguage   string `json:"target"`
			TextSource   string `json:"query"`
		} `json:"initial"`
	} `json:"pageProps"`
}
