package ai_baidu

import (
	"anto/domain/service/translator"
	"anto/lib/log"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-module/carbon"
	"sync"
)

var api = "https://aip.baidubce.com/rpc/2.0/mt/texttrans/v1"

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
		id:            "ai_baidu",
		name:          "百度智能云",
		sep:           "\n",
		langSupported: langSupported,
	}
}

type Translator struct {
	id            string
	name          string
	cfg           translator.ImplConfig
	qps           int
	procMax       int
	textMaxLen    int
	langSupported []translator.LangPair
	sep           string
}

func (customT *Translator) Init(cfg translator.ImplConfig) { customT.cfg = cfg }

func (customT *Translator) GetId() string                           { return customT.id }
func (customT *Translator) GetShortId() string                      { return "ai_bd" }
func (customT *Translator) GetName() string                         { return customT.name }
func (customT *Translator) GetCfg() translator.ImplConfig           { return customT.cfg }
func (customT *Translator) GetLangSupported() []translator.LangPair { return customT.langSupported }
func (customT *Translator) GetSep() string                          { return customT.sep }
func (customT *Translator) IsValid() bool {
	return customT.cfg != nil && customT.cfg.GetAK() != "" && customT.cfg.GetSK() != ""
}

func (customT *Translator) Translate(ctx context.Context, args *translator.TranslateArgs) (*translator.TranslateRes, error) {
	timeStart := carbon.Now()
	accessToken, err := customT.GetCfg().(*Config).GetAccessToken()
	if err != nil {
		return nil, err
	}
	urlQueried := fmt.Sprintf("%s?access_token=%s", api, accessToken)
	bodyContent := translateReq{Q: args.TextContent, From: args.FromLang, To: args.ToLang}
	respBytes, err := translator.RequestSimpleHttp(ctx, customT, urlQueried, true, bodyContent, nil)
	if err != nil {
		return nil, err
	}
	respObj := new(translateResp)
	if err = json.Unmarshal(respBytes, respObj); err != nil {
		log.Singleton().ErrorF("解析报文异常, 引擎: %s, 错误: %s", customT.GetName(), err)
		return nil, fmt.Errorf("解析报文出现异常, 错误: %s", err.Error())
	}

	ret := new(translator.TranslateRes)
	for _, transBlockArray := range respObj.Result.TransResult {
		ret.Results = append(ret.Results, &translator.TranslateResBlock{
			Id:             transBlockArray.Src,
			TextTranslated: transBlockArray.Dst,
		})
	}

	ret.TimeUsed = int(carbon.Now().DiffInSeconds(timeStart))
	return ret, nil
}

type translateReq struct {
	Q    string `json:"q"`
	From string `json:"from"`
	To   string `json:"to"`
}

type translateResp struct {
	Result struct {
		From        string `json:"from"`
		TransResult []struct {
			Dst string `json:"dst"`
			Src string `json:"src"`
		} `json:"trans_result"`
		To string `json:"to"`
	} `json:"result"`
	LogId int64 `json:"log_id"`
}
