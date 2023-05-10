package main

import (
	"anto/bootstrap"
	"anto/platform/cli"
	"context"
)

func main() {
	ctx := context.Background()
	bootstrap.Boot(ctx)

	cli.Execute()
}
