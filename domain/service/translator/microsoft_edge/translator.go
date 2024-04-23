package microsoft_edge

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

func API() *Translator {
	onceTranslator.Do(func() {
		apiTranslator = New()
	})
	return apiTranslator
}

func New() *Translator {
	return &Translator{
		id:            "edge",
		name:          "微软翻译",
		sep:           "\n",
		langSupported: langSupported,
	}
}

// Translator LingVA翻译已崩溃, 当前处于不可用状态, 所以直接禁用
type Translator struct {
	id            string
	name          string
	cfg           translator.ImplConfig
	langSupported []translator.LangPair
	sep           string
	token         string
}

func (customT *Translator) Init(cfg translator.ImplConfig) { customT.cfg = cfg }

func (customT *Translator) GetId() string                           { return customT.id }
func (customT *Translator) GetShortId() string                      { return "me" }
func (customT *Translator) GetName() string                         { return customT.name }
func (customT *Translator) GetCfg() translator.ImplConfig           { return customT.cfg }
func (customT *Translator) GetLangSupported() []translator.LangPair { return customT.langSupported }
func (customT *Translator) GetSep() string                          { return customT.sep }
func (customT *Translator) IsValid() bool                           { return true }

func (customT *Translator) Translate(ctx context.Context, args *translator.TranslateArgs) (*translator.TranslateRes, error) {
	timeStart := carbon.Now()
	token, err := customT.getToken(ctx)
	if err != nil {
		return nil, err
	}

	var queryUrl = fmt.Sprintf(
		"https://api-edge.cognitive.microsofttranslator.com/translate?from=%s&to=%s&api-version=3.0&includeSentenceLength=true",
		args.FromLang, args.ToLang,
	)
	req := []microsoftEdgeReq{}
	textReq := strings.Split(args.TextContent, customT.sep)
	for _, text := range textReq {
		req = append(req, microsoftEdgeReq{Text: text})
	}
	respBytes, err := translator.RequestSimpleHttp(ctx, customT, queryUrl, true, req, map[string]string{
		"Authorization": fmt.Sprintf("Bearer Bearer %s", token),
	})
	if err != nil {
		return nil, err
	}
	resp := []microsoftEdgeResp{}
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		log.Singleton().ErrorF("解析报文异常, 引擎: %s, 错误: %s", customT.GetName(), err)
		return nil, fmt.Errorf("解析报文出现异常, 错误: %s", err.Error())
	}
	textResp := []string{}
	for _, edgeResp := range resp {
		if edgeResp.Error.Code != 0 {
			// The request is not authorized because credentials are missing or invalid.
			if edgeResp.Error.Code == 401001 {
				customT.token = ""
			}
			return nil, fmt.Errorf("翻译异常: %s", edgeResp.Error.Message)
		}
		for _, translation := range edgeResp.Translations {
			textResp = append(textResp, translation.Text)
		}
	}
	if len(textReq) != len(textResp) {
		return nil, translator.ErrSrcAndTgtNotMatched
	}

	ret := new(translator.TranslateRes)
	for textIdx, textSource := range textReq {
		ret.Results = append(ret.Results, &translator.TranslateResBlock{
			Id:             textSource,
			TextTranslated: textResp[textIdx],
		})
	}

	ret.TimeUsed = int(carbon.Now().DiffAbsInSeconds(timeStart))
	return ret, nil

}

func (customT *Translator) getToken(ctx context.Context) (token string, err error) {
	if customT.token == "" {
		tokenBytes := []byte{}
		tokenBytes, err = translator.RequestSimpleHttp(ctx, customT, "https://edge.microsoft.com/translate/auth", false, nil, nil)
		if err != nil {
			return
		}
		customT.token = string(tokenBytes)
	}
	token = customT.token
	return
}

type microsoftEdgeReq struct {
	Text string `json:"text"`
}

type microsoftEdgeResp struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	DetectedLanguage struct {
		Language string  `json:"language"`
		Score    float64 `json:"score"`
	} `json:"detectedLanguage"`
	Translations []struct {
		Text    string `json:"text"`
		To      string `json:"to"`
		SentLen struct {
			SrcSentLen   []int `json:"srcSentLen"`
			TransSentLen []int `json:"transSentLen"`
		} `json:"sentLen"`
	} `json:"translations"`
}
