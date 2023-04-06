package ali_cloud_mt

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	alimt "github.com/aliyun/alibaba-cloud-sdk-go/services/alimt"
	"github.com/golang-module/carbon"
	"github.com/spf13/cast"
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
		id:            "ali_cloud_mt",
		name:          "阿里云MT",
		qps:           50,
		procMax:       20,
		textMaxLen:    3000,
		sep:           "\n",
		langSupported: langSupported,
		isClientOk:    false,
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
	mtClient      *alimt.Client
	isClientOk    bool
}

func (customT *Translator) Init(cfg interface{}) {
	customT.cfg = cfg.(*Cfg)
	customT.clientBuilder()
}

func (customT *Translator) GetId() string       { return customT.id }
func (customT *Translator) GetName() string     { return customT.name }
func (customT *Translator) GetCfg() interface{} { return nil }
func (customT *Translator) GetQPS() int         { return customT.qps }
func (customT *Translator) GetProcMax() int     { return customT.procMax }
func (customT *Translator) GetTextMaxLen() int {
	if customT.cfg.MaxSingleTextLength > 0 {
		return customT.cfg.MaxSingleTextLength
	}
	return customT.textMaxLen
}
func (customT *Translator) GetLangSupported() []tt_translator.LangK { return customT.langSupported }
func (customT *Translator) GetSep() string                          { return customT.sep }
func (customT *Translator) IsValid() bool {
	return customT.cfg != nil && customT.cfg.AKId != "" && customT.cfg.AKSecret != "" && customT.isClientOk == true
}

func (customT *Translator) Translate(args *tt_translator.TranslateArgs) (*tt_translator.TranslateRes, error) {
	timeStart := carbon.Now()
	ret := new(tt_translator.TranslateRes)

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
		resp, err := customT.mtClient.GetBatchTranslate(req)
		if err != nil {
			tt_log.GetInstance().Error(fmt.Sprintf("引擎: %s, 错误: %s", customT.GetName(), err))
			return nil, fmt.Errorf("引擎: %s, 错误: %s", customT.GetName(), err)
		}
		for _, translated := range resp.TranslatedList {
			if translated["code"] != "200" {
				tt_log.GetInstance().Error(fmt.Sprintf("引擎: %s, 错误: %s", customT.GetName(), translated["errorMsg"].(string)))
				return nil, fmt.Errorf("引擎: %s, 错误: %s", customT.GetName(), translated["errorMsg"].(string))
			}
			idx := cast.ToInt(translated["index"].(string))
			ret.Results = append(ret.Results, &tt_translator.TranslateResBlock{
				Id:             texts[idx],
				TextTranslated: translated["translated"].(string),
			})
		}
	}

	ret.TimeUsed = int(carbon.Now().DiffInSeconds(timeStart))
	return ret, nil
}

func (customT *Translator) clientBuilder() {
	if customT.cfg == nil || customT.cfg.AKId == "" || customT.cfg.AKSecret == "" {
		return
	}
	config := sdk.NewConfig()

	credential := credentials.NewAccessKeyCredential(customT.cfg.AKId, customT.cfg.AKSecret)
	region := customT.cfg.Region
	if region == "" {
		region = "cn-hangzhou"
	}
	client, err := alimt.NewClientWithOptions(region, config, credential)
	if err != nil {
		tt_log.GetInstance().Error("引擎: %s, 错误: 生成客户端失败(%s)", customT.GetName(), err)
		return
	}
	customT.mtClient = client
	customT.isClientOk = true
}
