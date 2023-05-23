package main

import (
	"anto/bootstrap"
	"anto/platform/cross_platform_fyne"
	"context"
)

func main() {
	ctx := context.Background()
	bootstrap.Boot(ctx)
	cross_platform_fyne.API().Init(nil)
	cross_platform_fyne.API().Run()
}
