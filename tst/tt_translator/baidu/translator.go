package baidu

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/golang-module/carbon"
	"net/url"
	"sync"
	"translator/tst/tt_log"
	"translator/tst/tt_translator"
	"translator/util"
)

var api = "https://fanyi-api.baidu.com/api/trans/vip/translate"

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
		id:            "baidu",
		name:          "百度翻译",
		qps:           1,
		procMax:       1,
		textMaxLen:    1000,
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
	langSupported []tt_translator.LangK
	sep           string
}

func (customT *Translator) Init(cfg interface{}) { customT.cfg = cfg.(*Cfg) }

func (customT *Translator) GetId() string                           { return customT.id }
func (customT *Translator) GetName() string                         { return customT.name }
func (customT *Translator) GetCfg() interface{}                     { return nil }
func (customT *Translator) GetQPS() int                             { return customT.qps }
func (customT *Translator) GetProcMax() int                         { return customT.procMax }
func (customT *Translator) GetTextMaxLen() int                      { return customT.textMaxLen }
func (customT *Translator) GetLangSupported() []tt_translator.LangK { return customT.langSupported }
func (customT *Translator) GetSep() string                          { return customT.sep }
func (customT *Translator) IsValid() bool {
	return customT.cfg != nil && customT.cfg.AppId != "" && customT.cfg.AppKey != ""
}

func (customT *Translator) Translate(args *tt_translator.TranslateArgs) (*tt_translator.TranslateRes, error) {
	timeStart := carbon.Now()
	salt := util.Uid()
	sign := customT.signBuilder(args.TextContent, salt)
	urlQueried := fmt.Sprintf(
		"%s?q=%s&from=%s&to=%s&appid=%s&salt=%s&sign=%s", api,
		url.QueryEscape(args.TextContent), args.FromLang, args.ToLang,
		customT.cfg.AppId, salt, sign,
	)
	respBytes, err := tt_translator.RequestSimpleGet(customT, urlQueried)
	if err != nil {
		return nil, err
	}
	respObj := new(remoteResp)
	if err := json.Unmarshal(respBytes, respObj); err != nil {
		tt_log.GetInstance().Error(fmt.Sprintf("解析报文异常, 引擎: %s, 错误: %s", customT.GetName(), err))
		return nil, fmt.Errorf("解析报文出现异常, 错误: %s", err.Error())
	}
	if respObj.ErrorCode != "" && respObj.ErrorCode != "52000" {
		tt_log.GetInstance().Error(fmt.Sprintf("接口响应异常, 引擎: %s, 代码: %s, 错误: %s", customT.GetName(), respObj.ErrorCode, respObj.ErrorMsg))
		return nil, fmt.Errorf("翻译异常, 代码: %s, 错误: %s", respObj.ErrorCode, respObj.ErrorMsg)
	}

	ret := new(tt_translator.TranslateRes)
	for _, transBlockArray := range respObj.Results {
		ret.Results = append(ret.Results, &tt_translator.TranslateResBlock{
			Id:             transBlockArray.Src,
			TextTranslated: transBlockArray.Dst,
		})
	}

	ret.TimeUsed = int(carbon.Now().DiffInSeconds(timeStart))
	return ret, nil
}

func (customT *Translator) signBuilder(strQuery string, salt string) string {
	tmpStr := fmt.Sprintf("%s%s%s%s", customT.cfg.AppId, strQuery, salt, customT.cfg.AppKey)
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
