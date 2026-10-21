package deepl_pro

import (
	"anto/domain/service/translator"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func testProvider(endpoint string) *Translator {
	provider := New()
	provider.endpoint = endpoint
	provider.Init(&Config{APIKey: "test-key", QPS: 1000, MaxCharNum: 2, MaxCoroutineNum: 1})
	return provider
}

func TestTranslateUsesProRequestAndPreservesSubtitleBlock(t *testing.T) {
	var sources []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v2/translate" {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("Authorization") != "DeepL-Auth-Key test-key" || r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("unexpected headers: %v", r.Header)
		}
		var request struct {
			Text               []string `json:"text"`
			SourceLang         string   `json:"source_lang"`
			TargetLang         string   `json:"target_lang"`
			SplitSentences     string   `json:"split_sentences"`
			PreserveFormatting bool     `json:"preserve_formatting"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("decode: %v", err)
		}
		if request.SourceLang != "EN" || request.TargetLang != "ZH" || request.SplitSentences != "0" || !request.PreserveFormatting || len(request.Text) != 1 {
			t.Errorf("unexpected payload: %+v", request)
		}
		sources = append(sources, request.Text[0])
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"translations": []map[string]string{{"text": "译"}}})
	}))
	defer server.Close()
	provider := testProvider(server.URL + "/v2/translate")
	result, err := provider.Translate(context.Background(), &translator.TranslateArgs{FromLang: "en", ToLang: "zh", TextContent: "a中b"})
	if err != nil {
		t.Fatal(err)
	}
	if len(sources) != 2 || sources[0] != "a中" || sources[1] != "b" {
		t.Fatalf("split lost rune boundaries: %q", sources)
	}
	if result.Results[0].Id != "a中b" || result.Results[0].TextTranslated != "译译" {
		t.Fatalf("unexpected result: %+v", result.Results[0])
	}
}

func TestLanguageCodesAndRequestSize(t *testing.T) {
	if code, ok := languageCode("en", true); !ok || code != "EN-US" {
		t.Fatalf("English target: %q, %v", code, ok)
	}
	if code, ok := languageCode("pt", true); !ok || code != "PT-PT" {
		t.Fatalf("Portuguese target: %q, %v", code, ok)
	}
	if !SupportsDirection("en", "zh") || SupportsDirection("en", "en") || SupportsDirection("unknown", "zh") {
		t.Fatal("direction validation failed")
	}
	parts := splitText(strings.Repeat("\n", 80000), 100000)
	if len(parts) < 2 {
		t.Fatalf("oversized request was not split: %d", len(parts))
	}
	for _, part := range parts {
		body, err := json.Marshal(map[string]interface{}{"text": []string{part}, "source_lang": "EN", "target_lang": "ZH", "split_sentences": "0", "preserve_formatting": true})
		if err != nil || len(body) > maxRequestBytes {
			t.Fatalf("body too large: %d, %v", len(body), err)
		}
	}
}

func TestRetryableAndQuotaErrors(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) == 1 {
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"error":{"message":"busy test-key"}}`))
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"translations": []map[string]string{{"text": "你好"}}})
	}))
	defer server.Close()
	provider := testProvider(server.URL)
	result, err := provider.Translate(context.Background(), &translator.TranslateArgs{FromLang: "en", ToLang: "zh", TextContent: "hi"})
	if err != nil || result.Results[0].TextTranslated != "你好" || calls.Load() != 2 {
		t.Fatalf("retry failed: %+v, %v, %d", result, err, calls.Load())
	}

	var quotaCalls atomic.Int32
	quota := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		quotaCalls.Add(1)
		w.WriteHeader(456)
		_, _ = w.Write([]byte(`{"message":"quota reached test-key"}`))
	}))
	defer quota.Close()
	provider.endpoint = quota.URL
	_, err = provider.Translate(context.Background(), &translator.TranslateArgs{FromLang: "en", ToLang: "zh", TextContent: "hi"})
	if err == nil || strings.Contains(err.Error(), "test-key") || quotaCalls.Load() != 1 || !strings.Contains(err.Error(), "456") {
		t.Fatalf("quota handling failed: %v, calls=%d", err, quotaCalls.Load())
	}
}

func TestRetryStopsOnCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "30")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"message":"busy"}`))
	}))
	defer server.Close()
	provider := testProvider(server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	started := time.Now()
	_, err := provider.Translate(ctx, &translator.TranslateArgs{FromLang: "en", ToLang: "zh", TextContent: "hi"})
	if err == nil || time.Since(started) > time.Second {
		t.Fatalf("cancellation was not prompt: %v", err)
	}
}
