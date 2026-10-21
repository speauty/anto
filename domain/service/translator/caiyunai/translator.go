package caiyunai

import (
	"anto/domain/service/translator"
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

	"golang.org/x/time/rate"
)

const publicEndpoint = "api.interpreter.caiyunai.com"
const translateURL = "https://" + publicEndpoint + "/v1/translator"

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
	limiter  *rate.Limiter
}

func New() *Translator {
	return &Translator{client: &http.Client{Timeout: 15 * time.Second}, endpoint: translateURL}
}

func (t *Translator) Init(config translator.ImplConfig) {
	t.cfg = config
	qps := config.GetQPS()
	if qps < 1 {
		qps = 1
	}
	t.limiter = rate.NewLimiter(rate.Limit(qps), 1)
}
func (t *Translator) GetId() string                           { return "caiyunai" }
func (t *Translator) GetShortId() string                      { return "cy" }
func (t *Translator) GetName() string                         { return "彩云小译" }
func (t *Translator) GetCfg() translator.ImplConfig           { return t.cfg }
func (t *Translator) GetLangSupported() []translator.LangPair { return langSupported }
func (t *Translator) GetSep() string                          { return "" }
func (t *Translator) IsValid() bool {
	return t.cfg != nil && strings.TrimSpace(t.cfg.GetAK()) != ""
}

// TestConnection checks reachability without sending a token or spending quota.
func (t *Translator) TestConnection(ctx context.Context, region string) (*translator.ConnectionResult, error) {
	if region != "" {
		return nil, fmt.Errorf("%w: 彩云小译无需地区参数", translator.ErrInvalidConnectionRegion)
	}
	ctx, cancel := context.WithTimeout(ctx, 7*time.Second)
	defer cancel()
	started := time.Now()
	dialer := &tls.Dialer{NetDialer: &net.Dialer{Timeout: 5 * time.Second}, Config: &tls.Config{ServerName: publicEndpoint, MinVersion: tls.VersionTLS12}}
	connection, err := dialer.DialContext(ctx, "tcp", net.JoinHostPort(publicEndpoint, "443"))
	if err != nil {
		return nil, fmt.Errorf("连接 %s 失败: %w", publicEndpoint, err)
	}
	_ = connection.Close()
	return &translator.ConnectionResult{Endpoint: publicEndpoint, LatencyMs: time.Since(started).Milliseconds()}, nil
}

func (t *Translator) Translate(ctx context.Context, args *translator.TranslateArgs) (*translator.TranslateRes, error) {
	if !t.IsValid() {
		return nil, fmt.Errorf("彩云小译尚未配置 API Token")
	}
	if args == nil {
		return nil, fmt.Errorf("翻译请求为空")
	}
	if !SupportsDirection(args.FromLang, args.ToLang) {
		return nil, fmt.Errorf("彩云小译不支持 %s 到 %s 的翻译方向", args.FromLang, args.ToLang)
	}
	started := time.Now()
	if args.TextContent == "" {
		return &translator.TranslateRes{Results: []*translator.TranslateResBlock{{Id: "", TextTranslated: ""}}}, nil
	}
	limit := t.cfg.GetMaxCharNum()
	if limit < 1 {
		limit = 1000
	}
	parts := splitText(args.TextContent, limit)
	translated := make([]string, 0, len(parts))
	for _, part := range parts {
		if err := t.limiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("彩云小译请求等待失败: %w", err)
		}
		result, err := t.translatePart(ctx, args.FromLang, args.ToLang, part)
		if err != nil {
			return nil, err
		}
		translated = append(translated, result)
	}
	return &translator.TranslateRes{TimeUsed: int(time.Since(started).Seconds()), Results: []*translator.TranslateResBlock{{Id: args.TextContent, TextTranslated: strings.Join(translated, "")}}}, nil
}

func splitText(value string, limit int) []string {
	runes := []rune(value)
	parts := make([]string, 0, (len(runes)+limit-1)/limit)
	for start := 0; start < len(runes); start += limit {
		end := start + limit
		if end > len(runes) {
			end = len(runes)
		}
		parts = append(parts, string(runes[start:end]))
	}
	return parts
}

func (t *Translator) translatePart(ctx context.Context, from, to, value string) (string, error) {
	body, err := json.Marshal(struct {
		Source    []string `json:"source"`
		TransType string   `json:"trans_type"`
		Detect    bool     `json:"detect"`
		Media     string   `json:"media"`
	}{Source: []string{value}, TransType: from + "2" + to, Detect: false, Media: "text"})
	if err != nil {
		return "", fmt.Errorf("生成彩云小译请求失败: %w", err)
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, t.endpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("创建彩云小译请求失败: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Authorization", "token "+t.cfg.GetAK())
	response, err := t.client.Do(request)
	if err != nil {
		return "", fmt.Errorf("请求彩云小译失败: %w", err)
	}
	defer response.Body.Close()
	var payload struct {
		Target  []string `json:"target"`
		Message string   `json:"message"`
		Error   string   `json:"error"`
	}
	if err := json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&payload); err != nil {
		return "", fmt.Errorf("解析彩云小译响应失败（HTTP %d）: %w", response.StatusCode, err)
	}
	if response.StatusCode != http.StatusOK {
		message := payload.Message
		if message == "" {
			message = payload.Error
		}
		message = strings.ReplaceAll(message, t.cfg.GetAK(), "********")
		if message == "" {
			message = "请检查 Token、额度或服务状态"
		}
		return "", fmt.Errorf("彩云小译翻译失败（HTTP %d）: %s", response.StatusCode, message)
	}
	if len(payload.Target) != 1 {
		return "", fmt.Errorf("彩云小译响应缺少翻译结果")
	}
	return payload.Target[0], nil
}
