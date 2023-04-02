package huawei_cloud_nlp

import (
	"errors"
	"fmt"
	"github.com/golang-module/carbon"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	nlp "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/nlp/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/nlp/v2/model"
	region "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/nlp/v2/region"
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
		id:         "huawei_cloud_nlp",
		name:       "华为云NLP",
		qps:        20,
		procMax:    10,
		textMaxLen: 2000,
		sep:        "\n",
		langSupported: []tt_translator.LangK{
			{"zh", "中文(简体)"},
			{"en", "英语"},
			{"ru", "俄语"},
			{"fr", "法语"},
			{"de", "德语"},
			{"ko", "韩语"},
			{"ja", "日语"},
			{"th", "泰语"},
			{"ar", "阿拉伯语"},
			{"pt", "葡萄牙语"},
			{"tr", "土耳其语"},
			{"es", "西班牙语"},
			{"vi", "越南语"},
		},
	}
}

var langTo = map[string]model.TextTranslationReqTo{
	"zh": model.GetTextTranslationReqToEnum().ZH,
	"en": model.GetTextTranslationReqToEnum().EN,
	"ru": model.GetTextTranslationReqToEnum().RU,
	"fr": model.GetTextTranslationReqToEnum().FR,
	"de": model.GetTextTranslationReqToEnum().DE,
	"ko": model.GetTextTranslationReqToEnum().KO,
	"ja": model.GetTextTranslationReqToEnum().JA,
	"th": model.GetTextTranslationReqToEnum().TH,
	"ar": model.GetTextTranslationReqToEnum().AR,
	"pt": model.GetTextTranslationReqToEnum().PT,
	"tr": model.GetTextTranslationReqToEnum().TR,
	"es": model.GetTextTranslationReqToEnum().ES,
	"vi": model.GetTextTranslationReqToEnum().VI,
}

var langFrom = map[string]model.TextTranslationReqFrom{
	"zh": model.GetTextTranslationReqFromEnum().ZH,
	"en": model.GetTextTranslationReqFromEnum().EN,
	"ru": model.GetTextTranslationReqFromEnum().RU,
	"fr": model.GetTextTranslationReqFromEnum().FR,
	"de": model.GetTextTranslationReqFromEnum().DE,
	"ko": model.GetTextTranslationReqFromEnum().KO,
	"ja": model.GetTextTranslationReqFromEnum().JA,
	"th": model.GetTextTranslationReqFromEnum().TH,
	"ar": model.GetTextTranslationReqFromEnum().AR,
	"pt": model.GetTextTranslationReqFromEnum().PT,
	"tr": model.GetTextTranslationReqFromEnum().TR,
	"es": model.GetTextTranslationReqFromEnum().ES,
	"vi": model.GetTextTranslationReqFromEnum().VI,
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

func (customT *Translator) Init(cfg interface{}) {
	customT.cfg = cfg.(*Cfg)
	if customT.cfg.Region == "" {
		customT.cfg.Region = "cn-north-4"
	}
}

func (customT *Translator) GetId() string                           { return customT.id }
func (customT *Translator) GetName() string                         { return customT.name }
func (customT *Translator) GetCfg() interface{}                     { return customT.cfg }
func (customT *Translator) GetQPS() int                             { return customT.qps }
func (customT *Translator) GetProcMax() int                         { return customT.procMax }
func (customT *Translator) GetTextMaxLen() int                      { return customT.textMaxLen }
func (customT *Translator) GetLangSupported() []tt_translator.LangK { return customT.langSupported }
func (customT *Translator) GetSep() string                          { return customT.sep }

func (customT *Translator) IsValid() bool {
	return customT.cfg != nil && customT.cfg.AKId != "" && customT.cfg.SkKey != "" && customT.cfg.ProjectId != ""
}

func (customT *Translator) Translate(args *tt_translator.TranslateArgs) (*tt_translator.TranslateRes, error) {
	timeStart := carbon.Now()
	ret := new(tt_translator.TranslateRes)

	request := &model.RunTextTranslationRequest{}
	sceneTextTranslationReq := model.GetTextTranslationReqSceneEnum().COMMON
	request.Body = &model.TextTranslationReq{
		Scene: &sceneTextTranslationReq,
		To:    langTo[args.ToLang],
		From:  langFrom[args.FromLang],
		Text:  args.TextContent,
	}

	resp, err := customT.getClient().RunTextTranslation(request)

	if err != nil {
		tt_log.GetInstance().Error(fmt.Sprintf("调用接口失败, 引擎: %s, 错误: %s", customT.GetName(), err))
		return nil, fmt.Errorf("调用接口失败(%s)", err)
	}
	if resp.ErrorCode != nil && *resp.ErrorCode != "" {
		tt_log.GetInstance().Error(fmt.Sprintf("接口响应错误, 引擎: %s, 错误: %s(%s)", customT.GetName(), *resp.ErrorCode, *resp.ErrorMsg))
		return nil, fmt.Errorf("响应错误(代码: %s, 错误: %s)", *resp.ErrorCode, *resp.ErrorMsg)
	}

	srcTexts := strings.Split(*resp.SrcText, customT.GetSep())
	translatedTexts := strings.Split(*resp.TranslatedText, customT.GetSep())
	if len(srcTexts) != len(translatedTexts) {
		tt_log.GetInstance().Error(fmt.Sprintf("响应解析错误, 引擎: %s, 错误: 译文和原文数量匹配失败", customT.GetName()))
		return nil, errors.New("译文和原文数量匹配失败")
	}
	for idx, text := range srcTexts {
		ret.Results = append(ret.Results, &tt_translator.TranslateResBlock{
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
