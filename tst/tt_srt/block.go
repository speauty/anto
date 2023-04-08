package tt_srt

import (
	"fmt"
	"github.com/golang-module/carbon"
	"strconv"
	"strings"
)

func (customB *SrtBlock) IsValid() bool {
	return customB.SeqNo > 0 &&
		customB.TimeStart != "" && customB.TimeEnd != "" && customB.TimeSep != "" &&
		customB.MainTrack != ""
}

func (customB *SrtBlock) decodeSeqNo(lineStr string) (isSeqNo bool, err error) {
	if lineStr == "" || customB.SeqNo > 0 {
		return
	}
	var tmpNum int
	tmpNum, err = strconv.Atoi(lineStr)
	if err != nil {
		return
	}
	if tmpNum <= 0 {
		err = fmt.Errorf("无效序列号[%d]", tmpNum)
		return
	}
	customB.SeqNo = tmpNum
	isSeqNo = true
	return
}

func (customB *SrtBlock) decodeTimeLine(lineStr string) (isTimeLine bool, err error) {
	if lineStr == "" || !strings.Contains(lineStr, timeSep) {
		return
	}
	if customB.SeqNo == 0 {
		return
	}
	timeSplit := strings.Split(lineStr, timeSep)
	if len(timeSplit) != 2 {
		return
	}
	firstBlock := strings.TrimSpace(timeSplit[0])
	secondBlock := strings.TrimSpace(timeSplit[1])
	layout := "15:04:05.999"
	if carbon.ParseByLayout(firstBlock, layout).IsInvalid() ||
		carbon.ParseByLayout(secondBlock, layout).IsInvalid() {
		return
	}
	customB.TimeStart = firstBlock
	customB.TimeEnd = secondBlock
	customB.TimeSep = timeSep
	isTimeLine = true
	return
}

func (customB *SrtBlock) decodeMainTrack(lineStr string) (isMain bool, err error) {
	if lineStr == "" || customB.MainTrack != "" {
		return
	}
	if customB.SeqNo == 0 || customB.TimeStart == "" {
		return
	}
	customB.MainTrack = lineStr
	isMain = true
	return
}

func (customB *SrtBlock) decodeSubTrack(lineStr string) (isSub bool, err error) {
	if lineStr == "" || customB.SubTrack != "" {
		return
	}
	if customB.SeqNo == 0 || customB.TimeStart == "" || customB.MainTrack == "" {
		return
	}
	customB.SubTrack = lineStr
	isSub = true
	return
}

func (customB *SrtBlock) encode(flagInverse bool, flagTrackExportedMode int) []byte {
	blockStr := fmt.Sprintf("%d\n%s %s %s\n", customB.SeqNo,
		customB.TimeStart, customB.TimeSep, customB.TimeEnd)
	if customB.SubTrack == "" {
		blockStr = fmt.Sprintf("%s%s\n", blockStr, customB.MainTrack)
	} else {
		if flagTrackExportedMode != 0 {
			if flagTrackExportedMode == 1 {
				if flagInverse == true {
					blockStr = fmt.Sprintf("%s%s\n", blockStr, customB.SubTrack)
				} else {
					blockStr = fmt.Sprintf("%s%s\n", blockStr, customB.MainTrack)
				}
			} else if flagTrackExportedMode == 2 {
				if flagInverse == true {
					blockStr = fmt.Sprintf("%s%s\n", blockStr, customB.MainTrack)
				} else {
					blockStr = fmt.Sprintf("%s%s\n", blockStr, customB.SubTrack)
				}
			}
		} else {
			if flagInverse == false {
				blockStr = fmt.Sprintf("%s%s\n%s\n", blockStr, customB.MainTrack, customB.SubTrack)
			} else {
				blockStr = fmt.Sprintf("%s%s\n%s\n", blockStr, customB.SubTrack, customB.MainTrack)
			}
		}
	}

	return []byte(blockStr)
}
