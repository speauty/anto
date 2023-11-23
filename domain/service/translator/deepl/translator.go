package deepl

import (
	"anto/domain/service/translator"
	"anto/lib/log"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-module/carbon"
	"strings"
	"sync"
)

var (
	apiTranslator  *Translator
	onceTranslator sync.Once
)

const (
	TRANSLATE_API_FREE string = "https://api-free.deepl.com/v2/translate"
	TRANSLATE_API_PRO  string = "https://api.deepl.com/v2/translate"
)

func API() *Translator {
	onceTranslator.Do(func() {
		apiTranslator = New()
	})
	return apiTranslator
}

func New() *Translator {
	return &Translator{
		id:            "deepl",
		name:          "DeepL",
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
func (customT *Translator) GetShortId() string                      { return "dl" }
func (customT *Translator) GetName() string                         { return customT.name }
func (customT *Translator) GetCfg() translator.ImplConfig           { return customT.cfg }
func (customT *Translator) GetLangSupported() []translator.LangPair { return customT.langSupported }
func (customT *Translator) GetSep() string                          { return customT.sep }
func (customT *Translator) IsValid() bool                           { return true }

func (customT *Translator) Translate(ctx context.Context, args *translator.TranslateArgs) (*translator.TranslateRes, error) {
	timeStart := carbon.Now()
	texts := strings.Split(args.TextContent, customT.GetSep())
	req := &translateReq{
		Text: texts, SourceLang: args.FromLang, TargetLang: args.ToLang, SplitSentences: "0",
	}

	reqBytes, _ := json.Marshal(req)
	respBytes, err := translator.RequestSimpleHttp(ctx, customT, TRANSLATE_API_FREE, true, reqBytes, map[string]string{
		"Authorization": fmt.Sprintf("DeepL-Auth-Key %s", customT.cfg.GetAK()),
	})
	if err != nil {
		return nil, err
	}

	resp := new(translateResp)
	if err = json.Unmarshal(respBytes, resp); err != nil {
		log.Singleton().ErrorF("解析报文异常, 引擎: %s, 错误: %s", customT.GetName(), err)
		return nil, fmt.Errorf("解析报文出现异常, 错误: %s", err.Error())
	}

	ret := new(translator.TranslateRes)

	for textIdx, textTarget := range resp.Translations {
		ret.Results = append(ret.Results, &translator.TranslateResBlock{
			Id:             texts[textIdx],
			TextTranslated: textTarget.Text,
		})
	}

	ret.TimeUsed = int(carbon.Now().DiffAbsInSeconds(timeStart))
	return ret, nil
}

type translateReq struct {
	Text           []string `json:"text"`
	SourceLang     string   `json:"source_lang"`
	TargetLang     string   `json:"target_lang"`
	SplitSentences string   `json:"split_sentences"`
}

type translateResp struct {
	Translations []struct {
		DetectedSourceLanguage string `json:"detected_source_language"`
		Text                   string `json:"text"`
	} `json:"translations"`
}
