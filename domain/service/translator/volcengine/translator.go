package volcengine

import (
	"anto/domain/service/translator"
	"anto/lib/log"
	"anto/lib/restrictor"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-module/carbon"
	"github.com/volcengine/volc-sdk-golang/base"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	host    = "open.volcengineapi.com"
	service = "translate"
	version = "2020-06-01"
)

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
		id:            "volc_engine",
		name:          "火山引擎",
		sep:           "\n",
		langSupported: langSupported,
	}
}

type Translator struct {
	id            string
	name          string
	cfg           translator.ImplConfig
	langSupported []translator.LangPair
	sep           string
	mtClient      *base.Client
}

func (customT *Translator) Init(cfg translator.ImplConfig) { customT.cfg = cfg }

func (customT *Translator) GetId() string                           { return customT.id }
func (customT *Translator) GetShortId() string                      { return "ve" }
func (customT *Translator) GetName() string                         { return customT.name }
func (customT *Translator) GetCfg() translator.ImplConfig           { return customT.cfg }
func (customT *Translator) GetLangSupported() []translator.LangPair { return customT.langSupported }
func (customT *Translator) GetSep() string                          { return customT.sep }
func (customT *Translator) IsValid() bool {
	return customT.cfg != nil && customT.cfg.GetAK() != "" && customT.cfg.GetSK() != ""
}

func (customT *Translator) Translate(ctx context.Context, args *translator.TranslateArgs) (*translator.TranslateRes, error) {
	timeStart := carbon.Now()

	params := &translateRequestParams{
		SourceLanguage: args.FromLang,
		TargetLanguage: args.ToLang,
	}
	params.TextList = append(params.TextList, args.TextContent)
	jsonBytes, _ := json.Marshal(params)
	if err := restrictor.Singleton().Wait(customT.GetId(), ctx); err != nil {
		return nil, fmt.Errorf("限流异常, 错误: %s", err.Error())
	}
	respBytes, _, err := customT.client().Json("TranslateText", nil, string(jsonBytes))
	if err != nil {
		return nil, err
	}

	resp := new(translateResponse)
	if err = json.Unmarshal(respBytes, resp); err != nil {
		log.Singleton().ErrorF("解析报文异常, 引擎: %s, 错误: %s", customT.GetName(), err)
		return nil, fmt.Errorf("解析报文出现异常, 错误: %s", err.Error())
	}

	if resp.ResponseMetaData.Error.Code != "" {
		log.Singleton().ErrorF("接口响应异常, 引擎: %s, 错误: %s(%s)", customT.GetName(), resp.ResponseMetaData.Error.Message, resp.ResponseMetaData.Error.Code)
		return nil, fmt.Errorf("接口响应异常, 引擎: %s, 错误: %s", customT.GetName(), resp.ResponseMetaData.Error.Message)
	}
	srcTexts := strings.Split(args.TextContent, customT.GetSep())
	tgtTexts := strings.Split(resp.TranslationList[0].Translation, customT.GetSep())
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

func (customT *Translator) client() *base.Client {
	if customT.mtClient == nil {
		serviceInfo := &base.ServiceInfo{
			Timeout: 5 * time.Second, Host: host,
			Header:      http.Header{"Accept": []string{"application/json"}},
			Credentials: base.Credentials{Region: base.RegionCnNorth1, Service: service},
		}

		apiInfoList := map[string]*base.ApiInfo{
			"TranslateText": {
				Method: http.MethodPost, Path: "/",
				Query: url.Values{"Action": []string{"TranslateText"}, "Version": []string{version}},
			},
		}
		client := base.NewClient(serviceInfo, apiInfoList)
		client.SetAccessKey(customT.cfg.GetAK())
		client.SetSecretKey(customT.cfg.GetSK())
		customT.mtClient = client
	}

	return customT.mtClient
}

type translateRequestParams struct {
	SourceLanguage string
	TargetLanguage string
	TextList       []string
}

type translateResponse struct {
	ResponseMetaData struct {
		RequestId string `json:"RequestId"`
		Action    string `json:"Action"`
		Version   string `json:"Version"`
		Service   string `json:"Service"`
		Region    string `json:"Region"`
		Error     struct {
			Code    string `json:"Code"`
			Message string `json:"Message"`
		} `json:"Error"`
	} `json:"ResponseMetaData"`
	TranslationList []struct {
		Translation string `json:"Translation"`
		//DetectedSourceLanguage string      `json:"DetectedSourceLanguage"`
		//Extra                  interface{} `json:"Extra"`
	} `json:"TranslationList"`
}
