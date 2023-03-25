package tt_srt

type EncodeOpt struct {
	FlagIsInverse bool
}

func (customS *Srt) Encode(opts *EncodeOpt) ([]byte, error) {
	var fsBytes []byte
	for _, block := range customS.Blocks {
		fsBytes = append(fsBytes, block.encode(opts.FlagIsInverse)...)
		fsBytes = append(fsBytes, '\n')
	}
	return fsBytes, nil

}
