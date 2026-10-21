package api

import (
	"anto/core"
	"context"
	"io"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestServerUsesRandomLoopbackPortAndToken(t *testing.T) {
	server, err := New(core.NewService())
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close(context.Background())
	serveDone := make(chan error, 1)
	go func() { serveDone <- server.Serve() }()

	ready := server.Ready()
	if ready.Port == 0 || ready.Token == "" {
		t.Fatalf("invalid ready payload: %+v", ready)
	}
	base := "http://127.0.0.1:" + strconv.Itoa(ready.Port)
	client := &http.Client{Timeout: time.Second}

	response, err := client.Get(base + "/api/v1/translators")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = io.Copy(io.Discard, response.Body)
	_ = response.Body.Close()
	if response.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected unauthorized, got %d", response.StatusCode)
	}

	request, _ := http.NewRequest(http.MethodGet, base+"/api/v1/translators", nil)
	request.Header.Set("X-Anto-Token", ready.Token)
	response, err = client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = io.Copy(io.Discard, response.Body)
	_ = response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected ok, got %d", response.StatusCode)
	}

	_ = server.Close(context.Background())
	if err := <-serveDone; err != nil {
		t.Fatal(err)
	}
}
