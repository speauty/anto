package xfyun

import (
	"anto/domain/service/translator"
	"anto/lib/log"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/golang-module/carbon"
	"go.uber.org/zap"
	"net/url"
	"strings"
	"sync"
)

const apiTranslate string = "https://itrans.xf-yun.com/v1/its"
const host string = "itrans.xf-yun.com"
const algorithm string = "hmac-sha256"
const headers string = "host date request-line"

var (
	apiSingleton  *Translator
	onceSingleton sync.Once
)

func API() *Translator {
	onceSingleton.Do(func() {
		apiSingleton = New()
	})
	return apiSingleton
}

func New() *Translator {
	return &Translator{
		id:            "xfyun",
		name:          "科大讯飞",
		sep:           "",
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
func (customT *Translator) GetShortId() string                      { return "xy" }
func (customT *Translator) GetName() string                         { return customT.name }
func (customT *Translator) GetCfg() translator.ImplConfig           { return customT.cfg }
func (customT *Translator) GetLangSupported() []translator.LangPair { return customT.langSupported }
func (customT *Translator) GetSep() string                          { return customT.sep }
func (customT *Translator) IsValid() bool {
	return customT.cfg != nil && customT.cfg.GetAK() != "" && customT.cfg.GetSK() != "" && customT.cfg.GetProjectKey() != ""
}

func (customT *Translator) signature(date string, requestLine string) string {
	sigOriginal := fmt.Sprintf("host: %s\ndate: %s\n%s", host, date, requestLine)
	mac := hmac.New(sha256.New, []byte(customT.GetCfg().GetSK()))
	mac.Write([]byte(sigOriginal))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func (customT *Translator) authorization(sig string) string {
	authorizationOriginal := fmt.Sprintf(
		"api_key=\"%s\", algorithm=\"%s\", headers=\"%s\", signature=\"%s\"",
		customT.cfg.GetAK(), algorithm, headers, sig,
	)
	return base64.StdEncoding.EncodeToString([]byte(authorizationOriginal))
}

func (customT *Translator) Translate(ctx context.Context, args *translator.TranslateArgs) (*translator.TranslateRes, error) {
	timeStart := carbon.Now()
	tr := &translateRequest{
		Header:    reqHeader{AppId: customT.GetCfg().GetProjectKey(), Status: 3},
		Parameter: reqParameter{reqParameterIts{From: args.FromLang, To: args.ToLang, Result: struct{}{}}},
		Payload:   reqPayload{InputData: reqPayloadInputData{Encoding: "utf8", Status: 3, Text: base64.StdEncoding.EncodeToString([]byte(args.TextContent))}},
	}
	sigTimeStr := carbon.Now().ToRfc1123String("GMT")
	sig := customT.signature(sigTimeStr, "POST /v1/its HTTP/1.1")
	authorization := customT.authorization(sig)
	reqBytes, _ := json.Marshal(tr)
	queryParams := fmt.Sprintf("host=%s&date=%s&authorization=%s", host, sigTimeStr, authorization)
	respBytes, err := translator.RequestSimpleHttp(ctx, customT, fmt.Sprintf("%s?%s", apiTranslate, url.PathEscape(queryParams)), true, reqBytes, nil)
	if err != nil {
		return nil, err
	}

	resp := new(translateResponse)
	if err = json.Unmarshal(respBytes, resp); err != nil {
		log.Singleton().ErrorF("解析报文异常, 引擎: %s, 错误: %s", customT.GetName(), err)
		return nil, fmt.Errorf("解析报文出现异常, 错误: %s", err.Error())
	}

	if resp.Header.Code != 0 {
		log.Singleton().ErrorF("接口响应异常, 引擎: %s, 错误: %s", customT.GetName(), err, zap.String("result", string(respBytes)))
		return nil, fmt.Errorf("翻译异常, 代码: %d", resp.Header.Message)
	}
	textTranslatedBytes, err := base64.StdEncoding.DecodeString(resp.Payload.Result.Text)
	if err != nil {
		log.Singleton().ErrorF("解析译文异常(base64), 引擎: %s, 错误: %s", customT.GetName(), err)
		return nil, fmt.Errorf("解析译文出现异常, 错误: %s", err.Error())
	}
	textTranslated := new(RespText)
	if err = json.Unmarshal(textTranslatedBytes, textTranslated); err != nil {
		log.Singleton().ErrorF("解析译文异常(json), 引擎: %s, 错误: %s", customT.GetName(), err)
		return nil, fmt.Errorf("解析译文出现异常, 错误: %s", err.Error())
	}

	ret := new(translator.TranslateRes)
	if customT.GetSep() == "" {
		ret.Results = append(ret.Results, &translator.TranslateResBlock{
			Id: textTranslated.TransResult.Src, TextTranslated: textTranslated.TransResult.Dst,
		})
	} else {
		srcTexts := strings.Split(textTranslated.TransResult.Src, customT.GetSep())
		tgtTexts := strings.Split(textTranslated.TransResult.Dst, customT.GetSep())
		if len(srcTexts) != len(tgtTexts) {
			return nil, translator.ErrSrcAndTgtNotMatched
		}

		for textIdx, textTarget := range tgtTexts {
			ret.Results = append(ret.Results, &translator.TranslateResBlock{
				Id: srcTexts[textIdx], TextTranslated: textTarget,
			})
		}
	}
	ret.TimeUsed = int(carbon.Now().DiffAbsInSeconds(timeStart))
	return ret, nil

}

type reqHeader struct {
	AppId  string `json:"app_id"` // 在平台申请的appid信息
	Status int    `json:"status"` // 请求状态，固定取值为：3（一次传完）
	ResId  string `json:"res_id"`
}

type reqParameter struct {
	Its reqParameterIts `json:"its"` // 用于上传功能参数
}

type reqParameterIts struct {
	From   string   `json:"from"`   // 源语种
	To     string   `json:"to"`     // 目标语种
	Result struct{} `json:"result"` // 响应结果字段名，固定传{}即可，为保留字段
}

type reqPayload struct {
	InputData reqPayloadInputData `json:"input_data"` // 用于上传相关数据
}

type reqPayloadInputData struct {
	Encoding string `json:"encoding"` // 文本编码，如utf8
	Status   int    `json:"status"`   // 数据状态，固定取值为：3（一次传完）
	Text     string `json:"text"`     // 待翻译文本的base64编码，字符要大于0且小于5000
}

type translateRequest struct {
	Header    reqHeader    `json:"header"`
	Parameter reqParameter `json:"parameter"`
	Payload   reqPayload   `json:"payload"`
}

type translateResponse struct {
	Header struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Sid     string `json:"sid"`
	} `json:"header"`
	Payload struct {
		Result struct {
			Seq    string `json:"seq"`
			Status string `json:"status"`
			Text   string `json:"text"`
		} `json:"result"`
	} `json:"payload"`
}

type RespText struct {
	TransResult struct {
		Dst string `json:"dst"`
		Src string `json:"src"`
	} `json:"trans_result"`
	From string `json:"from"`
	To   string `json:"to"`
}
