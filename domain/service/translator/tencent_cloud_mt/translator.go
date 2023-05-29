package tencent_cloud_mt

import (
	"anto/domain/service/translator"
	"anto/lib/log"
	"anto/lib/restrictor"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-module/carbon"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tencentHttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
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
		id:            "tencent_cloud_mt",
		name:          "腾讯云",
		sep:           "\n",
		isClientOk:    false,
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
	isClientOk    bool
	tencentClient *common.Client
}

func (customT *Translator) Init(cfg translator.ImplConfig) {
	customT.cfg = cfg
	customT.clientBuilder()
}

func (customT *Translator) GetId() string                           { return customT.id }
func (customT *Translator) GetShortId() string                      { return "tc" }
func (customT *Translator) GetName() string                         { return customT.name }
func (customT *Translator) GetCfg() translator.ImplConfig           { return customT.cfg }
func (customT *Translator) GetLangSupported() []translator.LangPair { return customT.langSupported }
func (customT *Translator) GetSep() string                          { return customT.sep }
func (customT *Translator) IsValid() bool {
	return customT.isClientOk && customT.cfg.GetAK() != "" && customT.cfg.GetSK() != "" && customT.cfg.GetRegion() != ""
}

func (customT *Translator) Translate(ctx context.Context, args *translator.TranslateArgs) (*translator.TranslateRes, error) {
	timeStart := carbon.Now()
	req := &remoteReq{
		BaseRequest: &tencentHttp.BaseRequest{},
		Source:      &args.FromLang,
		Target:      &args.ToLang,
		ProjectId:   customT.cfg.(*Config).GetProjectKeyPtr(),
	}
	texts := strings.Split(args.TextContent, customT.GetSep())
	for idx, _ := range texts {
		req.SourceTextList = append(req.SourceTextList, &texts[idx])
	}

	req.Init().WithApiInfo("tmt", "2018-03-21", "TextTranslateBatch")
	req.SetContext(ctx)
	if customT.tencentClient.GetCredential() == nil {
		log.Singleton().ErrorF("引擎: %s, 错误: 获取凭证失败", customT.GetName())
		return nil, fmt.Errorf("引擎: %s, 错误: 获取凭证失败", customT.GetName())
	}
	resp := &remoteResp{
		BaseResponse: &tencentHttp.BaseResponse{},
	}
	if err := restrictor.Singleton().Wait(customT.GetId(), ctx); err != nil {
		return nil, fmt.Errorf("限流异常, 错误: %s", err.Error())
	}
	err := customT.tencentClient.Send(req, resp)
	if err != nil {
		return nil, err
	}
	if len(resp.Response.TargetTextList) != len(texts) {
		return nil, translator.ErrSrcAndTgtNotMatched
	}
	ret := new(translator.TranslateRes)
	for idx, textTranslated := range resp.Response.TargetTextList {
		ret.Results = append(ret.Results, &translator.TranslateResBlock{
			Id:             texts[idx],
			TextTranslated: *textTranslated,
		})
	}

	ret.TimeUsed = int(carbon.Now().DiffInSeconds(timeStart))
	return ret, nil
}

func (customT *Translator) clientBuilder() {
	if customT.cfg.GetAK() == "" || customT.cfg.GetSK() == "" || customT.cfg.GetRegion() == "" {
		return
	}
	tmpClient, tmpErr := common.NewClientWithSecretId(customT.cfg.GetAK(), customT.cfg.GetSK(), customT.cfg.GetRegion())
	if tmpErr != nil {
		log.Singleton().ErrorF("引擎: %s, 错误: 初始化客户都安失败(%s)", customT.GetName(), tmpErr)
		return
	}

	tmpClient.WithProfile(profile.NewClientProfile())
	customT.tencentClient = tmpClient
	customT.isClientOk = true
	return
}

type remoteReq struct {
	*tencentHttp.BaseRequest
	Source         *string   `json:"Source,omitempty" name:"Source"`
	Target         *string   `json:"Target,omitempty" name:"Target"`
	ProjectId      *int64    `json:"ProjectId,omitempty" name:"ProjectId"`
	SourceTextList []*string `json:"SourceTextList,omitempty" name:"SourceTextList"`
}

type remoteResp struct {
	*tencentHttp.BaseResponse
	Response *remoteRespParam `json:"Response"`
}

func (r *remoteResp) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

func (r *remoteResp) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type remoteRespParam struct {
	Source         *string   `json:"Source,omitempty" name:"Source"`
	Target         *string   `json:"Target,omitempty" name:"Target"`
	TargetTextList []*string `json:"TargetTextList,omitempty" name:"TargetTextList"`
	RequestId      *string   `json:"RequestId,omitempty" name:"RequestId"`
}
