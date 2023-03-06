package lingva

import (
	"context"
	"encoding/json"
	"fmt"
	"gui.subtitle/src/srv/mt"
	"gui.subtitle/src/util/lang"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var api = "https://lingva.ml/_next/data/y9iL3S16LVaUrWMg9xpAa"

type MT struct {
}

func (lingvaMT *MT) GetId() mt.Id {
	return mt.IdLingVa
}

func (lingvaMT *MT) GetName() string {
	return mt.EngineLingVa.GetZH()
}

func (lingvaMT *MT) GetCfg() interface{} {
	return nil
}

func (lingvaMT *MT) Init(_ context.Context, _ interface{}) error {
	return nil
}

type TextTranslateArg struct {
	SourceText   string
	ToLanguage   string
	FromLanguage string
}

func (arg *TextTranslateArg) New(text string) *TextTranslateArg {
	arg.SourceText = text
	arg.ToLanguage = lang.ZH.ToString()
	arg.FromLanguage = lang.EN.ToString()
	return arg
}

type lingVaMTResp struct {
	State bool `json:"__N_SSG"`
	Props struct {
		Type           int    `json:"type"`
		TextTranslated string `json:"translation"`
		Params         struct {
			FromLanguage string `json:"source"`
			ToLanguage   string `json:"target"`
			TextSource   string `json:"query"`
		} `json:"initial"`
	} `json:"pageProps"`
}

func (lingvaMT *MT) TextTranslate(ctx context.Context, args interface{}) ([]mt.TextTranslateResp, error) {
	if _, ok := args.(*TextTranslateArg); !ok {
		return nil, fmt.Errorf("the args for lingvaMT.TextTranslateArg mismatched")
	}
	mtArgs := args.(*TextTranslateArg)
	queryUrl := fmt.Sprintf("%s/%s/%s/%s.json", api, mtArgs.FromLanguage, mtArgs.ToLanguage, url.PathEscape(mtArgs.SourceText))
	httpResp, err := http.DefaultClient.Get(queryUrl)
	defer func() {
		if httpResp.Body != nil {
			_ = httpResp.Body.Close()
		}
	}()
	if err != nil {
		return nil, fmt.Errorf("网络请求[%s]出现异常, 错误: %s", queryUrl, err.Error())
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
		return nil, fmt.Errorf("翻译异常, 链接: %s", queryUrl)
	}
	textTranslatedList := strings.Split(lingVaResp.Props.TextTranslated, mt.BlockSep)
	textSourceList := strings.Split(lingVaResp.Props.Params.TextSource, mt.BlockSep)
	if len(textSourceList) != len(textTranslatedList) {
		return nil, fmt.Errorf("翻译异常, 错误: 源文和译文数量不对等")
	}
	var resp []mt.TextTranslateResp
	for textIdx, textSource := range textSourceList {
		resp = append(resp, mt.TextTranslateResp{
			Idx:           textSource,
			StrTranslated: textTranslatedList[textIdx],
		})
	}
	return resp, nil
}

func (lingvaMT *MT) TextBatchTranslate(_ context.Context, _ interface{}) ([]mt.TextTranslateResp, error) {
	return nil, nil
}
