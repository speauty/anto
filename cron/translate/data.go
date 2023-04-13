package translate

import (
	"anto/cron/writer"
	"anto/dependency/service/translator"
	"anto/lib/srt"
	_type "anto/type"
	"fmt"
)

type SrtTranslateData struct {
	PrtSrt  *srt.Srt
	PtrOpts *SrtTranslateOpts
}

func (customData *SrtTranslateData) toSrtWriterData() *writer.SrtWriterData {
	tmpData := &writer.SrtWriterData{
		FileNameSaved: customData.fileNameSavedBuilder(),
		PrtSrt:        customData.PrtSrt,
		PtrOpts: &srt.EncodeOpt{
			FlagIsInverse:   customData.PtrOpts.MainTrackReport == _type.LangDirectionTo,
			FlagTrackExport: customData.PtrOpts.FlagTrackExport,
		},
	}
	return tmpData
}

func (customData *SrtTranslateData) fileNameSavedBuilder() string {
	newFileName := customData.PrtSrt.FilePath[0 : len(customData.PrtSrt.FilePath)-4]
	newFileName = fmt.Sprintf(
		"%s.%s2%s.srt", newFileName, customData.PtrOpts.FromLang, customData.PtrOpts.ToLang,
	)
	return newFileName
}

type SrtTranslateOpts struct {
	Translator      translator.InterfaceTranslator
	FromLang        string
	ToLang          string
	TranslateMode   _type.TranslateMode
	MainTrackReport _type.LangDirection
	FlagTrackExport int
}
