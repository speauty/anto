package detector

import (
	"anto/common"
	"anto/cron/reader"
	"anto/cron/translate"
	"anto/domain/service/translator"
)

type StrDetectorData struct {
	Translator      translator.ImplTranslator
	FromLang        string
	ToLang          string
	TranslateMode   common.TranslateMode
	MainTrackReport common.LangDirection
	SrtFile         string
	SrtDir          string
	FlagTrackExport int
}

func (customData StrDetectorData) toReaderData(filePath string) *reader.SrtReaderData {
	return &reader.SrtReaderData{
		FilePath: filePath,
		PtrTranslatorOpts: &translate.SrtTranslateOpts{
			Translator: customData.Translator,
			FromLang:   customData.FromLang, ToLang: customData.ToLang,
			TranslateMode: customData.TranslateMode, MainTrackReport: customData.MainTrackReport,
			FlagTrackExport: customData.FlagTrackExport,
		},
	}
}
