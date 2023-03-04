package aliyun

import (
	"context"
	"fmt"
	alimt20181012 "github.com/alibabacloud-go/alimt-20181012/v2/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/golang-module/carbon"
	"gui.subtitle/src/srv/mt"
	"gui.subtitle/src/util/lang"
)

type Cfg struct {
	AccessKeyId     string
	AccessKeySecret string
	Location        string
	Endpoint        string
	IsDebug         bool
}

type ALiMT struct {
	cfg      *Cfg
	mtClient *alimt20181012.Client
}

func (m *ALiMT) Init(_ context.Context, cfg interface{}) error {
	if _, ok := cfg.(*Cfg); !ok {
		return fmt.Errorf("the cfg's mismatched")
	}
	if m.cfg != nil { // 拒绝重复初始化
		return nil
	}
	cfg.(*Cfg).Endpoint = "mt.aliyuncs.com"
	m.cfg = cfg.(*Cfg)
	return m.initClient()
}

func (m *ALiMT) GetCfg() interface{} {
	return m.cfg
}

type TextBatchTranslateArg struct {
	Scene        string
	ApiType      string
	SourceText   string
	ToLanguage   string
	FromLanguage string
}

func (arg *TextBatchTranslateArg) New(text string) *TextBatchTranslateArg {
	arg.Scene = "general"
	arg.ApiType = "translate_standard"
	arg.SourceText = text
	arg.ToLanguage = lang.ZH.ToString()
	arg.FromLanguage = lang.EN.ToString()
	return arg
}

func (m *ALiMT) TextTranslate(context.Context, interface{}) ([]mt.TextTranslateResp, error) {
	return nil, nil
}

func (m *ALiMT) TextBatchTranslate(_ context.Context, args interface{}) ([]mt.TextTranslateResp, error) {
	if _, ok := args.(*TextBatchTranslateArg); !ok {
		return nil, fmt.Errorf("the args for ALiMT.TextBatchTranslate mismatched")
	}
	getBatchTranslateRequest := &alimt20181012.GetBatchTranslateRequest{
		FormatType: tea.String("text"), Scene: tea.String(args.(*TextBatchTranslateArg).Scene),
		ApiType:        tea.String(args.(*TextBatchTranslateArg).ApiType),
		SourceText:     tea.String(args.(*TextBatchTranslateArg).SourceText),
		TargetLanguage: tea.String(args.(*TextBatchTranslateArg).ToLanguage), SourceLanguage: tea.String(args.(*TextBatchTranslateArg).FromLanguage),
	}
	runtime := &util.RuntimeOptions{}
	resp, err := m.mtClient.GetBatchTranslateWithOptions(getBatchTranslateRequest, runtime)
	if err != nil {
		return nil, err
	}
	if tea.Int32Value(resp.Body.Code) != 200 {
		return nil, fmt.Errorf("current request broken, the code: %d", tea.Int32Value(resp.Body.Code))
	}
	var funcResp []mt.TextTranslateResp
	for blockIdx, blockTranslated := range resp.Body.TranslatedList {
		if blockTranslated["code"].(string) != "200" && m.cfg.IsDebug {
			return nil, fmt.Errorf(
				"[%s]%s, 错判翻译异常[%d], 索引: %s",
				carbon.Now(), m.GetName(), blockIdx, blockTranslated["index"],
			)
		}
		funcResp = append(funcResp, mt.TextTranslateResp{
			Idx:           blockTranslated["index"].(string),
			StrTranslated: blockTranslated["translated"].(string),
		})
	}
	return funcResp, nil
}

func (m *ALiMT) GetId() mt.Id {
	return mt.IdALiYun
}

func (m *ALiMT) GetName() string {
	return mt.EngineALiYun.GetZH()
}

func (m *ALiMT) initClient() error {
	mtCfg := &openapi.Config{
		AccessKeyId:     tea.String(m.cfg.AccessKeyId),
		AccessKeySecret: tea.String(m.cfg.AccessKeySecret),
		Endpoint:        tea.String(m.cfg.Endpoint),
	}
	result, err := alimt20181012.NewClient(mtCfg)
	if err != nil {
		return err
	}
	m.mtClient = result
	return nil
}
