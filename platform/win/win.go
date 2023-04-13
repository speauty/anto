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
	ui.GetInstance().RegisterMenus(GetInstance().GetMenus())

	ui.GetInstance().RegisterPages(page.GetSettings(), page.GetSubripTranslate())

	if err := ui.GetInstance().Init(cfg.GetInstance().UI); err != nil {
		panic(err)
	}

	_ = ui.GetInstance().GoPage(page.GetSubripTranslate().GetId())

	nohup.NewResident(
		ctx,
		detector.Singleton(), reader.Singleton(), translate.Singleton(), writer.Singleton(),
		ui.GetInstance(),
	)
}
