package translate

import (
	_const "anto/common"
	"anto/cron/writer"
	"anto/dependency/service/translator"
	"anto/lib/srt"
	"fmt"
	"strings"
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
			FlagIsInverse:   customData.PtrOpts.MainTrackReport == _const.LangDirectionTo,
			FlagTrackExport: customData.PtrOpts.FlagTrackExport,
		},
	}
	return tmpData
}

func (customData *SrtTranslateData) fileNameSavedBuilder() string {
	newFileName := customData.PrtSrt.FilePath[0 : len(customData.PrtSrt.FilePath)-4]
	newFileName = fmt.Sprintf(
		"%s.%s.%s.%s2%s.srt", newFileName, _const.AppName, customData.PtrOpts.Translator.GetShortId(), strings.ToLower(customData.PtrOpts.FromLang), strings.ToLower(customData.PtrOpts.ToLang),
	)
	return newFileName
}

type SrtTranslateOpts struct {
	Translator      translator.InterfaceTranslator
	FromLang        string
	ToLang          string
	TranslateMode   _const.TranslateMode
	MainTrackReport _const.LangDirection
	FlagTrackExport int
}
