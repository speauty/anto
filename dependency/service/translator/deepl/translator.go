package deepl

import (
	"anto/dependency/service/translator"
	"anto/lib/log"
	"encoding/json"
	"fmt"
	"github.com/golang-module/carbon"
	"strings"
	"sync"
)

var (
	apiTranslator  *Translator
	onceTranslator sync.Once
	api            = "https://www2.deepl.com/jsonrpc?method=LMT_handle_jobs"
)

func Singleton() *Translator {
	onceTranslator.Do(func() {
		apiTranslator = New()
	})
	return apiTranslator
}

func New() *Translator {
	return &Translator{
		id:            "deepl",
		name:          "DeepL",
		qps:           10,
		procMax:       20,
		textMaxLen:    1000,
		sep:           "\n",
		langSupported: langSupported,
	}
}

type Translator struct {
	id            string
	name          string
	qps           int
	procMax       int
	textMaxLen    int
	langSupported []translator.LangPair
	sep           string
}

func (customT *Translator) Init(_ interface{}) {}

func (customT *Translator) GetId() string                           { return customT.id }
func (customT *Translator) GetName() string                         { return customT.name }
func (customT *Translator) GetCfg() interface{}                     { return nil }
func (customT *Translator) GetQPS() int                             { return customT.qps }
func (customT *Translator) GetProcMax() int                         { return customT.procMax }
func (customT *Translator) GetTextMaxLen() int                      { return customT.textMaxLen }
func (customT *Translator) GetLangSupported() []translator.LangPair { return customT.langSupported }
func (customT *Translator) GetSep() string                          { return customT.sep }
func (customT *Translator) IsValid() bool                           { return false }

func (customT *Translator) Translate(args *translator.TranslateArgs) (*translator.TranslateRes, error) {
	timeStart := carbon.Now()
	texts := strings.Split(args.TextContent, customT.GetSep())

	req := new(translateReq).New()
	req.Params.Lang.SourceLangComputed = args.FromLang
	req.Params.Lang.TargetLang = args.ToLang

	for textIdx, text := range texts {
		req.Params.Jobs[0].Sentences = append(req.Params.Jobs[0].Sentences, &translateReqParamsSentence{
			Text: text,
			Id:   textIdx,
		})
	}

	respBytes, err := translator.RequestSimplePost(customT, api, req)
	if err != nil {
		return nil, err
	}

	resp := new(translateResp)
	if err = json.Unmarshal(respBytes, resp); err != nil {
		log.Singleton().ErrorF("解析报文异常, 引擎: %s, 错误: %s", customT.GetName(), err)
		return nil, fmt.Errorf("解析报文出现异常, 错误: %s", err.Error())
	}

	if resp.Error.Message != "" {
		log.Singleton().ErrorF("响应解析错误, 引擎: %s, 错误: %s", customT.GetName(), resp.Error.Message)
		return nil, fmt.Errorf("翻译异常, 错误: %s", resp.Error.Message)
	}

	if len(texts) != len(resp.Result.Translations[0].Beams[0].Sentences) {
		return nil, translator.ErrSrcAndTgtNotMatched
	}

	ret := new(translator.TranslateRes)
	for _, textSource := range resp.Result.Translations[0].Beams[0].Sentences {
		ret.Results = append(ret.Results, &translator.TranslateResBlock{
			Id:             texts[textSource.Ids[0]],
			TextTranslated: textSource.Text,
		})
	}

	ret.TimeUsed = int(carbon.Now().DiffAbsInSeconds(timeStart))
	return ret, nil

}

type translateReq struct {
	JsonRPC string              `json:"jsonrpc"` //
	Method  string              `json:"method"`
	Params  *translateReqParams `json:"params"`
}

type translateReqParams struct {
	Jobs      []*translateReqParamsJob `json:"jobs"`
	Lang      *translateReqParamsLang  `json:"lang"`
	Timestamp int64                    `json:"timestamp"`
}

type translateReqParamsJob struct {
	Kind              string                        `json:"kind"`
	Sentences         []*translateReqParamsSentence `json:"sentences"`
	PreferredNumBeams int                           `json:"preferred_num_beams"`
}

type translateReqParamsSentence struct {
	Text   string `json:"text"`
	Id     int    `json:"id"`
	Prefix string `json:"prefix"`
}

type translateReqParamsLang struct {
	SourceLangComputed string `json:"source_lang_computed"`
	TargetLang         string `json:"target_lang"`
}

func (customTR *translateReq) New() *translateReq {
	defaultJob := &translateReqParamsJob{Kind: "default", PreferredNumBeams: 0}
	req := &translateReq{
		JsonRPC: "2.0",
		Method:  "LMT_handle_jobs",
		Params: &translateReqParams{
			Lang:      &translateReqParamsLang{},
			Timestamp: carbon.Now().Timestamp(),
		},
	}
	req.Params.Jobs = append(req.Params.Jobs, defaultJob)
	return req
}

type translateResp struct {
	JsonRPC string `json:"jsonrpc"`
	Error   struct {
		Code    int    `json:"code,omitempty"`
		Message string `json:"message,omitempty"`
	} `json:"error,omitempty"`
	Result struct {
		Translations []struct {
			Beams []struct {
				Sentences []struct {
					Text string `json:"text"`
					Ids  []int  `json:"ids"`
				} `json:"sentences"`
				NumSymbols int `json:"num_symbols"`
			} `json:"beams"`
			Quality string `json:"quality"`
		} `json:"translations"`
		TargetLang            string   `json:"target_lang"`
		SourceLang            string   `json:"source_lang"`
		SourceLangIsConfident bool     `json:"source_lang_is_confident"`
		DetectedLanguages     struct{} `json:"detectedLanguages"`
	} `json:"result"`
}
