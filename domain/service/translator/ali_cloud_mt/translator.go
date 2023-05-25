package ali_cloud_mt

import (
	"anto/domain/service/translator"
	"anto/lib/log"
	"anto/lib/restrictor"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	alimt "github.com/aliyun/alibaba-cloud-sdk-go/services/alimt"
	"github.com/golang-module/carbon"
	"github.com/spf13/cast"
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
		id:            "ali_cloud_mt",
		name:          "阿里云",
		sep:           "\n",
		langSupported: langSupported,
		isClientOk:    false,
	}
}

type Translator struct {
	id            string
	name          string
	cfg           translator.ImplConfig
	langSupported []translator.LangPair
	sep           string
	mtClient      *alimt.Client
	isClientOk    bool
}

func (customT *Translator) Init(cfg translator.ImplConfig) {
	customT.cfg = cfg
	customT.clientBuilder()
}

func (customT *Translator) GetId() string                           { return customT.id }
func (customT *Translator) GetShortId() string                      { return "ac" }
func (customT *Translator) GetName() string                         { return customT.name }
func (customT *Translator) GetCfg() translator.ImplConfig           { return customT.cfg }
func (customT *Translator) GetLangSupported() []translator.LangPair { return customT.langSupported }
func (customT *Translator) GetSep() string                          { return customT.sep }
func (customT *Translator) IsValid() bool {
	return customT.cfg != nil && customT.cfg.GetAK() != "" && customT.cfg.GetSK() != "" && customT.isClientOk == true
}

func (customT *Translator) Translate(ctx context.Context, args *translator.TranslateArgs) (*translator.TranslateRes, error) {
	timeStart := carbon.Now()
	ret := new(translator.TranslateRes)

	texts := strings.Split(args.TextContent, customT.GetSep())
	var txtMap []map[int]string
	for idx, text := range texts {
		if idx%50 == 0 {
			txtMap = append(txtMap, map[int]string{})
		}
		txtMap[idx/50][idx] = text
	}
	for _, currentTxtBlock := range txtMap {
		bytes, _ := json.Marshal(currentTxtBlock)
		req := alimt.CreateGetBatchTranslateRequest()
		req.Scheme = "https"
		req.TargetLanguage = args.ToLang
		req.SourceLanguage = args.FromLang
		req.ApiType = "translate_standard"
		req.FormatType = "text"
		req.Scene = "general"
		req.SourceText = string(bytes)
		if err := restrictor.Singleton().Wait(customT.GetId(), ctx); err != nil {
			return nil, fmt.Errorf("限流异常, 错误: %s", err.Error())
		}
		resp, err := customT.mtClient.GetBatchTranslate(req)
		if err != nil {
			log.Singleton().ErrorF("引擎: %s, 错误: %s", customT.GetName(), err)
			return nil, fmt.Errorf("引擎: %s, 错误: %s", customT.GetName(), err)
		}
		for _, translated := range resp.TranslatedList {
			if translated["code"] != "200" {
				log.Singleton().ErrorF("引擎: %s, 错误: %s", customT.GetName(), translated["errorMsg"].(string))
				return nil, fmt.Errorf("引擎: %s, 错误: %s", customT.GetName(), translated["errorMsg"].(string))
			}
			idx := cast.ToInt(translated["index"].(string))
			ret.Results = append(ret.Results, &translator.TranslateResBlock{
				Id: texts[idx], TextTranslated: translated["translated"].(string),
			})
		}
	}

	ret.TimeUsed = int(carbon.Now().DiffInSeconds(timeStart))
	return ret, nil
}

func (customT *Translator) clientBuilder() {
	if customT.cfg == nil || customT.cfg.GetAK() == "" || customT.cfg.GetSK() == "" {
		return
	}
	config := sdk.NewConfig()

	credential := credentials.NewAccessKeyCredential(customT.cfg.GetAK(), customT.cfg.GetSK())
	client, err := alimt.NewClientWithOptions(customT.cfg.GetRegion(), config, credential)
	if err != nil {
		log.Singleton().ErrorF("引擎: %s, 错误: 生成客户端失败(%s)", customT.GetName(), err)
		return
	}
	customT.mtClient = client
	customT.isClientOk = true
}
