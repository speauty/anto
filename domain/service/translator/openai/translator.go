package openai

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

const apiTranslate = "https://api.openai.com/v1/chat/completions"

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
		id:            "openai",
		name:          "OpenAI",
		sep:           "\n",
		langSupported: langSupported,
	}
}

type Translator struct {
	id              string
	name            string
	cfg             translator.ImplConfig
	langSupported   []translator.LangPair
	sep             string
	currentFromLang string
	currentToLang   string
	currentMessages *translateRequest
}

func (customT *Translator) Init(cfg translator.ImplConfig) { customT.cfg = cfg }

func (customT *Translator) GetId() string                           { return customT.id }
func (customT *Translator) GetShortId() string                      { return "nt" }
func (customT *Translator) GetName() string                         { return customT.name }
func (customT *Translator) GetCfg() translator.ImplConfig           { return customT.cfg }
func (customT *Translator) GetLangSupported() []translator.LangPair { return customT.langSupported }
func (customT *Translator) GetSep() string                          { return customT.sep }
func (customT *Translator) IsValid() bool                           { return customT.cfg != nil && customT.cfg.GetAK() != "" }

func (customT *Translator) Translate(ctx context.Context, args *translator.TranslateArgs) (*translator.TranslateRes, error) {
	timeStart := carbon.Now()
	prompt := fmt.Sprintf(
		"You are an excellent simultaneous translator in both %s and %s. I will provide you with %s line by line, and you will translate it into %s line by line. The number of lines in both %s and %s is the same. Do you understand?",
		customT.currentFromLang, args.ToLang, customT.currentFromLang, args.ToLang, customT.currentFromLang, args.ToLang,
	)
	// 初始化
	if customT.currentMessages == nil || customT.currentMessages.Model == "" {
		customT.currentMessages = new(translateRequest)
		customT.currentMessages.Model = customT.cfg.GetProjectKey()
		customT.currentMessages.Messages = []MessageItem{}
	}
	// 预置条件, 语种必须保持一致, 否则就视为新开
	if customT.currentFromLang == "" || customT.currentToLang == "" || customT.currentFromLang != args.FromLang || customT.currentToLang != args.ToLang {
		if len(customT.currentMessages.Messages) > 0 {
			customT.currentMessages.Messages = []MessageItem{}
		}
		customT.currentMessages.Messages = append(customT.currentMessages.Messages, MessageItem{
			Role: "system", Content: prompt,
		})
	}
	customT.currentMessages.Messages = append(customT.currentMessages.Messages, MessageItem{
		Role: "user", Content: args.TextContent,
	})
	cntTokens := 0
	for _, message := range customT.currentMessages.Messages {
		cntTokens += len(message.Content)
	}

	if cntTokens < 2048 {

	} else {
		customT.currentMessages.Messages = []MessageItem{}
		customT.currentMessages.Messages = append(customT.currentMessages.Messages, MessageItem{
			Role: "system", Content: prompt,
		}, MessageItem{
			Role: "user", Content: args.TextContent,
		})
	}

	reqBytes, _ := json.Marshal(customT.currentMessages)
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", customT.cfg.GetAK()),
	}
	respBytes, err := translator.RequestSimpleHttp(ctx, customT, apiTranslate, true, reqBytes, headers)
	if err != nil {
		return nil, err
	}

	resp := new(translateResponse)
	if err = json.Unmarshal(respBytes, resp); err != nil {
		log.Singleton().ErrorF("解析报文异常, 引擎: %s, 错误: %s", customT.GetName(), err)
		return nil, fmt.Errorf("解析报文出现异常, 错误: %s", err.Error())
	}
	if resp.Usage.CompletionTokens <= 0 {
		log.Singleton().ErrorF("接口响应异常, 引擎: %s, 错误: 无响应, 数据: %s", customT.GetName(), string(respBytes))
		return nil, fmt.Errorf("接口响应异常, 引擎: %s, 错误: %s", customT.GetName(), "无响应")
	}

	srcTexts := strings.Split(args.TextContent, customT.GetSep())
	tgtTexts := strings.Split(resp.Choices[0].Message.Content, customT.GetSep())
	ret := new(translator.TranslateRes)
	if len(srcTexts) != len(tgtTexts) {
		ret.TimeUsed = int(carbon.Now().DiffAbsInSeconds(timeStart))
		return ret, nil
	}

	for textIdx, textTarget := range tgtTexts {
		ret.Results = append(ret.Results, &translator.TranslateResBlock{
			Id: srcTexts[textIdx], TextTranslated: textTarget,
		})
	}

	ret.TimeUsed = int(carbon.Now().DiffAbsInSeconds(timeStart))
	return ret, nil
}

type translateRequest struct {
	Model    string        `json:"model"`
	Messages []MessageItem `json:"messages"`
}

type MessageItem struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type translateResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int         `json:"index"`
		Message      MessageItem `json:"message"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}
