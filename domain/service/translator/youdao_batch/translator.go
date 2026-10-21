package youdao_batch

import (
	"anto/domain/service/translator"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"golang.org/x/time/rate"
)

const publicEndpoint = "openapi.youdao.com"
const translateURL = "https://" + publicEndpoint + "/v2/api"
const serviceCharLimit = 5000
const maxTextsPerRequest = 50 // A local bound on form size, not a provider quota.

var (
	apiTranslator  *Translator
	onceTranslator sync.Once
)

func API() *Translator {
	onceTranslator.Do(func() { apiTranslator = New() })
	return apiTranslator
}

type Translator struct {
	cfg      *Config
	client   *http.Client
	endpoint string
	limiter  *rate.Limiter
}

func New() *Translator {
	return &Translator{client: &http.Client{Timeout: 25 * time.Second}, endpoint: translateURL}
}

func (t *Translator) Init(config translator.ImplConfig) {
	t.cfg, _ = config.(*Config)
	qps := config.GetQPS()
	if qps < 1 {
		qps = 1
	}
	t.limiter = rate.NewLimiter(rate.Limit(qps), 1)
}
func (t *Translator) GetId() string                           { return "youdao_batch" }
func (t *Translator) GetShortId() string                      { return "ydb" }
func (t *Translator) GetName() string                         { return "批量翻译-有道智云" }
func (t *Translator) GetCfg() translator.ImplConfig           { return t.cfg }
func (t *Translator) GetLangSupported() []translator.LangPair { return langSupported }
func (t *Translator) GetSep() string                          { return "" }
func (t *Translator) IsValid() bool {
	return t.cfg != nil && strings.TrimSpace(t.cfg.AppKey) != "" && strings.TrimSpace(t.cfg.AppSecret) != ""
}

func (t *Translator) TestConnection(ctx context.Context, region string) (*translator.ConnectionResult, error) {
	if region != "" {
		return nil, fmt.Errorf("%w: 有道智云无需地区参数", translator.ErrInvalidConnectionRegion)
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
	if args == nil {
		return nil, fmt.Errorf("翻译请求为空")
	}
	started := time.Now()
	values, err := t.TranslateBatch(ctx, args.FromLang, args.ToLang, []string{args.TextContent})
	if err != nil {
		return nil, err
	}
	return &translator.TranslateRes{TimeUsed: int(time.Since(started).Seconds()), Results: []*translator.TranslateResBlock{{Id: args.TextContent, TextTranslated: values[0]}}}, nil
}

// TranslateBatch returns exactly one result per input, in the same position.
// A long subtitle is split into several queries and joined again locally.
func (t *Translator) TranslateBatch(ctx context.Context, fromLang, toLang string, texts []string) ([]string, error) {
	if !t.IsValid() {
		return nil, fmt.Errorf("有道智云尚未配置应用 ID 和应用密钥")
	}
	if !SupportsDirection(fromLang, toLang) {
		return nil, fmt.Errorf("有道智云不支持 %s 到 %s 的翻译方向", fromLang, toLang)
	}
	from, _ := languageCode(fromLang)
	to, _ := languageCode(toLang)
	limit := t.cfg.MaxCharNum
	if limit < 1 || limit > serviceCharLimit {
		limit = serviceCharLimit
	}
	results := make([]string, len(texts))
	type part struct {
		index int
		text  string
	}
	parts := make([]part, 0, len(texts))
	for i, value := range texts {
		for _, segment := range splitRunes(value, limit) {
			parts = append(parts, part{index: i, text: segment})
		}
	}
	for start := 0; start < len(parts); {
		end, chars := start, 0
		for end < len(parts) && end-start < maxTextsPerRequest {
			length := utf8.RuneCountInString(parts[end].text)
			if end > start && chars+length > limit {
				break
			}
			chars += length
			end++
		}
		queries := make([]string, 0, end-start)
		for _, segment := range parts[start:end] {
			queries = append(queries, segment.text)
		}
		translated, err := t.requestBatch(ctx, from, to, queries)
		if err != nil {
			return nil, err
		}
		for i, value := range translated {
			results[parts[start+i].index] += value
		}
		start = end
	}
	return results, nil
}

func splitRunes(value string, limit int) []string {
	if value == "" {
		return nil
	}
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

func signatureInput(queries []string) string {
	joined := strings.Join(queries, "")
	runes := []rune(joined)
	if len(runes) <= 20 {
		return joined
	}
	return string(runes[:10]) + strconv.Itoa(len(runes)) + string(runes[len(runes)-10:])
}

func (t *Translator) requestBatch(ctx context.Context, from, to string, queries []string) ([]string, error) {
	for attempt := 0; attempt < 3; attempt++ {
		if err := t.limiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("有道智云请求等待失败: %w", err)
		}
		values, err := t.formValues(from, to, queries)
		if err != nil {
			return nil, err
		}
		result, retry, err := t.send(ctx, values, queries)
		if err == nil {
			return result, nil
		}
		if !retry || attempt == 2 {
			return nil, err
		}
		timer := time.NewTimer(time.Duration(500*(1<<attempt)) * time.Millisecond)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil, fmt.Errorf("有道智云请求已取消: %w", ctx.Err())
		case <-timer.C:
		}
	}
	return nil, fmt.Errorf("有道智云翻译失败")
}

func (t *Translator) formValues(from, to string, queries []string) (url.Values, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("生成有道智云请求随机数失败: %w", err)
	}
	saltText := hex.EncodeToString(salt)
	curtime := strconv.FormatInt(time.Now().Unix(), 10)
	input := t.cfg.AppKey + signatureInput(queries) + saltText + curtime + t.cfg.AppSecret
	digest := sha256.Sum256([]byte(input))
	values := url.Values{"from": {from}, "to": {to}, "appKey": {t.cfg.AppKey}, "salt": {saltText}, "curtime": {curtime}, "signType": {"v3"}, "sign": {hex.EncodeToString(digest[:])}}
	for _, query := range queries {
		values.Add("q", query)
	}
	return values, nil
}

func (t *Translator) send(ctx context.Context, values url.Values, queries []string) ([]string, bool, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, t.endpoint, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, false, fmt.Errorf("创建有道智云请求失败: %w", err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	response, err := t.client.Do(request)
	if err != nil {
		return nil, false, fmt.Errorf("请求有道智云失败: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		status := response.StatusCode
		return nil, status == 429 || status >= 500, fmt.Errorf("有道智云翻译失败（HTTP %d）", status)
	}
	var payload struct {
		ErrorCode        json.RawMessage `json:"errorCode"`
		ErrorIndex       []int           `json:"errorIndex"`
		TranslateResults []struct {
			Query       string `json:"query"`
			Translation string `json:"translation"`
		} `json:"translateResults"`
	}
	if err := json.NewDecoder(io.LimitReader(response.Body, 4<<20)).Decode(&payload); err != nil {
		return nil, false, fmt.Errorf("解析有道智云响应失败: %w", err)
	}
	code := strings.Trim(string(payload.ErrorCode), "\"")
	if code != "0" || len(payload.ErrorIndex) != 0 {
		return nil, code == "411" || code == "412", fmt.Errorf("有道智云翻译失败（错误码 %s，失败序号 %v）", code, payload.ErrorIndex)
	}
	if len(payload.TranslateResults) != len(queries) {
		return nil, false, fmt.Errorf("有道智云返回 %d 条译文，预期 %d 条", len(payload.TranslateResults), len(queries))
	}
	result := make([]string, len(queries))
	for i, item := range payload.TranslateResults {
		if item.Query != queries[i] {
			return nil, false, fmt.Errorf("有道智云第 %d 条结果与原文不匹配", i+1)
		}
		result[i] = item.Translation
	}
	return result, false, nil
}
