package openapi_youdao

import (
	"anto/domain/service/translator"
	"anto/lib/log"
	"anto/lib/util"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/golang-module/carbon"
	"github.com/google/go-querystring/query"
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
		id:            "openapi_youdao",
		name:          "有道智云",
		api:           "https://openapi.youdao.com/v2/api",
		qps:           1,
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
	api           string
	qps           int
	procMax       int
	textMaxLen    int
	langSupported []translator.LangPair
	sep           string
}

func (customT *Translator) Init(cfg interface{}) { customT.cfg = cfg.(*Cfg) }

func (customT *Translator) GetId() string       { return customT.id }
func (customT *Translator) GetShortId() string  { return "oy" }
func (customT *Translator) GetName() string     { return customT.name }
func (customT *Translator) GetCfg() interface{} { return nil }
func (customT *Translator) GetQPS() int         { return customT.qps }
func (customT *Translator) GetProcMax() int     { return customT.procMax }
func (customT *Translator) GetTextMaxLen() int {
	if customT.cfg.MaxSingleTextLength > 0 {
		return customT.cfg.MaxSingleTextLength
	}
	return customT.textMaxLen
}
func (customT *Translator) GetLangSupported() []translator.LangPair { return customT.langSupported }
func (customT *Translator) GetSep() string                          { return customT.sep }
func (customT *Translator) IsValid() bool {
	return customT.cfg != nil && customT.cfg.AppKey != "" && customT.cfg.AppSecret != ""
}

func (customT *Translator) Translate(ctx context.Context, args *translator.TranslateArgs) (*translator.TranslateRes, error) {
	timeStart := carbon.Now()
	texts := strings.Split(args.TextContent, customT.GetSep())
	newReq := &remoteReq{
		TextQuery: texts,
		From:      args.FromLang, To: args.ToLang,
		AppKey: customT.cfg.AppKey, Salt: util.Uid(), SignType: "v3",
		CurrentTime: fmt.Sprintf("%d", carbon.Now().Timestamp()),
	}
	newReq.Sign = customT.signBuilder(strings.Join(texts, ""), newReq.Salt, newReq.CurrentTime)
	params, _ := query.Values(newReq)
	urlQueried := fmt.Sprintf("%s?%s", customT.api, params.Encode())
	respBytes, err := translator.RequestSimpleGet(ctx, customT, urlQueried)
	if err != nil {
		return nil, err
	}
	newResp := new(remoteResp)
	if err = json.Unmarshal(respBytes, newResp); err != nil {
		log.Singleton().ErrorF("引擎: %s, 错误: 解析报文异常(%s)", customT.GetName(), err)
		return nil, fmt.Errorf("错误: 解析报文出现异常(%s)", err.Error())
	}

	if newResp.ErrorCode != "0" {
		errMsg := errorMap[newResp.ErrorCode]
		log.Singleton().ErrorF("引擎: %s, 错误: 接口响应异常(%s:%s)", customT.GetName(), newResp.ErrorCode, errMsg)
		return nil, fmt.Errorf("错误: 翻译异常(%s:%s)", newResp.ErrorCode, errMsg)
	}

	ret := new(translator.TranslateRes)
	for _, transBlock := range newResp.TranslateResults {
		ret.Results = append(ret.Results, &translator.TranslateResBlock{
			Id:             transBlock.Query,
			TextTranslated: transBlock.Translation,
		})
	}

	ret.TimeUsed = int(carbon.Now().DiffInSeconds(timeStart))
	return ret, nil
}

func (customT *Translator) signBuilder(textQuery, salt, currentTime string) string {
	tmpQuery := textQuery

	if tmpLen := len(textQuery); tmpLen > 20 {
		tmpQuery = fmt.Sprintf("%s%d%s", textQuery[0:10], tmpLen, textQuery[tmpLen-10:])
	}
	tmpQuery = fmt.Sprintf("%s%s%s%s%s", customT.cfg.AppKey, tmpQuery, salt, currentTime, customT.cfg.AppSecret)
	newSha := sha256.New()
	newSha.Write([]byte(tmpQuery))

	return hex.EncodeToString(newSha.Sum(nil))
}

type remoteReq struct {
	TextQuery   []string `url:"q"`        // 要翻译的文本.可指定多个
	From        string   `url:"from"`     // 源语言
	To          string   `url:"to"`       // 目标语言
	AppKey      string   `url:"appKey"`   // 应用标识（应用ID）
	Salt        string   `url:"salt"`     // 随机字符串，可使用UUID进行生产
	CurrentTime string   `url:"curtime"`  // 随机字符串，可使用UUID进行生产
	Sign        string   `url:"sign"`     // 签名信息sha256(appKey+q+salt+密钥)
	SignType    string   `url:"signType"` // 签名类型(v3)
}

type remoteResp struct {
	ErrorCode        string `json:"errorCode"`
	ErrorIndex       []int  `json:"errorIndex"`
	TranslateResults []struct {
		Query        string `json:"query"`
		Translation  string `json:"translation"`
		Type         string `json:"type"`
		VerifyResult string `json:"verifyResult"`
	} `json:"translateResults"`
}
