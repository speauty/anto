package huawei_cloud_nlp

import (
	"anto/dependency/service/translator"
	"anto/lib/log"
	"anto/lib/restrictor"
	"context"
	"fmt"
	"github.com/golang-module/carbon"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	nlp "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/nlp/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/nlp/v2/model"
	region "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/nlp/v2/region"
	"strings"
	"sync"
)

var (
	apiTranslator  *Translator
	onceTranslator sync.Once
)

func Singleton() *Translator {
	onceTranslator.Do(func() {
		apiTranslator = New()
	})
	return apiTranslator
}

func New() *Translator {
	return &Translator{
		id:            "huawei_cloud_nlp",
		name:          "华为云NLP",
		qps:           20,
		procMax:       10,
		textMaxLen:    2000,
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
	langSupported []translator.LangPair
	sep           string
}

func (customT *Translator) Init(cfg interface{}) {
	customT.cfg = cfg.(*Cfg)
	if customT.cfg.Region == "" {
		customT.cfg.Region = "cn-north-4"
	}
}

func (customT *Translator) GetId() string       { return customT.id }
func (customT *Translator) GetName() string     { return customT.name }
func (customT *Translator) GetCfg() interface{} { return customT.cfg }
func (customT *Translator) GetQPS() int         { return customT.qps }
func (customT *Translator) GetProcMax() int     { return customT.procMax }
func (customT *Translator) GetTextMaxLen() int {
	if customT.cfg.MaxSingleTextLength > 0 {
		return customT.cfg.MaxSingleTextLength
	}
	return customT.textMaxLen
}
func (customT *Translator) GetLangSupported() []translator.LangPair { return customT.langSupported }
func (customT *Translator) GetSep() string                          { return customT.sep }

func (customT *Translator) IsValid() bool {
	return customT.cfg != nil && customT.cfg.AKId != "" && customT.cfg.SkKey != "" && customT.cfg.ProjectId != ""
}

func (customT *Translator) Translate(ctx context.Context, args *translator.TranslateArgs) (*translator.TranslateRes, error) {
	timeStart := carbon.Now()
	ret := new(translator.TranslateRes)

	request := &model.RunTextTranslationRequest{}
	sceneTextTranslationReq := model.GetTextTranslationReqSceneEnum().COMMON
	request.Body = &model.TextTranslationReq{
		Scene: &sceneTextTranslationReq,
		To:    langTo[args.ToLang],
		From:  langFrom[args.FromLang],
		Text:  args.TextContent,
	}
	if err := restrictor.Singleton().Wait(customT.GetId(), ctx); err != nil {
		return nil, fmt.Errorf("限流异常, 错误: %s", err.Error())
	}
	resp, err := customT.getClient().RunTextTranslation(request)

	if err != nil {
		log.Singleton().ErrorF("调用接口失败, 引擎: %s, 错误: %s", customT.GetName(), err)
		return nil, fmt.Errorf("调用接口失败(%s)", err)
	}
	if resp.ErrorCode != nil && *resp.ErrorCode != "" {
		log.Singleton().ErrorF("接口响应错误, 引擎: %s, 错误: %s(%s)", customT.GetName(), *resp.ErrorCode, *resp.ErrorMsg)
		return nil, fmt.Errorf("响应错误(代码: %s, 错误: %s)", *resp.ErrorCode, *resp.ErrorMsg)
	}

	srcTexts := strings.Split(*resp.SrcText, customT.GetSep())
	translatedTexts := strings.Split(*resp.TranslatedText, customT.GetSep())
	if len(srcTexts) != len(translatedTexts) {
		return nil, translator.ErrSrcAndTgtNotMatched
	}
	for idx, text := range srcTexts {
		ret.Results = append(ret.Results, &translator.TranslateResBlock{
			Id:             text,
			TextTranslated: translatedTexts[idx],
		})
	}

	ret.TimeUsed = int(carbon.Now().DiffAbsInSeconds(timeStart))
	return ret, nil

}

func (customT *Translator) getAuth() *basic.Credentials {
	return basic.NewCredentialsBuilder().WithAk(customT.cfg.AKId).WithSk(customT.cfg.SkKey).Build()
}

func (customT *Translator) getClient() *nlp.NlpClient {
	return nlp.NewNlpClient(
		nlp.NlpClientBuilder().
			WithRegion(region.ValueOf("cn-north-4")).
			WithCredential(customT.getAuth()).
			Build(),
	)
}
