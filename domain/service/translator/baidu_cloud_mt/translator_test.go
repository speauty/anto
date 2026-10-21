package baidu_cloud_mt

import (
	"anto/domain/service/translator"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func configuredTranslator(endpoint, tier string) *Translator {
	engine := New()
	engine.Init(&Config{AppID: "test-app", AppSecret: "test-secret", ServiceTier: tier, MaxCoroutineNum: 1})
	engine.endpoint = endpoint
	return engine
}

func TestTierLimitsAndSignature(t *testing.T) {
	for _, item := range []struct {
		tier  string
		chars int
		qps   int
	}{
		{TierStandard, 1000, 1}, {TierAdvanced, 6000, 10}, {TierPremium, 6000, 100},
	} {
		chars, qps, valid := TierLimits(item.tier)
		if !valid || chars != item.chars || qps != item.qps {
			t.Fatalf("unexpected limits for %s: %d chars, %d QPS", item.tier, chars, qps)
		}
		engine := configuredTranslator("http://127.0.0.1", item.tier)
		if engine.GetCfg().GetMaxCharNum() != chars || int(engine.limiter.Limit()) != qps {
			t.Fatalf("tier %s was not applied to requests", item.tier)
		}
	}
	if got := sign("2015063000000001", "apple", "1435660288", "12345678"); got != "f89f9594663708c1605f3d736d01d2d4" {
		t.Fatalf("official signature example mismatch: %s", got)
	}
}

func TestTranslateUsesFormAndPreservesSubtitle(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost || request.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			t.Errorf("unexpected request: %s %s", request.Method, request.Header.Get("Content-Type"))
		}
		if err := request.ParseForm(); err != nil {
			t.Fatal(err)
		}
		values := request.PostForm
		if values.Get("q") != "Hello\nworld" || values.Get("from") != "en" || values.Get("to") != "zh" || values.Get("appid") != "test-app" {
			t.Errorf("unexpected request fields: %v", values)
		}
		if values.Get("sign") != sign("test-app", values.Get("q"), values.Get("salt"), "test-secret") {
			t.Error("request signature differs from unencoded query")
		}
		_, _ = writer.Write([]byte(`{"trans_result":[{"src":"Hello","dst":"你好"},{"src":"world","dst":"世界"}]}`))
	}))
	defer server.Close()
	engine := configuredTranslator(server.URL, TierStandard)
	result, err := engine.Translate(context.Background(), &translator.TranslateArgs{FromLang: "en", ToLang: "zh", TextContent: "Hello\nworld"})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Results) != 1 || result.Results[0].Id != "Hello\nworld" || result.Results[0].TextTranslated != "你好\n世界" {
		t.Fatalf("unexpected translation: %+v", result.Results)
	}
}

func TestChunksRespectTierCharactersAndUTF8Bytes(t *testing.T) {
	for _, item := range []struct {
		text  string
		tier  string
		parts int
	}{
		{strings.Repeat("a", 1001), TierStandard, 2},
		{strings.Repeat("界", 2001), TierAdvanced, 2},
	} {
		limit, _, _ := TierLimits(item.tier)
		parts := splitText(item.text, limit)
		if len(parts) != item.parts || strings.Join(parts, "") != item.text {
			t.Fatalf("unexpected chunks for %s: %d", item.tier, len(parts))
		}
		for _, part := range parts {
			if len(part) > 6000 || len([]rune(part)) > limit {
				t.Fatalf("chunk exceeded documented limits: %d bytes", len(part))
			}
		}
	}
}

func TestProviderErrorAndCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusForbidden)
		_, _ = fmt.Fprint(writer, `{"error_code":"54001","error_msg":"test-secret invalid"}`)
	}))
	defer server.Close()
	engine := configuredTranslator(server.URL, TierStandard)
	_, err := engine.Translate(context.Background(), &translator.TranslateArgs{FromLang: "en", ToLang: "zh", TextContent: "Hello"})
	if err == nil || !strings.Contains(err.Error(), "54001") || strings.Contains(err.Error(), "test-secret") {
		t.Fatalf("unexpected provider error: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = engine.Translate(ctx, &translator.TranslateArgs{FromLang: "en", ToLang: "zh", TextContent: "Hello"})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation was lost: %v", err)
	}
	_, err = engine.TestConnection(context.Background(), "ap-guangzhou")
	if !errors.Is(err, translator.ErrInvalidConnectionRegion) {
		t.Fatalf("unexpected region error: %v", err)
	}
}
