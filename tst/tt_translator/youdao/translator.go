package youdao

import (
	"anto/tst/tt_log"
	"anto/tst/tt_translator"
	"encoding/json"
	"fmt"
	"github.com/golang-module/carbon"
	"go.uber.org/zap"
	"net/url"
	"strings"
	"sync"
)

var api = "https://fanyi.youdao.com/translate?&doctype=json"

var (
	apiTranslator  *Translator
	onceTranslator sync.Once
)

func GetInstance() *Translator {
	onceTranslator.Do(func() {
		apiTranslator = New()
	})
	return apiTranslator
}

func New() *Translator {
	return &Translator{
		id:            "youdao",
		name:          "有道翻译",
		qps:           50,
		procMax:       20,
		textMaxLen:    2000,
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
	langSupported []tt_translator.LangK
	sep           string
}

func (customT *Translator) Init(_ interface{}) {}

func (customT *Translator) GetId() string                           { return customT.id }
func (customT *Translator) GetName() string                         { return customT.name }
func (customT *Translator) GetCfg() interface{}                     { return nil }
func (customT *Translator) GetQPS() int                             { return customT.qps }
func (customT *Translator) GetProcMax() int                         { return customT.procMax }
func (customT *Translator) GetTextMaxLen() int                      { return customT.textMaxLen }
func (customT *Translator) GetLangSupported() []tt_translator.LangK { return customT.langSupported }
func (customT *Translator) GetSep() string                          { return customT.sep }
func (customT *Translator) IsValid() bool                           { return true }

func (customT *Translator) Translate(args *tt_translator.TranslateArgs) (*tt_translator.TranslateRes, error) {
	timeStart := carbon.Now()
	urlQueried := fmt.Sprintf(
		"%s&type=%s2%s&i=%s", api,
		strings.ToUpper(args.FromLang), strings.ToUpper(args.ToLang),
		url.QueryEscape(args.TextContent),
	)

	respBytes, err := tt_translator.RequestSimpleGet(customT, urlQueried)
	if err != nil {
		return nil, err
	}
	youDaoResp := new(youDaoMTResp)
	if err := json.Unmarshal(respBytes, youDaoResp); err != nil {
		tt_log.GetInstance().Error(fmt.Sprintf("解析报文异常, 引擎: %s, 错误: %s", customT.GetName(), err))
		return nil, fmt.Errorf("解析报文出现异常, 错误: %s", err.Error())
	}
	if youDaoResp.ErrorCode != 0 {
		tt_log.GetInstance().Error(fmt.Sprintf("接口响应异常, 引擎: %s, 错误: %s", customT.GetName(), err), zap.String("result", string(respBytes)))
		return nil, fmt.Errorf("翻译异常, 代码: %d", youDaoResp.ErrorCode)
	}

	ret := new(tt_translator.TranslateRes)
	for _, transBlockArray := range youDaoResp.TransResult {
		for _, block := range transBlockArray {
			ret.Results = append(ret.Results, &tt_translator.TranslateResBlock{
				Id:             block.Src,
				TextTranslated: block.Tgt,
			})
		}
	}

	ret.TimeUsed = int(carbon.Now().DiffInSeconds(timeStart))
	return ret, nil
}

type youDaoMTResp struct {
	Type        string `json:"type"`
	ErrorCode   int    `json:"errorCode"`
	ElapsedTime int    `json:"elapsedTime"`
	TransResult [][]struct {
		Src string `json:"src,omitempty"` // 原文
		Tgt string `json:"tgt,omitempty"` // 译文
	} `json:"translateResult"`
}
