package win

import (
	"anto/cfg"
	"anto/cron/detector"
	"anto/cron/reader"
	"anto/cron/translate"
	"anto/cron/writer"
	"anto/lib/log"
	"anto/lib/nohup"
	"anto/platform/win/page"
	"anto/platform/win/ui"
	"context"
)

func Run(ctx context.Context) {
	ui.Singleton().RegisterMenus(GetInstance().GetMenus())

	ui.Singleton().RegisterPages(page.GetSettings(), page.GetSubripTranslate(), page.GetAboutUs())

	if err := ui.Singleton().Init(cfg.Singleton().UI); err != nil {
		log.Singleton().ErrorF("UI启动崩溃, 错误: %s", err)
		panic(err)
	}

	_ = ui.Singleton().GoPage(page.GetSubripTranslate().GetId())

	nohup.NewResident(
		ctx,
		detector.Singleton(), reader.Singleton(), translate.Singleton(), writer.Singleton(),
		ui.Singleton(),
	)
}
