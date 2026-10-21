package main

import (
	"anto/bootstrap"
	"anto/core"
	"anto/core/api"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	bootstrap.Boot(ctx)
	service := core.NewService()
	server, err := api.New(service)
	if err != nil {
		panic(err)
	}
	ready, _ := json.Marshal(server.Ready())
	// Flutter reads this one machine-readable line to learn the random port and token.
	fmt.Printf("ANTO_CORE_READY %s\n", ready)
	serveErr := make(chan error, 1)
	go func() { serveErr <- server.Serve() }()
	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = server.Close(shutdownCtx)
	case err := <-serveErr:
		if err != nil {
			panic(err)
		}
	}
}
