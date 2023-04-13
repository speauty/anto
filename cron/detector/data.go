package detector

import (
	_type "anto/common"
	"anto/cron/reader"
	"anto/cron/translate"
	"anto/dependency/service/translator"
)

type StrDetectorData struct {
	Translator      translator.InterfaceTranslator
	FromLang        string
	ToLang          string
	TranslateMode   _type.TranslateMode
	MainTrackReport _type.LangDirection
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
