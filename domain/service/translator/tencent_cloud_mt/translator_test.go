package tencent_cloud_mt

import (
	"anto/domain/service/translator"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func configuredTranslator(endpoint string) *Translator {
	t := New()
	t.Init(&Config{SecretID: "test-id", SecretKey: "test-secret", Region: "ap-guangzhou", QPS: 5, MaxCharNum: 2000, MaxCoroutineNum: 5})
	t.endpoint = endpoint
	return t
}

func TestTranslateSendsTencentTextRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/" {
			t.Errorf("request route = %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("X-TC-Action") != "TextTranslate" || r.Header.Get("X-TC-Version") != "2018-03-21" || r.Header.Get("X-TC-Region") != "ap-guangzhou" {
			t.Errorf("Tencent headers are missing: %+v", r.Header)
		}
		if !strings.HasPrefix(r.Header.Get("Authorization"), "TC3-HMAC-SHA256 Credential=test-id/") {
			t.Errorf("authorization header is missing TC3 signature")
		}
		var payload struct {
			SourceText string `json:"SourceText"`
			Source     string `json:"Source"`
			Target     string `json:"Target"`
			ProjectID  int    `json:"ProjectId"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("decode request: %v", err)
		}
		if payload.SourceText != "Hello" || payload.Source != "en" || payload.Target != "zh" || payload.ProjectID != 0 {
			t.Errorf("unexpected payload: %+v", payload)
		}
		_, _ = w.Write([]byte(`{"Response":{"TargetText":"你好","RequestId":"request-1"}}`))
	}))
	defer server.Close()
	engine := configuredTranslator(server.URL)
	result, err := engine.Translate(context.Background(), &translator.TranslateArgs{FromLang: "en", ToLang: "zh", TextContent: "Hello"})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Results) != 1 || result.Results[0].Id != "Hello" || result.Results[0].TextTranslated != "你好" {
		t.Fatalf("unexpected translation: %+v", result.Results)
	}
}

func TestTranslateReportsProviderErrorWithoutExposingKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"Response":{"Error":{"Code":"FailedOperation.UserNotRegistered","Message":"service unavailable"},"RequestId":"request-2"}}`))
	}))
	defer server.Close()
	engine := configuredTranslator(server.URL)
	_, err := engine.Translate(context.Background(), &translator.TranslateArgs{FromLang: "en", ToLang: "zh", TextContent: "Hello"})
	if err == nil || !strings.Contains(err.Error(), "FailedOperation.UserNotRegistered") || !strings.Contains(err.Error(), "request-2") {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(err.Error(), "test-secret") {
		t.Fatal("provider secret leaked into error")
	}
}

func TestTranslationRespectsContextAndCharacterLimit(t *testing.T) {
	parts := splitText(strings.Repeat("界", 2001), 2000)
	if len(parts) != 2 || len([]rune(parts[0])) != 2000 || len([]rune(parts[1])) != 1 {
		t.Fatalf("unexpected chunks: %d", len(parts))
	}
	engine := configuredTranslator("http://127.0.0.1:1")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := engine.Translate(ctx, &translator.TranslateArgs{FromLang: "en", ToLang: "zh", TextContent: "Hello"})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation was lost: %v", err)
	}
	_, err = engine.TestConnection(context.Background(), "invalid")
	if !errors.Is(err, translator.ErrInvalidConnectionRegion) {
		t.Fatalf("unexpected region error: %v", err)
	}
}

func TestSignatureUsesUTCDay(t *testing.T) {
	at := time.Date(2026, 10, 20, 0, 30, 0, 0, time.FixedZone("CST", 8*3600))
	got := authorization("test-id", "test-secret", publicEndpoint, []byte(`{"SourceText":"Hi"}`), at)
	if !strings.Contains(got, "Credential=test-id/2026-10-19/tmt/tc3_request") || !strings.Contains(got, "SignedHeaders=content-type;host") {
		t.Fatalf("unexpected signature scope: %s", got)
	}
	if !strings.HasSuffix(got, "Signature=8787d1c5f5dc98141742117e62910fa00ddc392cb3754ad1d12054397b279535") {
		t.Fatalf("unexpected TC3 signature: %s", got)
	}
}
