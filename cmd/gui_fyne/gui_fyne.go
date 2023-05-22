package main

import (
	"anto/bootstrap"
	"anto/platform/e_fyne"
	"context"
)

func main() {
	ctx := context.Background()
	bootstrap.Boot(ctx)
	e_fyne.API().Init(nil)
	e_fyne.API().Run()
}
