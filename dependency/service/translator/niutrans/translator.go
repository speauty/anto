package niutrans

import (
	"anto/dependency/service/translator"
	"anto/lib/log"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-module/carbon"
	"strings"
	"sync"
)

const apiTranslate = "https://api.niutrans.com/NiuTransServer/translation"

var (
	apiSingleton  *Translator
	onceSingleton sync.Once
)

func Singleton() *Translator {
	onceSingleton.Do(func() {
		apiSingleton = New()
	})
	return apiSingleton
}

func New() *Translator {
	return &Translator{
		id:            "niutrans",
		name:          "小牛翻译",
		qps:           50,
		procMax:       20,
		textMaxLen:    5000,
		sep:           "\n",
		langSupported: langSupported,
	}
}

type Translator struct {
	id            string
	name          string
	cfg           *Cfg
	qps           int
	procMax       int
	textMaxLen    int
	langSupported []translator.LangPair
	sep           string
}

func (customT *Translator) Init(cfg interface{}) { customT.cfg = cfg.(*Cfg) }

func (customT *Translator) GetId() string                           { return customT.id }
func (customT *Translator) GetShortId() string                      { return "nt" }
func (customT *Translator) GetName() string                         { return customT.name }
func (customT *Translator) GetCfg() interface{}                     { return nil }
func (customT *Translator) GetQPS() int                             { return customT.qps }
func (customT *Translator) GetProcMax() int                         { return customT.procMax }
func (customT *Translator) GetTextMaxLen() int                      { return customT.textMaxLen }
func (customT *Translator) GetLangSupported() []translator.LangPair { return customT.langSupported }
func (customT *Translator) GetSep() string                          { return customT.sep }
func (customT *Translator) IsValid() bool                           { return customT.cfg != nil && customT.cfg.AppKey != "" }

func (customT *Translator) Translate(ctx context.Context, args *translator.TranslateArgs) (*translator.TranslateRes, error) {
	timeStart := carbon.Now()
	tr := &translateRequest{
		Apikey:  customT.cfg.AppKey,
		SrcText: args.TextContent,
		From:    args.FromLang,
		To:      args.ToLang,
	}
	respBytes, err := translator.RequestSimplePost(ctx, customT, apiTranslate, tr)
	if err != nil {
		return nil, err
	}

	resp := new(translateResponse)
	if err = json.Unmarshal(respBytes, resp); err != nil {
		log.Singleton().ErrorF("解析报文异常, 引擎: %s, 错误: %s", customT.GetName(), err)
		return nil, fmt.Errorf("解析报文出现异常, 错误: %s", err.Error())
	}

	if resp.ErrorMsg != "" {
		log.Singleton().ErrorF("接口响应异常, 引擎: %s, 错误: %s(%s)", customT.GetName(), resp.ErrorMsg, resp.ErrorCode)
		return nil, fmt.Errorf("接口响应异常, 引擎: %s, 错误: %s", customT.GetName(), resp.ErrorMsg)
	}

	srcTexts := strings.Split(args.TextContent, customT.GetSep())
	tgtTexts := strings.Split(resp.TgtText, customT.GetSep())
	if len(srcTexts) != len(tgtTexts) {
		return nil, translator.ErrSrcAndTgtNotMatched
	}

	ret := new(translator.TranslateRes)

	for textIdx, textTarget := range tgtTexts {
		ret.Results = append(ret.Results, &translator.TranslateResBlock{
			Id: srcTexts[textIdx], TextTranslated: textTarget,
		})
	}

	ret.TimeUsed = int(carbon.Now().DiffAbsInSeconds(timeStart))
	return ret, nil

}

type translateRequest struct {
	Apikey  string `json:"apikey"`
	SrcText string `json:"src_text"`
	From    string `json:"from"`
	To      string `json:"to"`
}

type translateResponse struct {
	TgtText   string `json:"tgt_text,omitempty"`
	To        string `json:"to"`
	From      string `json:"from"`
	ErrorCode string `json:"error_code,omitempty"`
	ErrorMsg  string `json:"error_msg,omitempty"`
}
