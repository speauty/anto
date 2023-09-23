package baidu

import (
	"anto/domain/service/translator"
	"anto/lib/log"
	"anto/lib/util"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/golang-module/carbon"
	"net/url"
	"sync"
)

var api = "https://fanyi-api.baidu.com/api/trans/vip/translate"

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
		id:            "baidu",
		name:          "百度翻译",
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
func (customT *Translator) GetShortId() string                      { return "bd" }
func (customT *Translator) GetName() string                         { return customT.name }
func (customT *Translator) GetCfg() translator.ImplConfig           { return customT.cfg }
func (customT *Translator) GetLangSupported() []translator.LangPair { return customT.langSupported }
func (customT *Translator) GetSep() string                          { return customT.sep }
func (customT *Translator) IsValid() bool {
	return customT.cfg != nil && customT.cfg.GetAK() != "" && customT.cfg.GetSK() != ""
}

func (customT *Translator) Translate(ctx context.Context, args *translator.TranslateArgs) (*translator.TranslateRes, error) {
	timeStart := carbon.Now()
	salt := util.Uid()
	sign := customT.signBuilder(args.TextContent, salt)
	urlQueried := fmt.Sprintf(
		"%s?q=%s&from=%s&to=%s&appid=%s&salt=%s&sign=%s", api,
		url.QueryEscape(args.TextContent), args.FromLang, args.ToLang,
		customT.cfg.GetAK(), salt, sign,
	)
	respBytes, err := translator.RequestSimpleHttp(ctx, customT, urlQueried, false, nil, nil)
	if err != nil {
		return nil, err
	}
	respObj := new(remoteResp)
	if err = json.Unmarshal(respBytes, respObj); err != nil {
		log.Singleton().ErrorF("解析报文异常, 引擎: %s, 错误: %s", customT.GetName(), err)
		return nil, fmt.Errorf("解析报文出现异常, 错误: %s", err.Error())
	}
	if respObj.ErrorCode != "" && respObj.ErrorCode != "52000" {
		log.Singleton().ErrorF("接口响应异常, 引擎: %s, 错误: %s(%s)", customT.GetName(), respObj.ErrorCode, respObj.ErrorMsg)
		return nil, fmt.Errorf("翻译异常, 代码: %s, 错误: %s", respObj.ErrorCode, respObj.ErrorMsg)
	}

	ret := new(translator.TranslateRes)
	for _, transBlockArray := range respObj.Results {
		ret.Results = append(ret.Results, &translator.TranslateResBlock{
			Id:             transBlockArray.Src,
			TextTranslated: transBlockArray.Dst,
		})
	}

	ret.TimeUsed = int(carbon.Now().DiffInSeconds(timeStart))
	return ret, nil
}

func (customT *Translator) signBuilder(strQuery string, salt string) string {
	tmpStr := fmt.Sprintf("%s%s%s%s", customT.cfg.GetAK(), strQuery, salt, customT.cfg.GetSK())
	tmpMD5 := md5.New()
	tmpMD5.Write([]byte(tmpStr))
	return fmt.Sprintf("%x", tmpMD5.Sum(nil))
}

type remoteResp struct {
	ErrorCode string `json:"error_code,omitempty"`
	ErrorMsg  string `json:"error_msg,omitempty"`
	From      string `json:"from"`
	To        string `json:"to"`
	Results   []struct {
		Src string `json:"src"`
		Dst string `json:"dst"`
	} `json:"trans_result"`
}
