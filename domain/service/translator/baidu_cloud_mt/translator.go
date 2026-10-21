package baidu_cloud_mt

import (
	"anto/domain/service/translator"
	"context"
	"crypto/md5"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const publicEndpoint = "fanyi-api.baidu.com"
const translatePath = "/api/trans/vip/translate"

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
	return &Translator{
		client:   &http.Client{Timeout: 15 * time.Second},
		endpoint: "https://" + publicEndpoint + translatePath,
	}
}

func (t *Translator) Init(config translator.ImplConfig) {
	t.cfg = config
	t.limiter = rate.NewLimiter(rate.Limit(config.GetQPS()), 1)
}
func (t *Translator) GetId() string                 { return "baidu_cloud_mt" }
func (t *Translator) GetShortId() string            { return "bd" }
func (t *Translator) GetName() string               { return "机器翻译-百度云" }
func (t *Translator) GetCfg() translator.ImplConfig { return t.cfg }
func (t *Translator) GetLangSupported() []translator.LangPair {
	return langSupported
}

// One request per subtitle preserves the association between source and
// translated blocks even when the API splits a multiline query internally.
func (t *Translator) GetSep() string { return "" }
func (t *Translator) IsValid() bool {
	if t.cfg == nil || strings.TrimSpace(t.cfg.GetAK()) == "" ||
		strings.TrimSpace(t.cfg.GetSK()) == "" {
		return false
	}
	if config, ok := t.cfg.(*Config); ok {
		_, _, valid := TierLimits(config.ServiceTier)
		return valid
	}
	return true
}

// TestConnection checks the public API's TLS connection without consuming quota.
func (t *Translator) TestConnection(ctx context.Context, region string) (*translator.ConnectionResult, error) {
	if region != "" {
		return nil, fmt.Errorf("%w: 百度云无需地区参数", translator.ErrInvalidConnectionRegion)
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
		return nil, fmt.Errorf("百度云翻译引擎未配置 App ID 或密钥")
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
		if err := t.limiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("百度云翻译等待限流失败: %w", err)
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
	if limit < 1 || limit > 6000 {
		limit = 1000
	}
	parts := make([]string, 0, len(value)/1000+1)
	var current strings.Builder
	chars := 0
	for _, letter := range value {
		encoded := string(letter)
		if chars >= limit || current.Len()+len(encoded) > 6000 {
			parts = append(parts, current.String())
			current.Reset()
			chars = 0
		}
		current.WriteString(encoded)
		chars++
	}
	if current.Len() > 0 {
		parts = append(parts, current.String())
	}
	return parts
}

func sign(appID, query, salt, secret string) string {
	sum := md5.Sum([]byte(appID + query + salt + secret))
	return hex.EncodeToString(sum[:])
}

func newSalt() (string, error) {
	value := make([]byte, 16)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return hex.EncodeToString(value), nil
}

func (t *Translator) translatePart(ctx context.Context, from, to, query string) (string, error) {
	salt, err := newSalt()
	if err != nil {
		return "", fmt.Errorf("生成百度云翻译随机数失败: %w", err)
	}
	form := url.Values{
		"q":     {query},
		"from":  {from},
		"to":    {to},
		"appid": {t.cfg.GetAK()},
		"salt":  {salt},
		"sign":  {sign(t.cfg.GetAK(), query, salt, t.cfg.GetSK())},
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, t.endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("创建百度云翻译请求失败: %w", err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := t.client.Do(request)
	if err != nil {
		return "", fmt.Errorf("请求百度云翻译服务失败: %w", err)
	}
	defer response.Body.Close()
	var payload struct {
		ErrorCode string `json:"error_code"`
		ErrorMsg  string `json:"error_msg"`
		Results   []struct {
			Dst string `json:"dst"`
		} `json:"trans_result"`
	}
	if err := json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&payload); err != nil {
		return "", fmt.Errorf("解析百度云翻译响应失败（HTTP %d）: %w", response.StatusCode, err)
	}
	if payload.ErrorCode != "" {
		message := strings.ReplaceAll(payload.ErrorMsg, t.cfg.GetSK(), "********")
		return "", fmt.Errorf("百度云翻译失败（%s）: %s", payload.ErrorCode, message)
	}
	if response.StatusCode != http.StatusOK || len(payload.Results) == 0 {
		return "", fmt.Errorf("百度云翻译响应无效（HTTP %d）", response.StatusCode)
	}
	results := make([]string, 0, len(payload.Results))
	for _, item := range payload.Results {
		results = append(results, item.Dst)
	}
	return strings.Join(results, "\n"), nil
}
