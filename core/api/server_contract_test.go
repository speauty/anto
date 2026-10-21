package api

import (
	"anto/cfg"
	"anto/core"
	"anto/domain/repository"
	"anto/domain/service/translator"
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestServerCoreEndpointsAndJSONErrors(t *testing.T) {
	server, base, done := startContractServer(t)
	defer stopContractServer(t, server, done)
	client := &http.Client{Timeout: time.Second}

	response, err := client.Get(base + "/api/v1/health")
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatalf("health status = %d", response.StatusCode)
	}
	_ = response.Body.Close()
	request, _ := http.NewRequest(http.MethodPost, base+"/api/v1/health", nil)
	response, err = client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("health method status = %d", response.StatusCode)
	}
	_ = response.Body.Close()

	request, _ = http.NewRequest(http.MethodGet, base+"/api/v1/config", nil)
	request.Header.Set("X-Anto-Token", server.Ready().Token)
	response, err = client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatalf("config status = %d", response.StatusCode)
	}
	var config core.ConfigDTO
	if err := json.NewDecoder(response.Body).Decode(&config); err != nil {
		t.Fatal(err)
	}
	_ = response.Body.Close()
	if config.App.Concurrency == 0 {
		t.Fatal("config response did not include app defaults")
	}

	configPath := filepath.Join(t.TempDir(), "cfg.yml")
	if err := os.WriteFile(configPath, []byte("{}\n"), 0600); err != nil {
		t.Fatal(err)
	}
	cfg.Singleton().GetViper().SetConfigFile(configPath)
	request, _ = http.NewRequest(http.MethodPut, base+"/api/v1/config", strings.NewReader(`{"app":{"default_from_lang":"en","default_to_lang":"zh","concurrency":3}}`))
	request.Header.Set("X-Anto-Token", server.Ready().Token)
	request.Header.Set("Content-Type", "application/json")
	response, err = client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatalf("config update status = %d", response.StatusCode)
	}
	var updated core.ConfigDTO
	if err := json.NewDecoder(response.Body).Decode(&updated); err != nil {
		t.Fatal(err)
	}
	_ = response.Body.Close()
	if updated.App.Concurrency != 3 || updated.App.DefaultFromLang != "en" {
		t.Fatalf("config update was not applied: %+v", updated.App)
	}
	request, _ = http.NewRequest(http.MethodPut, base+"/api/v1/config", strings.NewReader(`{"app":{"output_directory":"output-only","keep_original":false}}`))
	request.Header.Set("X-Anto-Token", server.Ready().Token)
	request.Header.Set("Content-Type", "application/json")
	response, err = client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatalf("partial config update status = %d", response.StatusCode)
	}
	var partial core.ConfigDTO
	if err := json.NewDecoder(response.Body).Decode(&partial); err != nil {
		t.Fatal(err)
	}
	_ = response.Body.Close()
	if partial.App.OutputDirectory != "output-only" || partial.App.KeepOriginal || partial.App.Concurrency != 3 || partial.App.DefaultFromLang != "en" {
		t.Fatalf("partial config update lost existing fields: %+v", partial.App)
	}

	request, _ = http.NewRequest(http.MethodPut, base+"/api/v1/config", strings.NewReader(`{"app":{"concurrency":0}}`))
	request.Header.Set("X-Anto-Token", server.Ready().Token)
	request.Header.Set("Content-Type", "application/json")
	response, err = client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("invalid app config status = %d", response.StatusCode)
	}
	_ = response.Body.Close()

	request, _ = http.NewRequest(http.MethodPut, base+"/api/v1/config", strings.NewReader(`{"providers":{"contract-test":{"qps":0}}}`))
	request.Header.Set("X-Anto-Token", server.Ready().Token)
	request.Header.Set("Content-Type", "application/json")
	response, err = client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("invalid provider config status = %d", response.StatusCode)
	}
	_ = response.Body.Close()

	request, _ = http.NewRequest(http.MethodGet, base+"/api/v1/tasks", nil)
	request.Header.Set("X-Anto-Token", server.Ready().Token)
	response, err = client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatalf("tasks status = %d", response.StatusCode)
	}
	_ = response.Body.Close()

	request, _ = http.NewRequest(http.MethodPost, base+"/api/v1/tasks", strings.NewReader("{"))
	request.Header.Set("X-Anto-Token", server.Ready().Token)
	request.Header.Set("Content-Type", "application/json")
	response, err = client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("invalid task status = %d", response.StatusCode)
	}
	var apiError struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.NewDecoder(response.Body).Decode(&apiError); err != nil {
		t.Fatal(err)
	}
	_ = response.Body.Close()
	if apiError.Error.Code != "invalid_request" {
		t.Fatalf("unexpected error code %q", apiError.Error.Code)
	}

	request, _ = http.NewRequest(http.MethodPost, base+"/api/v1/tasks/missing/cancel", nil)
	request.Header.Set("X-Anto-Token", server.Ready().Token)
	response, err = client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("cancel status = %d", response.StatusCode)
	}
	_ = response.Body.Close()
}

func TestTranslatorConnectionEndpoint(t *testing.T) {
	engine := &connectionContractTranslator{contractTranslator: newContractTranslator()}
	engine.id = "connection-contract-test"
	repository.GetTranslators().Register(engine)
	server, base, done := startContractServer(t)
	defer stopContractServer(t, server, done)
	client := &http.Client{Timeout: time.Second}
	url := base + "/api/v1/translators/connection-contract-test/test-connection"

	request, _ := http.NewRequest(http.MethodPost, url, strings.NewReader(`{"region":"cn-hangzhou"}`))
	response, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	_ = response.Body.Close()
	if response.StatusCode != http.StatusUnauthorized {
		t.Fatalf("unauthenticated status = %d", response.StatusCode)
	}

	request, _ = http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set("X-Anto-Token", server.Ready().Token)
	response, err = client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	_ = response.Body.Close()
	if response.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("GET status = %d", response.StatusCode)
	}

	request, _ = http.NewRequest(http.MethodPost, url, strings.NewReader(`{"region":"cn-hangzhou"}`))
	request.Header.Set("X-Anto-Token", server.Ready().Token)
	response, err = client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	var result translator.ConnectionResult
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	_ = response.Body.Close()
	if response.StatusCode != http.StatusOK || engine.testedRegion != "cn-hangzhou" || result.Endpoint != "example.invalid" {
		t.Fatalf("connection result: status=%d, region=%q, result=%+v", response.StatusCode, engine.testedRegion, result)
	}

	for region, wantStatus := range map[string]int{"invalid": http.StatusBadRequest, "offline": http.StatusBadGateway} {
		request, _ = http.NewRequest(http.MethodPost, url, strings.NewReader(`{"region":"`+region+`"}`))
		request.Header.Set("X-Anto-Token", server.Ready().Token)
		response, err = client.Do(request)
		if err != nil {
			t.Fatal(err)
		}
		_ = response.Body.Close()
		if response.StatusCode != wantStatus {
			t.Fatalf("region %s status = %d, want %d", region, response.StatusCode, wantStatus)
		}
	}
}

func TestServerEventsBeginWithSnapshot(t *testing.T) {
	server, base, done := startContractServer(t)
	defer stopContractServer(t, server, done)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	request, _ := http.NewRequestWithContext(ctx, http.MethodGet, base+"/api/v1/events", nil)
	request.Header.Set("X-Anto-Token", server.Ready().Token)
	response, err := (&http.Client{Timeout: time.Second}).Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	line, err := bufio.NewReader(response.Body).ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(line) != "event: tasks.snapshot" {
		t.Fatalf("unexpected SSE event line %q", line)
	}
}

func TestServerCreatesAndTracksTask(t *testing.T) {
	repository.GetTranslators().Register(newContractTranslator())
	server, base, done := startContractServer(t)
	defer stopContractServer(t, server, done)

	file := filepath.Join(t.TempDir(), "source.srt")
	if err := os.WriteFile(file, []byte("1\n00:00:00,000 --> 00:00:01,000\nhello\n"), 0600); err != nil {
		t.Fatal(err)
	}
	payloadBytes, err := json.Marshal(core.TaskRequest{
		Files:           []string{file},
		TranslatorID:    "contract-test",
		FromLang:        "en",
		ToLang:          "zh",
		TranslateMode:   "全量翻译",
		MainTrackReport: "来源语种",
		FlagTrackExport: 2,
	})
	if err != nil {
		t.Fatal(err)
	}
	payload := string(payloadBytes)
	request, _ := http.NewRequest(http.MethodPost, base+"/api/v1/tasks", strings.NewReader(payload))
	request.Header.Set("X-Anto-Token", server.Ready().Token)
	request.Header.Set("Content-Type", "application/json")
	response, err := (&http.Client{Timeout: time.Second}).Do(request)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusAccepted {
		t.Fatalf("create task status = %d", response.StatusCode)
	}
	var created core.Task
	if err := json.NewDecoder(response.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}
	_ = response.Body.Close()

	client := &http.Client{Timeout: time.Second}
	var current core.Task
	for attempt := 0; attempt < 30; attempt++ {
		request, _ = http.NewRequest(http.MethodGet, base+"/api/v1/tasks/"+created.ID, nil)
		request.Header.Set("X-Anto-Token", server.Ready().Token)
		response, err = client.Do(request)
		if err != nil {
			t.Fatal(err)
		}
		if err := json.NewDecoder(response.Body).Decode(&current); err != nil {
			t.Fatal(err)
		}
		_ = response.Body.Close()
		if current.Status == core.TaskCompleted || current.Status == core.TaskFailed {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if current.Status != core.TaskCompleted {
		t.Fatalf("task did not complete: %+v", current)
	}
	if current.OutputDirectory == "" || len(current.OutputFiles) != 1 {
		t.Fatalf("task output metadata missing: %+v", current)
	}
}

func TestServerCancelsRunningTask(t *testing.T) {
	repository.GetTranslators().Register(newSlowContractTranslator())
	server, base, done := startContractServer(t)
	defer stopContractServer(t, server, done)

	file := filepath.Join(t.TempDir(), "source.srt")
	if err := os.WriteFile(file, []byte("1\n00:00:00,000 --> 00:00:01,000\nhello\n"), 0600); err != nil {
		t.Fatal(err)
	}
	payload, err := json.Marshal(core.TaskRequest{
		Files:           []string{file},
		TranslatorID:    "contract-slow",
		FromLang:        "en",
		ToLang:          "zh",
		TranslateMode:   "全量翻译",
		MainTrackReport: "来源语种",
		FlagTrackExport: 2,
	})
	if err != nil {
		t.Fatal(err)
	}
	request, _ := http.NewRequest(http.MethodPost, base+"/api/v1/tasks", strings.NewReader(string(payload)))
	request.Header.Set("X-Anto-Token", server.Ready().Token)
	request.Header.Set("Content-Type", "application/json")
	response, err := (&http.Client{Timeout: time.Second}).Do(request)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusAccepted {
		t.Fatalf("create task status = %d", response.StatusCode)
	}
	var created core.Task
	if err := json.NewDecoder(response.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}
	_ = response.Body.Close()

	request, _ = http.NewRequest(http.MethodPost, base+"/api/v1/tasks/"+created.ID+"/cancel", nil)
	request.Header.Set("X-Anto-Token", server.Ready().Token)
	response, err = (&http.Client{Timeout: time.Second}).Do(request)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusAccepted {
		t.Fatalf("cancel status = %d", response.StatusCode)
	}
	_ = response.Body.Close()

	client := &http.Client{Timeout: time.Second}
	for attempt := 0; attempt < 30; attempt++ {
		request, _ = http.NewRequest(http.MethodGet, base+"/api/v1/tasks/"+created.ID, nil)
		request.Header.Set("X-Anto-Token", server.Ready().Token)
		response, err = client.Do(request)
		if err != nil {
			t.Fatal(err)
		}
		var current core.Task
		if err := json.NewDecoder(response.Body).Decode(&current); err != nil {
			t.Fatal(err)
		}
		_ = response.Body.Close()
		if current.Status == core.TaskCancelled {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("task did not reach cancelled status")
}

type contractConfig struct{ translator.DefaultConfig }

type connectionContractTranslator struct {
	*contractTranslator
	testedRegion string
}

func (t *connectionContractTranslator) TestConnection(_ context.Context, region string) (*translator.ConnectionResult, error) {
	t.testedRegion = region
	switch region {
	case "invalid":
		return nil, translator.ErrInvalidConnectionRegion
	case "offline":
		return nil, errors.New("endpoint unreachable")
	}
	return &translator.ConnectionResult{Endpoint: "example.invalid", LatencyMs: 12}, nil
}

func (contractConfig) GetQPS() int        { return 1 }
func (contractConfig) GetMaxCharNum() int { return 1000 }

type contractTranslator struct {
	id            string
	config        *contractConfig
	waitForCancel bool
}

func newContractTranslator() *contractTranslator {
	return &contractTranslator{id: "contract-test", config: &contractConfig{}}
}

func newSlowContractTranslator() *contractTranslator {
	return &contractTranslator{id: "contract-slow", config: &contractConfig{}, waitForCancel: true}
}

func (t *contractTranslator) Init(translator.ImplConfig)              {}
func (t *contractTranslator) GetId() string                           { return t.id }
func (t *contractTranslator) GetShortId() string                      { return "ct" }
func (t *contractTranslator) GetName() string                         { return "Contract test" }
func (t *contractTranslator) GetCfg() translator.ImplConfig           { return t.config }
func (t *contractTranslator) GetLangSupported() []translator.LangPair { return nil }
func (t *contractTranslator) GetSep() string                          { return "" }
func (t *contractTranslator) IsValid() bool                           { return true }
func (t *contractTranslator) Translate(ctx context.Context, args *translator.TranslateArgs) (*translator.TranslateRes, error) {
	if t.waitForCancel {
		<-ctx.Done()
	}
	return &translator.TranslateRes{Results: []*translator.TranslateResBlock{{Id: args.TextContent, TextTranslated: "你好"}}}, nil
}

func startContractServer(t *testing.T) (*Server, string, chan error) {
	t.Helper()
	server, err := New(core.NewService())
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- server.Serve() }()
	base := "http://127.0.0.1:" + strconv.Itoa(server.Ready().Port)
	return server, base, done
}

func stopContractServer(t *testing.T, server *Server, done chan error) {
	t.Helper()
	if err := server.Close(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
}
