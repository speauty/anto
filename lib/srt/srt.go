package srt

import "strings"

const timeSep = "-->"

type Srt struct {
	FileName       string
	FilePath       string
	FileSign       string
	FileSize       int
	CntBlock       int
	Blocks         []*Block
	FlagTranslated bool
}

type Block struct {
	SeqNo     int
	TimeStart string
	TimeEnd   string
	TimeSep   string
	MainTrack string
	SubTrack  string
}

// FileNameSync 从FilePath中解析文件名称-强行覆盖
func (customS *Srt) FileNameSync() {
	if customS.FilePath == "" {
		return
	}
	filepathArray := strings.Split(customS.FilePath, "/")
	filename := filepathArray[len(filepathArray)-1] // 可能要做其他处理, 先保留这个中间转换
	customS.FileName = filename
}
