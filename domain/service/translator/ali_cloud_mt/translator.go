package ali_cloud_mt

import (
	"anto/domain/service/translator"
	"anto/lib/log"
	"anto/lib/restrictor"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	alimt "github.com/aliyun/alibaba-cloud-sdk-go/services/alimt"
	"github.com/golang-module/carbon"
	"github.com/spf13/cast"
	"net"
	"strings"
	"sync"
	"time"
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
		name:          "机器翻译-阿里云",
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

// TestConnection checks the selected public endpoint without making a
// translation request or consuming provider quota.
func (customT *Translator) TestConnection(ctx context.Context, region string) (*translator.ConnectionResult, error) {
	if region == "" && customT.cfg != nil {
		region = customT.cfg.GetRegion()
	}
	endpoint := publicEndpoints[region]
	if endpoint == "" {
		return nil, fmt.Errorf("%w: %s", translator.ErrInvalidConnectionRegion, region)
	}
	ctx, cancel := context.WithTimeout(ctx, 7*time.Second)
	defer cancel()
	started := time.Now()
	dialer := &tls.Dialer{
		NetDialer: &net.Dialer{Timeout: 5 * time.Second},
		Config:    &tls.Config{ServerName: endpoint, MinVersion: tls.VersionTLS12},
	}
	connection, err := dialer.DialContext(ctx, "tcp", net.JoinHostPort(endpoint, "443"))
	if err != nil {
		return nil, fmt.Errorf("连接 %s 失败: %w", endpoint, err)
	}
	_ = connection.Close()
	return &translator.ConnectionResult{Endpoint: endpoint, LatencyMs: time.Since(started).Milliseconds()}, nil
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
	region := customT.cfg.GetRegion()
	client, err := alimt.NewClientWithOptions(region, config, credential)
	if err != nil {
		log.Singleton().ErrorF("引擎: %s, 错误: 生成客户端失败(%s)", customT.GetName(), err)
		return
	}
	if endpoint := publicEndpoints[region]; endpoint != "" {
		client.Domain = endpoint
	}
	customT.mtClient = client
	customT.isClientOk = true
}
