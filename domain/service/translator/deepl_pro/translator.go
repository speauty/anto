package deepl_pro

import (
	"anto/domain/service/translator"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const publicEndpoint = "api.deepl.com"
const translateURL = "https://" + publicEndpoint + "/v2/translate"
const maxRequestBytes = 128 * 1024
const maxTextJSONBytes = maxRequestBytes - 512

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
	return &Translator{client: &http.Client{Timeout: 20 * time.Second}, endpoint: translateURL}
}

func (t *Translator) Init(config translator.ImplConfig) {
	t.cfg = config
	qps := config.GetQPS()
	if qps < 1 {
		qps = 1
	}
	t.limiter = rate.NewLimiter(rate.Limit(qps), 1)
}
func (t *Translator) GetId() string                           { return "deepl_pro" }
func (t *Translator) GetShortId() string                      { return "dlp" }
func (t *Translator) GetName() string                         { return "DeepL Pro" }
func (t *Translator) GetCfg() translator.ImplConfig           { return t.cfg }
func (t *Translator) GetLangSupported() []translator.LangPair { return langSupported }
func (t *Translator) GetSep() string                          { return "" }
func (t *Translator) IsValid() bool {
	return t.cfg != nil && strings.TrimSpace(t.cfg.GetAK()) != ""
}

// TestConnection checks the Pro host without sending credentials or using quota.
func (t *Translator) TestConnection(ctx context.Context, region string) (*translator.ConnectionResult, error) {
	if region != "" {
		return nil, fmt.Errorf("%w: DeepL Pro 无需地区参数", translator.ErrInvalidConnectionRegion)
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
		return nil, fmt.Errorf("DeepL Pro 尚未配置 API Key")
	}
	if args == nil {
		return nil, fmt.Errorf("翻译请求为空")
	}
	if !SupportsDirection(args.FromLang, args.ToLang) {
		return nil, fmt.Errorf("DeepL Pro 不支持 %s 到 %s 的翻译方向", args.FromLang, args.ToLang)
	}
	from, _ := languageCode(args.FromLang, false)
	to, _ := languageCode(args.ToLang, true)
	started := time.Now()
	if args.TextContent == "" {
		return &translator.TranslateRes{Results: []*translator.TranslateResBlock{{Id: "", TextTranslated: ""}}}, nil
	}
	limit := t.cfg.GetMaxCharNum()
	if limit < 1 {
		limit = 3000
	}
	parts := splitText(args.TextContent, limit)
	translated := make([]string, 0, len(parts))
	for _, part := range parts {
		value, err := t.translatePart(ctx, from, to, part)
		if err != nil {
			return nil, err
		}
		translated = append(translated, value)
	}
	return &translator.TranslateRes{TimeUsed: int(time.Since(started).Seconds()), Results: []*translator.TranslateResBlock{{Id: args.TextContent, TextTranslated: strings.Join(translated, "")}}}, nil
}

// splitText respects both the local character setting and DeepL's 128 KiB
// request-body limit, including JSON escapes for control characters.
func splitText(value string, limit int) []string {
	if limit < 1 {
		limit = 3000
	}
	parts := make([]string, 0, len(value)/limit+1)
	var current strings.Builder
	chars, encoded := 0, 0
	for _, letter := range value {
		letterJSON, _ := json.Marshal(string(letter))
		letterBytes := len(letterJSON) - 2
		if chars > 0 && (chars >= limit || encoded+letterBytes > maxTextJSONBytes) {
			parts = append(parts, current.String())
			current.Reset()
			chars, encoded = 0, 0
		}
		current.WriteRune(letter)
		chars++
		encoded += letterBytes
	}
	if current.Len() > 0 {
		parts = append(parts, current.String())
	}
	return parts
}

func (t *Translator) translatePart(ctx context.Context, from, to, value string) (string, error) {
	body, err := json.Marshal(struct {
		Text           []string `json:"text"`
		SourceLang     string   `json:"source_lang"`
		TargetLang     string   `json:"target_lang"`
		SplitSentences string   `json:"split_sentences"`
		PreserveFormat bool     `json:"preserve_formatting"`
	}{Text: []string{value}, SourceLang: from, TargetLang: to, SplitSentences: "0", PreserveFormat: true})
	if err != nil {
		return "", fmt.Errorf("生成 DeepL Pro 请求失败: %w", err)
	}
	if len(body) > maxRequestBytes {
		return "", fmt.Errorf("DeepL Pro 请求超过 128 KiB 限制")
	}
	for attempt := 0; attempt < 3; attempt++ {
		if err := t.limiter.Wait(ctx); err != nil {
			return "", fmt.Errorf("DeepL Pro 请求等待失败: %w", err)
		}
		result, retry, delay, err := t.send(ctx, body)
		if err == nil {
			return result, nil
		}
		if !retry || attempt == 2 {
			return "", err
		}
		if delay <= 0 {
			delay = time.Duration(500*(1<<attempt)+rand.Intn(200)) * time.Millisecond
		}
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return "", fmt.Errorf("DeepL Pro 请求已取消: %w", ctx.Err())
		case <-timer.C:
		}
	}
	return "", fmt.Errorf("DeepL Pro 请求失败")
}

func (t *Translator) send(ctx context.Context, body []byte) (string, bool, time.Duration, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, t.endpoint, bytes.NewReader(body))
	if err != nil {
		return "", false, 0, fmt.Errorf("创建 DeepL Pro 请求失败: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "DeepL-Auth-Key "+t.cfg.GetAK())
	response, err := t.client.Do(request)
	if err != nil {
		return "", false, 0, fmt.Errorf("请求 DeepL Pro 失败: %w", err)
	}
	defer response.Body.Close()
	var payload struct {
		Translations []struct {
			Text string `json:"text"`
		} `json:"translations"`
		Message string `json:"message"`
		Error   struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	decodeErr := json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&payload)
	if response.StatusCode != http.StatusOK {
		message := payload.Message
		if message == "" {
			message = payload.Error.Message
		}
		if message == "" {
			message = http.StatusText(response.StatusCode)
		}
		message = strings.ReplaceAll(message, t.cfg.GetAK(), "********")
		status := response.StatusCode
		retry := status == 429 || status == 529 || status == 500 || status == 503 || status == 504
		return "", retry, retryAfter(response.Header.Get("Retry-After")), fmt.Errorf("DeepL Pro 翻译失败（HTTP %d）: %s", status, message)
	}
	if decodeErr != nil {
		return "", false, 0, fmt.Errorf("解析 DeepL Pro 响应失败: %w", decodeErr)
	}
	if len(payload.Translations) != 1 {
		return "", false, 0, fmt.Errorf("DeepL Pro 响应缺少翻译结果")
	}
	return payload.Translations[0].Text, false, 0, nil
}

func retryAfter(value string) time.Duration {
	if seconds, err := strconv.Atoi(strings.TrimSpace(value)); err == nil && seconds > 0 {
		delay := time.Duration(seconds) * time.Second
		if delay > 30*time.Second {
			return 30 * time.Second
		}
		return delay
	}
	if date, err := http.ParseTime(value); err == nil {
		delay := time.Until(date)
		if delay < 0 {
			return 0
		}
		if delay > 30*time.Second {
			return 30 * time.Second
		}
		return delay
	}
	return 0
}
