package aliyun

import (
	"fmt"
	alimt20181012 "github.com/alibabacloud-go/alimt-20181012/v2/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/golang-module/carbon"
	"gui.subtitle/src/srv/mt"
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

func (m *ALiMT) Init(cfg interface{}) error {
	if _, ok := cfg.(*Cfg); !ok {
		return fmt.Errorf("the cfg's mismatched")
	}
	if m.cfg != nil { // 拒绝重复初始化
		return nil
	}
	if cfg.(*Cfg).Endpoint == "" { // 拼接域名
		cfg.(*Cfg).Endpoint = fmt.Sprintf("mt.%s.aliyuncs.com", cfg.(*Cfg).Location)
	}
	m.cfg = cfg.(*Cfg)
	return m.initClient()
}

type TextBatchTranslateArg struct {
	Scene      string
	ApiType    string
	SourceText string
	TargetLang string
	SourceLang string
}

func (arg *TextBatchTranslateArg) New(text string) *TextBatchTranslateArg {
	arg.Scene = "general"
	arg.ApiType = "translate_standard"
	arg.SourceText = text
	arg.TargetLang = "zh"
	arg.SourceLang = "en"
	return arg
}

func (m *ALiMT) TextBatchTranslate(args interface{}) ([]mt.TextTranslateResp, error) {
	if _, ok := args.(*TextBatchTranslateArg); !ok {
		return nil, fmt.Errorf("the args for ALiMT.TextBatchTranslate mismatched")
	}
	getBatchTranslateRequest := &alimt20181012.GetBatchTranslateRequest{
		FormatType: tea.String("text"), Scene: tea.String(args.(*TextBatchTranslateArg).Scene),
		ApiType:        tea.String(args.(*TextBatchTranslateArg).ApiType),
		SourceText:     tea.String(args.(*TextBatchTranslateArg).SourceText),
		TargetLanguage: tea.String(args.(*TextBatchTranslateArg).TargetLang), SourceLanguage: tea.String(args.(*TextBatchTranslateArg).SourceLang),
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
			fmt.Println(fmt.Errorf(
				"[%s]%s, 错判翻译异常[%d], 索引: %s",
				carbon.Now(), m.GetName(), blockIdx, blockTranslated["index"],
			))
			continue
		}
		funcResp = append(funcResp, mt.TextTranslateResp{
			Idx:           blockTranslated["index"].(string),
			StrTranslated: blockTranslated["translated"].(string),
		})
	}
	return funcResp, nil
}

func (m *ALiMT) GetName() string {
	return "阿里云-机器翻译"
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
