package main

import (
	"anto/bootstrap"
	"anto/cron/detector"
	"anto/cron/reader"
	"anto/cron/translate"
	"anto/cron/writer"
	"anto/lib/nohup"
	"anto/platform/cross_platform_fyne"
	"context"
)

func main() {
	ctx := context.Background()
	bootstrap.Boot(ctx)
	cross_platform_fyne.API().Init(nil)

	{ // 绑定消息通道
		detector.Singleton().SetMsgRedirect(cross_platform_fyne.API().ChanConsole())
		translate.Singleton().SetMsgRedirect(cross_platform_fyne.API().ChanConsole())
	}

	nohup.NewResident(
		ctx,
		detector.Singleton(), reader.Singleton(), translate.Singleton(), writer.Singleton(),
		cross_platform_fyne.API(),
	)
}
