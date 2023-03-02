package youdao

import (
	"context"
	"encoding/json"
	"fmt"
	"gui.subtitle/src/srv/mt"
	"gui.subtitle/src/util/lang"
	"io"
	"net/http"
	url2 "net/url"
	"strings"
)

var api = "https://fanyi.youdao.com/translate?&doctype=json"

type MT struct {
}

func (youDaoMT *MT) GetId() mt.Id {
	return mt.IdYouDao
}

func (youDaoMT *MT) GetName() string {
	return mt.EngineYouDao.GetZH()
}

func (youDaoMT *MT) GetCfg() interface{} {
	return nil
}

func (youDaoMT *MT) Init(_ context.Context, _ interface{}) error {
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

type youDaoMTResp struct {
	Type        string `json:"type"`
	ErrorCode   int    `json:"errorCode"`
	ElapsedTime int    `json:"elapsedTime"`
	TransResult [][]struct {
		Src string `json:"src,omitempty"` // 原文
		Tgt string `json:"tgt,omitempty"` // 译文
	} `json:"translateResult"`
}

func (youDaoMT *MT) TextTranslate(ctx context.Context, args interface{}) ([]mt.TextTranslateResp, error) {
	if _, ok := args.(*TextTranslateArg); !ok {
		return nil, fmt.Errorf("the args for YoudaoMT.TextTranslateArg mismatched")
	}
	mtArgs := args.(*TextTranslateArg)
	argType := youDaoMT.convertLanguage2Type(mtArgs.FromLanguage, mtArgs.ToLanguage)
	mtArgs.SourceText = url2.QueryEscape(mtArgs.SourceText)
	url := fmt.Sprintf("%s&type=%s&i=%s", api, argType, mtArgs.SourceText)
	httpResp, err := http.DefaultClient.Get(url)
	defer func() {
		if httpResp.Body != nil {
			_ = httpResp.Body.Close()
		}
	}()
	if err != nil {
		return nil, fmt.Errorf("网络请求[%s]出现异常, 错误: %s", url, err.Error())
	}
	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取报文出现异常, 错误: %s", err.Error())
	}
	youDaoResp := new(youDaoMTResp)
	if err := json.Unmarshal(respBytes, youDaoResp); err != nil {
		fmt.Println(string(respBytes))
		return nil, fmt.Errorf("解析报文出现异常, 错误: %s", err.Error())
	}
	if youDaoResp.ErrorCode != 0 {
		return nil, fmt.Errorf("翻译异常, 代码: %d", youDaoResp.ErrorCode)
	}
	var resp []mt.TextTranslateResp
	for _, transBlockArray := range youDaoResp.TransResult {
		for _, transBlock := range transBlockArray {
			resp = append(resp, mt.TextTranslateResp{
				Idx:           transBlock.Src,
				StrTranslated: transBlock.Tgt,
			})
		}
	}
	return resp, nil
}

func (youDaoMT *MT) TextBatchTranslate(_ context.Context, _ interface{}) ([]mt.TextTranslateResp, error) {
	return nil, nil
}

// ZH_CN2EN 中文　»　英语
// ZH_CN2JA 中文　»　日语
// ZH_CN2KR 中文　»　韩语
// ZH_CN2FR 中文　»　法语
// ZH_CN2RU 中文　»　俄语
// ZH_CN2SP 中文　»　西语
// EN2ZH_CN 英语　»　中文
// JA2ZH_CN 日语　»　中文
// KR2ZH_CN 韩语　»　中文
// FR2ZH_CN 法语　»　中文
// RU2ZH_CN 俄语　»　中文
// SP2ZH_CN 西语　»　中文
func (youDaoMT *MT) convertLanguage2Type(fromLanguage string, toLanguage string) string {
	fromLanguage = strings.ToUpper(fromLanguage)
	toLanguage = strings.ToUpper(toLanguage)
	if fromLanguage == "ZH" {
		fromLanguage = "ZH_CN"
	}
	if toLanguage == "ZH" {
		toLanguage = "ZH_CN"
	}
	return fmt.Sprintf("%s2%s", fromLanguage, toLanguage)
}
