package caiyunai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang-module/carbon"
	"io"
	"net/http"
	"strings"
	"sync"
	"translator/tst/tt_log"
	"translator/tst/tt_translator"
	"translator/util"
)

var (
	apiTranslator  *Translator
	onceTranslator sync.Once
	api            = "http://api.interpreter.caiyunai.com/v1/translator"
)

func GetInstance() *Translator {
	onceTranslator.Do(func() {
		apiTranslator = New()
	})
	return apiTranslator
}

func New() *Translator {
	return &Translator{
		id:            "caiyun_ai",
		name:          "彩云小译",
		qps:           10,
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
	langSupported []tt_translator.LangK
	sep           string
}

func (customT *Translator) Init(cfg interface{}) { customT.cfg = cfg.(*Cfg) }

func (customT *Translator) GetId() string       { return customT.id }
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
func (customT *Translator) GetLangSupported() []tt_translator.LangK { return customT.langSupported }
func (customT *Translator) GetSep() string                          { return customT.sep }
func (customT *Translator) IsValid() bool                           { return customT.cfg != nil && customT.cfg.Token != "" }

func (customT *Translator) Translate(args *tt_translator.TranslateArgs) (*tt_translator.TranslateRes, error) {
	timeStart := carbon.Now()
	texts := strings.Split(args.TextContent, customT.GetSep())
	req := new(translateReq)
	req.Source = texts
	req.RequestId = util.Uid()
	req.TransType = fmt.Sprintf("%s2%s", args.FromLang, args.ToLang)

	reqBytes, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest(http.MethodPost, api, bytes.NewReader(reqBytes))
	httpReq.Header.Set("content-type", "application/json")
	httpReq.Header.Set("x-authorization", fmt.Sprintf("token %s", customT.cfg.Token))
	httpResp, err := new(http.Client).Do(httpReq)
	defer func() {
		if httpResp != nil && httpResp.Body != nil {
			_ = httpResp.Body.Close()
		}
	}()
	if err != nil {
		tt_log.GetInstance().Error(fmt.Sprintf("调用接口失败, 引擎: %s, 错误: %s", customT.GetName(), err))
		return nil, fmt.Errorf("网络请求出现异常, 错误: %s", err.Error())
	}
	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		tt_log.GetInstance().Error(fmt.Sprintf("读取报文异常, 引擎: %s, 错误: %s", customT.GetName(), err))
		return nil, fmt.Errorf("读取报文出现异常, 错误: %s", err.Error())
	}

	resp := new(translateResp)
	if err = json.Unmarshal(respBytes, resp); err != nil {
		tt_log.GetInstance().Error(fmt.Sprintf("解析报文异常, 引擎: %s, 错误: %s", customT.GetName(), err))
		return nil, fmt.Errorf("解析报文出现异常, 错误: %s", err.Error())
	}

	if resp.Msg != "" {
		tt_log.GetInstance().Error(fmt.Sprintf("接口响应异常, 引擎: %s, 错误: %s", customT.GetName(), resp.Msg))
		return nil, fmt.Errorf("接口响应异常, 引擎: %s, 错误: %s", customT.GetName(), resp.Msg)
	}

	if len(texts) != len(resp.Target) {
		tt_log.GetInstance().Error(fmt.Sprintf("响应解析错误, 引擎: %s, 错误: 译文和原文数量匹配失败", customT.GetName()))
		return nil, fmt.Errorf("翻译异常, 错误: 源文和译文数量不对等")
	}

	ret := new(tt_translator.TranslateRes)

	for textIdx, textTarget := range resp.Target {
		ret.Results = append(ret.Results, &tt_translator.TranslateResBlock{
			Id:             texts[textIdx],
			TextTranslated: textTarget,
		})
	}

	ret.TimeUsed = int(carbon.Now().DiffAbsInSeconds(timeStart))
	return ret, nil

}

type translateReq struct {
	Source    []string `json:"source"`
	TransType string   `json:"trans_type"`
	RequestId string   `json:"request_id"`
}

type translateResp struct {
	Msg        string   `json:"message,omitempty"`
	SrcTgt     []string `json:"src_tgt,omitempty"`
	Target     []string `json:"target,omitempty"`
	Confidence float64  `json:"confidence,omitempty"` // 可信度?
	Rc         int      `json:"rc,omitempty"`
}
