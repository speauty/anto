package tencent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tencentHttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"gui.subtitle/src/srv/mt"
)

const APIVersion = "2018-03-21"

type Cfg struct {
	SecretId  string
	SecretKey string
	Region    string
}

type MT struct {
	tencentClient *common.Client
	cfg           *Cfg
}

func (m *MT) GetId() mt.Id {
	return mt.IdTencent
}

func (m *MT) GetName() string {
	return mt.EngineTencent.GetZH()
}

func (m *MT) GetCfg() interface{} {
	return m.cfg
}

func (m *MT) Init(_ context.Context, cfg interface{}) error {
	if _, ok := cfg.(*Cfg); !ok {
		return fmt.Errorf("the cfg's mismatched")
	}
	if m.cfg != nil || m.tencentClient != nil {
		return nil
	}
	m.cfg = cfg.(*Cfg)
	tmpClient, tmpErr := common.NewClientWithSecretId(m.cfg.SecretId, m.cfg.SecretKey, m.cfg.Region)
	if tmpErr != nil {
		return tmpErr
	}

	tmpClient.WithProfile(profile.NewClientProfile())
	m.tencentClient = tmpClient
	return nil
}

func (m *MT) TextTranslate(ctx context.Context, args interface{}) ([]mt.TextTranslateResp, error) {
	return nil, nil
}

type TextTranslateBatchRequest struct {
	*tencentHttp.BaseRequest
	Source         *string   `json:"Source,omitempty" name:"Source"`
	Target         *string   `json:"Target,omitempty" name:"Target"`
	ProjectId      *int64    `json:"ProjectId,omitempty" name:"ProjectId"`
	SourceTextList []*string `json:"SourceTextList,omitempty" name:"SourceTextList"`
}

type TextTranslateBatchResponseParams struct {
	Source         *string   `json:"Source,omitempty" name:"Source"`
	Target         *string   `json:"Target,omitempty" name:"Target"`
	TargetTextList []*string `json:"TargetTextList,omitempty" name:"TargetTextList"`
	RequestId      *string   `json:"RequestId,omitempty" name:"RequestId"`
}

type TextTranslateBatchResponse struct {
	*tencentHttp.BaseResponse
	Response *TextTranslateBatchResponseParams `json:"Response"`
}

func (r *TextTranslateBatchResponse) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

func (r *TextTranslateBatchResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type TextBatchTranslateArg struct {
	FromLanguage string
	ToLanguage   string
	TextList     []string
}

func (m *MT) TextBatchTranslate(ctx context.Context, args interface{}) ([]mt.TextTranslateResp, error) {
	if _, ok := args.(*TextBatchTranslateArg); !ok {
		return nil, fmt.Errorf("the args for ALiMT.TextBatchTranslate mismatched")
	}
	projectId := int64(0)
	request := &TextTranslateBatchRequest{
		BaseRequest: &tencentHttp.BaseRequest{},
		Source:      &args.(*TextBatchTranslateArg).FromLanguage,
		Target:      &args.(*TextBatchTranslateArg).ToLanguage,
		ProjectId:   &projectId,
	}
	for idx, _ := range args.(*TextBatchTranslateArg).TextList {
		request.SourceTextList = append(request.SourceTextList, &args.(*TextBatchTranslateArg).TextList[idx])
	}

	request.Init().WithApiInfo("tmt", APIVersion, "TextTranslateBatch")
	request.SetContext(ctx)
	if m.tencentClient.GetCredential() == nil {
		return nil, errors.New("文本批量翻译需要凭证")
	}
	response := &TextTranslateBatchResponse{
		BaseResponse: &tencentHttp.BaseResponse{},
	}
	err := m.tencentClient.Send(request, response)
	if err != nil {
		return nil, err
	}
	var resp []mt.TextTranslateResp
	for sourceIdx, sourceText := range args.(*TextBatchTranslateArg).TextList {
		resp = append(resp, mt.TextTranslateResp{
			Idx:           sourceText,
			StrTranslated: *response.Response.TargetTextList[sourceIdx],
		})
	}
	return resp, nil
}
