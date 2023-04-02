package ling_va

import (
	"encoding/json"
	"fmt"
	"github.com/golang-module/carbon"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
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
		id:         "lingva",
		name:       "Lingva",
		qps:        10,
		procMax:    20,
		textMaxLen: 2000,
		sep:        "\n",
		langSupported: []tt_translator.LangK{
			{"zh", "中文"},
			{"en", "英语"},
			{"ja", "日语"},
			{"ru", "俄语"},
			{"fr", "法语"},
			{"ko", "韩语"},
			{"de", "德语"},
			{"da", "丹麦语"},
			{"nl", "荷兰语"},
			{"it", "意大利语"},
			{"ar", "阿拉伯语"},
			{"be", "白俄罗斯语"},
		},
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
}

func (customT *Translator) Init(cfg interface{}) { customT.cfg = cfg.(*Cfg) }

func (customT *Translator) GetId() string                           { return customT.id }
func (customT *Translator) GetName() string                         { return customT.name }
func (customT *Translator) GetCfg() interface{}                     { return nil }
func (customT *Translator) GetQPS() int                             { return customT.qps }
func (customT *Translator) GetProcMax() int                         { return customT.procMax }
func (customT *Translator) GetTextMaxLen() int                      { return customT.textMaxLen }
func (customT *Translator) GetLangSupported() []tt_translator.LangK { return customT.langSupported }
func (customT *Translator) GetSep() string                          { return customT.sep }
func (customT *Translator) IsValid() bool                           { return customT.cfg.DataId != "" }

func (customT *Translator) Translate(args *tt_translator.TranslateArgs) (*tt_translator.TranslateRes, error) {
	timeStart := carbon.Now()

	var api = fmt.Sprintf("https://lingva.ml/_next/data/%s", customT.cfg.DataId)
	queryUrl := fmt.Sprintf(
		"%s/%s/%s/%s.json", api,
		args.FromLang, args.ToLang, url.PathEscape(args.TextContent),
	)
	httpResp, err := http.DefaultClient.Get(queryUrl)
	defer func() {
		if httpResp != nil {
			_ = httpResp.Body.Close()
		}
	}()
	if err != nil {
		return nil, fmt.Errorf("网络请求出现异常, 错误: %s", err.Error())
	}
	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取报文出现异常, 错误: %s", err.Error())
	}
	lingVaResp := new(lingVaMTResp)
	if err := json.Unmarshal(respBytes, lingVaResp); err != nil {
		return nil, fmt.Errorf("解析报文出现异常, 错误: %s", err.Error())
	}
	if lingVaResp.State == false {
		return nil, fmt.Errorf("翻译异常")
	}
	textTranslatedList := strings.Split(lingVaResp.Props.TextTranslated, customT.sep)
	textSourceList := strings.Split(lingVaResp.Props.Params.TextSource, customT.sep)
	if len(textSourceList) != len(textTranslatedList) {
		return nil, fmt.Errorf("翻译异常, 错误: 源文和译文数量不对等")
	}

	ret := new(tt_translator.TranslateRes)
	for textIdx, textSource := range textSourceList {
		ret.Results = append(ret.Results, &tt_translator.TranslateResBlock{
			Id:             textSource,
			TextTranslated: textTranslatedList[textIdx],
		})
	}

	ret.TimeUsed = int(carbon.Now().DiffAbsInSeconds(timeStart))
	return ret, nil

}

type lingVaMTResp struct {
	State bool `json:"__N_SSG"`
	Props struct {
		//Type           int    `json:"type"`
		TextTranslated string `json:"translation"`
		Params         struct {
			//FromLanguage string `json:"source"`
			//ToLanguage   string `json:"target"`
			TextSource string `json:"query"`
		} `json:"initial"`
	} `json:"pageProps"`
}
