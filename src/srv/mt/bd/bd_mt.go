// Package bd 百度翻译-通用文本翻译 @link https://api.fanyi.baidu.com/doc/21
package bd

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/golang-module/carbon"
	"gui.subtitle/src/srv/mt"
	"gui.subtitle/src/util/lang"
	"io"
	"net/http"
	url2 "net/url"
)

var apiUrl = "http://api.fanyi.baidu.com/api/trans/vip/translate"

type Cfg struct {
	AppId     string
	AppSecret string
}

type MT struct {
	cfg *Cfg
}

func (m *MT) Init(_ context.Context, cfg interface{}) error {
	if _, ok := cfg.(*Cfg); !ok {
		return fmt.Errorf("the cfg's mismatched")
	}
	if m.cfg != nil { // 拒绝重复初始化
		return nil
	}
	m.cfg = cfg.(*Cfg)
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

type bdMTResp struct {
	From        string `json:"from,omitempty"`
	To          string `json:"to,omitempty"`
	TransResult []struct {
		Src string `json:"src,omitempty"`
		Dst string `json:"dst,omitempty"`
	} `json:"trans_result"`
	ErrorCode string `json:"error_code,omitempty"`
	ErrorMsg  string `json:"error_msg,omitempty"`
}

func (m *MT) TextTranslate(_ context.Context, args interface{}) (*mt.TextTranslateResp, error) {
	if _, ok := args.(*TextTranslateArg); !ok {
		return nil, fmt.Errorf("the args for BaiduMT.TextTranslateArg mismatched")
	}
	ttArgs := args.(*TextTranslateArg)
	salt := carbon.Now().Format("20060102150405")
	sign := m.GenSign(ttArgs.SourceText, salt)
	url := fmt.Sprintf(
		"%s?q=%s&from=%s&to=%s&appid=%s&salt=%s&sign=%s",
		apiUrl, url2.QueryEscape(ttArgs.SourceText), ttArgs.FromLanguage, ttArgs.ToLanguage,
		m.cfg.AppId, salt, sign,
	)
	httpResp, err := http.DefaultClient.Get(url)
	defer func() { _ = httpResp.Body.Close() }()
	if err != nil {
		return nil, fmt.Errorf("网络请求[%s]出现异常, 错误: %s", url, err.Error())
	}
	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取报文出现异常, 错误: %s", err.Error())
	}
	bdMT := &bdMTResp{}
	if err = json.Unmarshal(respBytes, bdMT); err != nil {
		return nil, fmt.Errorf("解析报文出现异常, 错误: %s", err.Error())
	}

	if bdMT.ErrorMsg != "" {
		return nil, fmt.Errorf("翻译异常, 代码: %s, 错误: %s", bdMT.ErrorCode, bdMT.ErrorMsg)
	}
	resp := &mt.TextTranslateResp{
		Idx:           bdMT.TransResult[0].Src,
		StrTranslated: bdMT.TransResult[0].Dst,
	}
	return resp, nil
}

func (m *MT) TextBatchTranslate(context.Context, interface{}) ([]mt.TextTranslateResp, error) {
	return nil, nil
}

func (m *MT) GetName() string {
	return "百度翻译"
}

func (m *MT) GetId() mt.Id {
	return mt.BAIDU
}

func (m *MT) GenSign(queryStr, saltStr string) string {
	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%s%s%s%s", m.cfg.AppId, queryStr, saltStr, m.cfg.AppSecret)))
	return hex.EncodeToString(h.Sum(nil))
}
