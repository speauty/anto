package deepseek

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestTranslateBatchPreservesPositionsAndConfig(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.URL.Path != "/v1/chat/completions" || r.Method != http.MethodPost {
			t.Errorf("unexpected endpoint: %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer local-key" {
			t.Errorf("authorization = %q", r.Header.Get("Authorization"))
		}
		var request struct {
			Model    string `json:"model"`
			Messages []struct {
				Content string `json:"content"`
			} `json:"messages"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Error(err)
			return
		}
		if request.Model != "my-deepseek" || len(request.Messages) != 2 || !strings.Contains(request.Messages[1].Content, "自然、简洁") {
			t.Errorf("request settings were lost: %+v", request)
		}
		text := strings.TrimPrefix(request.Messages[1].Content, "翻译风格要求：自然、简洁\n字幕 JSON：")
		var input []string
		if err := json.Unmarshal([]byte(text), &input); err != nil {
			t.Error(err)
			return
		}
		values := make([]string, len(input))
		for i, item := range input {
			values[i] = "[" + item + "]"
		}
		content, _ := json.Marshal(map[string]interface{}{"translations": values})
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"choices": []interface{}{map[string]interface{}{"finish_reason": "stop", "message": map[string]string{"content": string(content)}}}})
	}))
	defer server.Close()
	config := new(Config).Default().(*Config)
	config.BaseURL = server.URL + "/v1"
	config.APIKey = "local-key"
	config.Model = "my-deepseek"
	config.Prompt = "自然、简洁"
	config.MaxCharNum = 4
	engine := New()
	engine.Init(config)
	got, err := engine.TranslateBatch(context.Background(), "en", "zh", []string{"same", "", "same", "abcdef"})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"[same]", "", "[same]", "[abcd][ef]"}
	if !reflect.DeepEqual(got, want) || requests != 4 {
		t.Fatalf("translated = %v, requests = %d; want %v and 4", got, requests, want)
	}
}

func TestTranslateBatchRejectsMismatchedOutput(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"choices":[{"finish_reason":"stop","message":{"content":"{\"translations\":[\"one\"]}"}}]}`)
	}))
	defer server.Close()
	config := new(Config).Default().(*Config)
	config.BaseURL = server.URL
	engine := New()
	engine.Init(config)
	_, err := engine.TranslateBatch(context.Background(), "en", "zh", []string{"a", "b"})
	if err == nil || !strings.Contains(err.Error(), "条数不匹配") {
		t.Fatalf("expected count mismatch, got %v", err)
	}
}

func TestValidateBaseURL(t *testing.T) {
	for _, valid := range []string{"https://api.deepseek.com", "https://example.com/v1", "http://127.0.0.1:8000/v1", "http://localhost:11434/v1"} {
		if err := ValidateBaseURL(valid); err != nil {
			t.Errorf("valid base URL %q rejected: %v", valid, err)
		}
	}
	for _, invalid := range []string{"", "http://example.com/v1", "https://name:secret@example.com", "https://example.com/v1?key=secret"} {
		if err := ValidateBaseURL(invalid); err == nil {
			t.Errorf("invalid base URL %q accepted", invalid)
		}
	}
}
