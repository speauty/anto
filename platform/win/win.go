package win

import (
	"anto/cfg"
	"anto/cron/detector"
	"anto/cron/reader"
	"anto/cron/translate"
	"anto/cron/writer"
	"anto/lib/nohup"
	"anto/platform/win/page"
	"anto/platform/win/ui"
	"context"
)

func Run(ctx context.Context) {
	ui.Singleton().RegisterMenus(GetInstance().GetMenus())

	ui.Singleton().RegisterPages(page.GetSettings(), page.GetSubripTranslate())

	if err := ui.Singleton().Init(cfg.Singleton().UI); err != nil {
		panic(err)
	}

	_ = ui.Singleton().GoPage(page.GetSubripTranslate().GetId())

	nohup.NewResident(
		ctx,
		detector.Singleton(), reader.Singleton(), translate.Singleton(), writer.Singleton(),
		ui.Singleton(),
	)
}
