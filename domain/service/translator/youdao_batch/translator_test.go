package youdao_batch

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestBatchRequestSignsAllTextsAndPreservesPosition(t *testing.T) {
	want := []string{"你好\n世界", "重复", "重复"}
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.Method != http.MethodPost || r.URL.Path != "/v2/api" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(r.PostForm["q"], want) {
			t.Errorf("q fields = %q, want %q", r.PostForm["q"], want)
		}
		if r.PostForm.Get("from") != "en" || r.PostForm.Get("to") != "zh-CHS" || r.PostForm.Get("signType") != "v3" {
			t.Errorf("unexpected language/sign type fields: %+v", r.PostForm)
		}
		joined := []rune(strings.Join(want, ""))
		input := string(joined)
		if len(joined) > 20 {
			input = string(joined[:10]) + strconv.Itoa(len(joined)) + string(joined[len(joined)-10:])
		}
		digest := sha256.Sum256([]byte("test-app" + input + r.PostForm.Get("salt") + r.PostForm.Get("curtime") + "test-secret"))
		if r.PostForm.Get("sign") != hex.EncodeToString(digest[:]) {
			t.Error("invalid v3 signature")
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"errorCode": "0",
			"translateResults": []map[string]string{
				{"query": want[0], "translation": "first"},
				{"query": want[1], "translation": "second"},
				{"query": want[2], "translation": "third"},
			},
		})
	}))
	defer server.Close()
	engine := New()
	engine.endpoint = server.URL + "/v2/api"
	engine.Init(&Config{AppKey: "test-app", AppSecret: "test-secret", QPS: 100, MaxCharNum: 5000})
	got, err := engine.TranslateBatch(context.Background(), "en", "zh", want)
	if err != nil {
		t.Fatal(err)
	}
	if calls != 1 || !reflect.DeepEqual(got, []string{"first", "second", "third"}) {
		t.Fatalf("calls=%d results=%q", calls, got)
	}
}

func TestBatchPartialFailureDoesNotWriteMisalignedResults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"errorCode":"0","errorIndex":[1],"translateResults":[{"query":"one","translation":"一"}]}`))
	}))
	defer server.Close()
	engine := New()
	engine.endpoint = server.URL
	engine.Init(&Config{AppKey: "test-app", AppSecret: "test-secret", QPS: 100, MaxCharNum: 5000})
	result, err := engine.TranslateBatch(context.Background(), "en", "zh", []string{"one", "two"})
	if err == nil || result != nil || !strings.Contains(err.Error(), "失败序号 [1]") {
		t.Fatalf("partial failure should abort batch: result=%q err=%v", result, err)
	}
}

func TestSignatureTruncatesByUnicodeCodePoints(t *testing.T) {
	input := signatureInput([]string{"你好世界测试汉字一二三四五六七八九十甲乙丙丁戊"})
	if input != "你好世界测试汉字一二23六七八九十甲乙丙丁戊" {
		t.Fatalf("unexpected truncated input: %q", input)
	}
}

func TestLongTextIsSplitWithinLocalLimit(t *testing.T) {
	var calls [][]string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		queries := append([]string(nil), r.PostForm["q"]...)
		calls = append(calls, queries)
		results := make([]map[string]string, 0, len(queries))
		for _, query := range queries {
			results = append(results, map[string]string{"query": query, "translation": "X"})
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"errorCode": "0", "translateResults": results})
	}))
	defer server.Close()
	engine := New()
	engine.endpoint = server.URL
	engine.Init(&Config{AppKey: "test-app", AppSecret: "test-secret", QPS: 100, MaxCharNum: 5000})
	got, err := engine.TranslateBatch(context.Background(), "en", "zh", []string{strings.Repeat("界", 5001)})
	if err != nil {
		t.Fatal(err)
	}
	if len(calls) != 2 || len([]rune(calls[0][0])) != 5000 || calls[1][0] != "界" || !reflect.DeepEqual(got, []string{"XX"}) {
		t.Fatalf("unexpected split calls=%d result=%q", len(calls), got)
	}
}
