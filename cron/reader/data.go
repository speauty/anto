package reader

import (
	"anto/cron/translate"
	"anto/lib/srt"
)

type SrtReaderData struct {
	FilePath          string
	PrtSrt            *srt.Srt
	PtrTranslatorOpts *translate.SrtTranslateOpts
}

func (customData *SrtReaderData) toTranslateData() *translate.SrtTranslateData {
	return &translate.SrtTranslateData{
		PrtSrt:  customData.PrtSrt,
		PtrOpts: customData.PtrTranslatorOpts,
	}
}
