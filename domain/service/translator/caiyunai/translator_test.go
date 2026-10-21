package caiyunai

import (
	"anto/domain/service/translator"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTranslateUsesDocumentedRequestAndPreservesBlocks(t *testing.T) {
	var sources []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/translator" {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("X-Authorization") != "token test-token" || r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("unexpected headers: %v", r.Header)
		}
		var request struct {
			Source    []string `json:"source"`
			TransType string   `json:"trans_type"`
			Detect    bool     `json:"detect"`
			Media     string   `json:"media"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("decode: %v", err)
		}
		if request.TransType != "en2zh" || request.Detect || request.Media != "text" || len(request.Source) != 1 {
			t.Errorf("unexpected payload: %+v", request)
		}
		sources = append(sources, request.Source[0])
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"target": []string{"译"}})
	}))
	defer server.Close()
	provider := New()
	provider.endpoint = server.URL + "/v1/translator"
	provider.Init(&Config{Token: "test-token", QPS: 1000, MaxCharNum: 2, MaxCoroutineNum: 1})
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

func TestDirectionsAndErrors(t *testing.T) {
	for _, item := range []struct {
		from, to  string
		supported bool
	}{
		{"en", "zh", true}, {"zh-Hant", "ja", true}, {"de", "zh-Hant", true},
		{"en", "ja", false}, {"de", "en", false}, {"zh", "zh", false},
	} {
		if got := SupportsDirection(item.from, item.to); got != item.supported {
			t.Errorf("%s to %s: %v", item.from, item.to, got)
		}
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"message":"bad test-token"}`))
	}))
	defer server.Close()
	provider := New()
	provider.endpoint = server.URL
	provider.Init(&Config{Token: "test-token", QPS: 1000, MaxCharNum: 1000, MaxCoroutineNum: 1})
	_, err := provider.Translate(context.Background(), &translator.TranslateArgs{FromLang: "en", ToLang: "zh", TextContent: "hello"})
	if err == nil || strings.Contains(err.Error(), "test-token") {
		t.Fatalf("error leaked token or missing: %v", err)
	}
	_, err = provider.Translate(context.Background(), &translator.TranslateArgs{FromLang: "en", ToLang: "ja", TextContent: "hello"})
	if err == nil || !strings.Contains(err.Error(), "不支持") {
		t.Fatalf("invalid direction not rejected: %v", err)
	}
}
