package tt_srt

type EncodeOpt struct {
	FlagTrackExport int // 导出轨道模式 0-全轨 1-主轨 2-副轨
	FlagIsInverse   bool
}

func (customS *Srt) Encode(opts *EncodeOpt) ([]byte, error) {
	var fsBytes []byte
	for _, block := range customS.Blocks {
		fsBytes = append(fsBytes, block.encode(opts.FlagIsInverse, opts.FlagTrackExport)...)
		fsBytes = append(fsBytes, '\n')
	}
	return fsBytes, nil

}
