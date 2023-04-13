package writer

import "anto/lib/srt"

type SrtWriterData struct {
	FileNameSaved string
	PrtSrt        *srt.Srt
	PtrOpts       *srt.EncodeOpt
}
