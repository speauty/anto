package tencent_cloud_mt

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-module/carbon"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tencentHttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"strings"
	"sync"
	"translator/tst/tt_log"
	"translator/tst/tt_translator"
)

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
		id:            "tecent_mt",
		name:          "腾讯云MT",
		qps:           50,
		procMax:       20,
		textMaxLen:    2000,
		sep:           "\n",
		isClientOk:    false,
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
	isClientOk    bool
	tencentClient *common.Client
}

func (customT *Translator) Init(cfg interface{}) {
	customT.cfg = cfg.(*Cfg)
	if customT.cfg == nil {
		tt_log.GetInstance().Error(fmt.Sprintf("引擎: %s, 错误: 未设置必需参数", customT.GetName()))
		return
	}
	customT.clientBuilder()
}

func (customT *Translator) GetId() string                           { return customT.id }
func (customT *Translator) GetName() string                         { return customT.name }
func (customT *Translator) GetCfg() interface{}                     { return nil }
func (customT *Translator) GetQPS() int                             { return customT.qps }
func (customT *Translator) GetProcMax() int                         { return customT.procMax }
func (customT *Translator) GetTextMaxLen() int                      { return customT.textMaxLen }
func (customT *Translator) GetLangSupported() []tt_translator.LangK { return customT.langSupported }
func (customT *Translator) GetSep() string                          { return customT.sep }
func (customT *Translator) IsValid() bool {
	return customT.isClientOk && customT.cfg.SecretId != "" && customT.cfg.SecretKey != "" && customT.cfg.Region != ""
}

func (customT *Translator) Translate(args *tt_translator.TranslateArgs) (*tt_translator.TranslateRes, error) {
	timeStart := carbon.Now()
	req := &remoteReq{
		BaseRequest: &tencentHttp.BaseRequest{},
		Source:      &args.FromLang,
		Target:      &args.ToLang,
		ProjectId:   &customT.cfg.ProjectId,
	}
	texts := strings.Split(args.TextContent, customT.GetSep())
	for idx, _ := range texts {
		req.SourceTextList = append(req.SourceTextList, &texts[idx])
	}

	req.Init().WithApiInfo("tmt", "2018-03-21", "TextTranslateBatch")
	req.SetContext(context.Background())
	if customT.tencentClient.GetCredential() == nil {
		tt_log.GetInstance().Error(fmt.Sprintf("引擎: %s, 错误: 获取凭证失败", customT.GetName()))
		return nil, fmt.Errorf("引擎: %s, 错误: 获取凭证失败", customT.GetName())
	}
	resp := &remoteResp{
		BaseResponse: &tencentHttp.BaseResponse{},
	}
	err := customT.tencentClient.Send(req, resp)
	if err != nil {
		return nil, err
	}
	if len(resp.Response.TargetTextList) != len(texts) {
		return nil, fmt.Errorf("引擎: %s, 错误: 译文和原文数量匹配失败", customT.GetName())
	}
	ret := new(tt_translator.TranslateRes)
	for idx, textTranslated := range resp.Response.TargetTextList {
		ret.Results = append(ret.Results, &tt_translator.TranslateResBlock{
			Id:             texts[idx],
			TextTranslated: *textTranslated,
		})
	}

	ret.TimeUsed = int(carbon.Now().DiffInSeconds(timeStart))
	return ret, nil
}

func (customT *Translator) clientBuilder() {
	if customT.cfg.SecretId == "" || customT.cfg.SecretKey == "" || customT.cfg.Region == "" {
		return
	}
	tmpClient, tmpErr := common.NewClientWithSecretId(customT.cfg.SecretId, customT.cfg.SecretKey, customT.cfg.Region)
	if tmpErr != nil {
		tt_log.GetInstance().Error(fmt.Sprintf("引擎: %s, 错误: 初始化客户都安失败(%s)", customT.GetName(), tmpErr))
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
