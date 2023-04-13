package srt

import (
	"anto/lib/util"
	"bufio"
	"fmt"
	"io"
	"strings"
)

func (customS *Srt) Decode(fileStream io.Reader) (err error) {
	var lineBytes []byte
	var currentLine string
	var currentBlock *Block

	fileReader := bufio.NewReader(fileStream)
	isHeader := true
	for {
		lineBytes = []byte{}
		currentLine = ""

		lineBytes, err = fileReader.ReadBytes('\n')
		if isHeader {
			if util.HasUTF8Dom(lineBytes) {
				lineBytes = lineBytes[3:]
			}

			isHeader = false
		}
		currentLine = strings.TrimSpace(string(lineBytes))

		if customS.decodeBroken(err) {
			return
		}

		if customS.decodeEOFDone(err, currentLine) {
			err = nil
			break
		}

		if currentLine == "" { // 检测到空行 字幕块间隔行
			if currentBlock != nil {
				if !currentBlock.IsValid() {
					err = fmt.Errorf("字幕块[序列号: %d]无效", currentBlock.SeqNo)
					return
				}
				customS.Blocks = append(customS.Blocks, currentBlock)
				currentBlock = nil
			}
			continue // 可能存在连续空行?
		}

		if currentBlock == nil { // isSeqNo
			currentBlock = new(Block)
			if _, err = currentBlock.decodeSeqNo(currentLine); err != nil {
				err = fmt.Errorf("解析序列号异常: %s", err)
				return
			}
			continue
		}
		if isTimeLine, _ := currentBlock.decodeTimeLine(currentLine); isTimeLine {
			continue
		}
		if isMain, _ := currentBlock.decodeMainTrack(currentLine); isMain {
			continue
		}
		if isSub, _ := currentBlock.decodeSubTrack(currentLine); isSub {
			continue
		}
		err = fmt.Errorf("解析字幕块异常, 出现多余行[%s][最近序列号: %d]", currentLine, currentBlock.SeqNo)
		return
	}

	if currentBlock != nil {
		customS.Blocks = append(customS.Blocks, currentBlock)
		currentBlock = nil
	}
	return
}

func (customS *Srt) decodeBroken(err error) bool {
	return err != nil && err != io.EOF
}

func (customS *Srt) decodeEOFDone(err error, currentLine string) bool {
	return err == io.EOF && currentLine == ""
}
