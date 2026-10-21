package tencent_cloud_mt

import (
	"anto/domain/service/translator"
	"anto/lib/restrictor"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	apiTranslator  *Translator
	onceTranslator sync.Once
)

func API() *Translator {
	onceTranslator.Do(func() { apiTranslator = New() })
	return apiTranslator
}

type Translator struct {
	cfg      translator.ImplConfig
	client   *http.Client
	endpoint string
}

func New() *Translator {
	return &Translator{
		client:   &http.Client{Timeout: 15 * time.Second},
		endpoint: "https://" + publicEndpoint,
	}
}

func (t *Translator) Init(config translator.ImplConfig) { t.cfg = config }
func (t *Translator) GetId() string                     { return "tencent_cloud_mt" }
func (t *Translator) GetShortId() string                { return "tc" }
func (t *Translator) GetName() string                   { return "机器翻译-腾讯云" }
func (t *Translator) GetCfg() translator.ImplConfig     { return t.cfg }
func (t *Translator) GetLangSupported() []translator.LangPair {
	return langSupported
}
func (t *Translator) GetSep() string { return "" }
func (t *Translator) IsValid() bool {
	return t.cfg != nil && strings.TrimSpace(t.cfg.GetAK()) != "" &&
		strings.TrimSpace(t.cfg.GetSK()) != "" && validRegion(t.cfg.GetRegion())
}

// TestConnection checks only the public API's TLS connectivity. It does not
// send a translation request or consume the account's quota.
func (t *Translator) TestConnection(ctx context.Context, region string) (*translator.ConnectionResult, error) {
	if region == "" && t.cfg != nil {
		region = t.cfg.GetRegion()
	}
	if !validRegion(region) {
		return nil, fmt.Errorf("%w: %s", translator.ErrInvalidConnectionRegion, region)
	}
	ctx, cancel := context.WithTimeout(ctx, 7*time.Second)
	defer cancel()
	started := time.Now()
	dialer := &tls.Dialer{
		NetDialer: &net.Dialer{Timeout: 5 * time.Second},
		Config:    &tls.Config{ServerName: publicEndpoint, MinVersion: tls.VersionTLS12},
	}
	connection, err := dialer.DialContext(ctx, "tcp", net.JoinHostPort(publicEndpoint, "443"))
	if err != nil {
		return nil, fmt.Errorf("连接 %s 失败: %w", publicEndpoint, err)
	}
	_ = connection.Close()
	return &translator.ConnectionResult{Endpoint: publicEndpoint, LatencyMs: time.Since(started).Milliseconds()}, nil
}

func (t *Translator) Translate(ctx context.Context, args *translator.TranslateArgs) (*translator.TranslateRes, error) {
	if !t.IsValid() {
		return nil, fmt.Errorf("腾讯云翻译引擎未配置 SecretId、SecretKey 或有效地域")
	}
	if args == nil {
		return nil, fmt.Errorf("翻译请求为空")
	}
	started := time.Now()
	if args.TextContent == "" {
		return &translator.TranslateRes{Results: []*translator.TranslateResBlock{{Id: "", TextTranslated: ""}}}, nil
	}
	parts := splitText(args.TextContent, t.cfg.GetMaxCharNum())
	translated := make([]string, 0, len(parts))
	for _, part := range parts {
		if err := restrictor.Singleton().Wait(t.GetId(), ctx); err != nil {
			return nil, fmt.Errorf("腾讯云翻译等待限流失败: %w", err)
		}
		result, err := t.translatePart(ctx, args.FromLang, args.ToLang, part)
		if err != nil {
			return nil, err
		}
		translated = append(translated, result)
	}
	return &translator.TranslateRes{
		TimeUsed: int(time.Since(started).Seconds()),
		Results: []*translator.TranslateResBlock{{
			Id:             args.TextContent,
			TextTranslated: strings.Join(translated, ""),
		}},
	}, nil
}

func splitText(value string, configuredLimit int) []string {
	limit := configuredLimit
	if limit < 1 || limit > 2000 {
		limit = 2000
	}
	runes := []rune(value)
	if len(runes) == 0 {
		return []string{""}
	}
	parts := make([]string, 0, (len(runes)+limit-1)/limit)
	for len(runes) > 0 {
		count := len(runes)
		if count > limit {
			count = limit
		}
		parts = append(parts, string(runes[:count]))
		runes = runes[count:]
	}
	return parts
}

func (t *Translator) translatePart(ctx context.Context, from, to, text string) (string, error) {
	body, err := json.Marshal(struct {
		SourceText string `json:"SourceText"`
		Source     string `json:"Source"`
		Target     string `json:"Target"`
		ProjectID  int    `json:"ProjectId"`
	}{SourceText: text, Source: from, Target: to, ProjectID: 0})
	if err != nil {
		return "", fmt.Errorf("编码腾讯云翻译请求失败: %w", err)
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, t.endpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("创建腾讯云翻译请求失败: %w", err)
	}
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Set("X-TC-Action", "TextTranslate")
	request.Header.Set("X-TC-Version", "2018-03-21")
	request.Header.Set("X-TC-Region", t.cfg.GetRegion())
	requestTime := time.Now()
	request.Header.Set("X-TC-Timestamp", fmt.Sprint(requestTime.Unix()))
	request.Header.Set("Authorization", authorization(t.cfg.GetAK(), t.cfg.GetSK(), request.Host, body, requestTime))
	response, err := t.client.Do(request)
	if err != nil {
		return "", fmt.Errorf("请求腾讯云翻译服务失败: %w", err)
	}
	defer response.Body.Close()
	var envelope struct {
		Response struct {
			TargetText *string `json:"TargetText"`
			RequestID  string  `json:"RequestId"`
			Error      *struct {
				Code    string `json:"Code"`
				Message string `json:"Message"`
			} `json:"Error"`
		} `json:"Response"`
	}
	if err := json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&envelope); err != nil {
		return "", fmt.Errorf("解析腾讯云翻译响应失败（HTTP %d）: %w", response.StatusCode, err)
	}
	if envelope.Response.Error != nil {
		return "", fmt.Errorf("腾讯云翻译失败（%s，RequestId: %s）: %s", envelope.Response.Error.Code, envelope.Response.RequestID, envelope.Response.Error.Message)
	}
	if response.StatusCode != http.StatusOK || envelope.Response.TargetText == nil {
		return "", fmt.Errorf("腾讯云翻译响应无效（HTTP %d，RequestId: %s）", response.StatusCode, envelope.Response.RequestID)
	}
	return *envelope.Response.TargetText, nil
}
