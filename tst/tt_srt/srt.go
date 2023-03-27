package tt_srt

const timeSep = "-->"

type Srt struct {
	FileName       string
	FilePath       string
	FileSign       string
	FileSize       int
	CntBlock       int
	Blocks         []*SrtBlock
	FlagTranslated bool
}

type SrtBlock struct {
	SeqNo     int
	TimeStart string
	TimeEnd   string
	TimeSep   string
	MainTrack string
	SubTrack  string
}
