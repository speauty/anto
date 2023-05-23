package ling_va

import (
	"anto/domain/service/translator"
	"anto/lib/log"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-module/carbon"
	"go.uber.org/zap"
	"net/url"
	"strings"
	"sync"
)

var (
	apiTranslator  *Translator
	onceTranslator sync.Once
)

func Singleton() *Translator {
	onceTranslator.Do(func() {
		apiTranslator = New()
	})
	return apiTranslator
}

func New() *Translator {
	return &Translator{
		id:            "lingva",
		name:          "Lingva",
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
func (customT *Translator) GetShortId() string                      { return "lv" }
func (customT *Translator) GetName() string                         { return customT.name }
func (customT *Translator) GetCfg() translator.ImplConfig           { return customT.cfg }
func (customT *Translator) GetLangSupported() []translator.LangPair { return customT.langSupported }
func (customT *Translator) GetSep() string                          { return customT.sep }
func (customT *Translator) IsValid() bool                           { return customT.cfg.GetAK() != "" }

func (customT *Translator) Translate(ctx context.Context, args *translator.TranslateArgs) (*translator.TranslateRes, error) {
	timeStart := carbon.Now()

	var api = fmt.Sprintf("https://lingva.ml/_next/data/%s/", customT.cfg.GetAK())
	queryUrl := fmt.Sprintf(
		"%s/%s/%s/%s.json", api,
		args.FromLang, args.ToLang, url.PathEscape(args.TextContent),
	)
	respBytes, err := translator.RequestSimpleGet(ctx, customT, queryUrl)
	if err != nil {
		return nil, err
	}
	lingVaResp := new(lingVaMTResp)
	if err := json.Unmarshal(respBytes, lingVaResp); err != nil {
		fmt.Println(string(respBytes))
		log.Singleton().ErrorF("解析报文异常, 引擎: %s, 错误: %s", customT.GetName(), err)
		return nil, fmt.Errorf("解析报文出现异常, 错误: %s", err.Error())
	}
	if lingVaResp.State == false {
		log.Singleton().ErrorF("接口响应异常, 引擎: %s, 错误: %s", customT.GetName(), err, zap.String("result", string(respBytes)))
		return nil, fmt.Errorf("翻译异常")
	}
	textTranslatedList := strings.Split(lingVaResp.Props.TextTranslated, customT.sep)
	textSourceList := strings.Split(lingVaResp.Props.Params.TextSource, customT.sep)
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
