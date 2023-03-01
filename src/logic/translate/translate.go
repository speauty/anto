package translate

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang-module/carbon"
	mt2 "gui.subtitle/src/srv/mt"
	aliyun2 "gui.subtitle/src/srv/mt/aliyun"
	"io"
	"regexp"
	"strings"
	"sync"
	"unicode"
)

type Block struct {
	Idx       string
	TimeStr   string
	TextCH    string
	TextEN    string
	TextENCnt int
}

func Reader(ir io.Reader) ([]*Block, error) {
	fileReader := bufio.NewReader(ir)
	var subtitles []*Block

	tmpST := new(Block)

	for {
		lineBytes, err := fileReader.ReadBytes('\n')

		//去掉字符串首尾空白字符，返回字符串
		line := strings.TrimSpace(string(lineBytes))

		if err != nil && err != io.EOF {
			return nil, err
		}
		if err == io.EOF && line == "" {
			break
		}
		lineLen := len(line)
		if lineLen == 0 { // 跳过空行
			continue
		}

		if tmpST.Idx == "" {
			tmpST.Idx = line
		} else {
			if strings.Contains(line, " --> ") {
				tmpST.TimeStr = line
			} else {
				if hasZH(line) {
					tmpST.TextCH = line
				} else {
					tmpST.TextEN = line
					tmpST.TextENCnt = lineLen
				}
			}
		}
		if tmpST.TextEN != "" {
			subtitles = append(subtitles, tmpST)
			tmpST = new(Block)
		}
	}

	return subtitles, nil
}

func Writer(iw io.Writer, contents []*Block) (int, error) {
	cnt := 0
	for _, content := range contents {
		buffer := bytes.NewBufferString(fmt.Sprintf("%s\n%s\n", content.Idx, content.TimeStr))
		if content.TextCH != "" {
			buffer.WriteString(content.TextCH)
			buffer.WriteByte('\n')
		}
		if content.TextEN != "" {
			buffer.WriteString(content.TextEN)
			buffer.WriteByte('\n')
		}
		buffer.WriteByte('\n')
		write, err := iw.Write(buffer.Bytes())
		if err != nil {
			return 0, err
		}
		cnt += write
	}
	return cnt, nil
}

func Translate(mt interface{}, contents []*Block) ([]string, int) {
	var contentsChunked []map[string]string
	tmpMap := map[string]string{}
	tmpLen := 0
	for _, content := range contents {
		if content.TextCH != "" { // 略过已经翻译了的
			continue
		}
		if tmpLen+content.TextENCnt >= 5000 { // 单次批量翻译最大值
			contentsChunked = append(contentsChunked, tmpMap)
			tmpLen = 0
			tmpMap = map[string]string{}
		}
		tmpLen += content.TextENCnt
		tmpMap[content.Idx] = content.TextEN
	}
	if tmpLen != 0 {
		contentsChunked = append(contentsChunked, tmpMap)
	}
	wg := sync.WaitGroup{}
	var results []string
	cntError := 0
	for idx, m := range contentsChunked {
		wg.Add(1)
		localIdx := idx
		mChunked := m
		go func() {
			defer wg.Done()
			timeStartedAt := carbon.Now()
			marshal, _ := json.Marshal(mChunked)
			translates, err := mt.(mt2.MT).TextBatchTranslate(new(aliyun2.TextBatchTranslateArg).New(string(marshal)))
			if err != nil {
				results = append(results, fmt.Sprintf("[%s]%s翻译失败, 错误: %s", carbon.Now(), mt.(mt2.MT).GetName(), err))
				cntError++
				return
			}
			lineMatchedCnt := 0
			for _, blockTranslated := range translates {
				for contentIdx, content := range contents {
					if blockTranslated.Idx == content.Idx {
						contents[contentIdx].TextCH = blockTranslated.StrTranslated
						lineMatchedCnt++
					}
				}
			}
			results = append(results, fmt.Sprintf(
				"[%s]%s翻译成功, 序号: %d, 字幕行数: %d, 耗时(s): %d",
				carbon.Now(), mt.(mt2.MT).GetName(), localIdx+1, lineMatchedCnt, carbon.Now().DiffAbsInSeconds(timeStartedAt),
			))
		}()

	}
	wg.Wait()
	return results, cntError
}

func hasZH(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) || (regexp.MustCompile("[\u3002\uff1b\uff0c\uff1a\u201c\u201d\uff08\uff09\u3001\uff1f\u300a\u300b]").MatchString(string(r))) {
			return true
		}
	}
	return false
}
