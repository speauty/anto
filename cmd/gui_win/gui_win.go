package main

import (
	"anto/bootstrap"
	"anto/platform/win"
	"context"
)

func main() {
	ctx := context.Background()
	bootstrap.Boot(ctx)
	win.Run(ctx)
}
