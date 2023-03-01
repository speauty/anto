package translate

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-module/carbon"
	mt2 "gui.subtitle/src/srv/mt"
	aliyun2 "gui.subtitle/src/srv/mt/aliyun"
	"gui.subtitle/src/srv/mt/bd"
	"gui.subtitle/src/util/lang"
	"io"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
	"unicode"
)

type Block struct {
	Idx       string
	TimeStr   string
	TextZH    string
	TextZHCnt int
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
					tmpST.TextZH = line
					tmpST.TextZHCnt = lineLen
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
		if content.TextZH != "" {
			buffer.WriteString(content.TextZH)
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

func Translate(ctx context.Context, mt interface{}, contents []*Block, fromLanguage lang.StrLang, toLanguage lang.StrLang) ([]string, int, error) {
	if err := preCheckBlocks(contents, fromLanguage); err != nil {
		return nil, 0, fmt.Errorf("预检字幕失败, 错误: %s", err.Error())
	}
	wg := &sync.WaitGroup{}
	var results []string
	cntError := 0
	var lastError error
	cntBlock := len(contents)

	maxCoroutine := 10

	switch mt.(mt2.MT).GetId() {
	case mt2.ALI:
		contentsChunked, err := chunkBlocksForALi(contents, toLanguage)
		if err != nil {
			return nil, 0, fmt.Errorf("阿里翻译字幕分包异常, 错误: %s", err.Error())
		}
		for idx, m := range contentsChunked {
			wg.Add(1)
			localIdx := idx
			mChunked := m
			go func() {
				defer wg.Done()
				timeStartedAt := carbon.Now()
				marshal, _ := json.Marshal(mChunked)
				args := new(aliyun2.TextBatchTranslateArg).New(string(marshal))
				args.FromLanguage = fromLanguage.ToString()
				args.ToLanguage = toLanguage.ToString()
				translates, err := mt.(mt2.MT).TextBatchTranslate(ctx, args)
				if err != nil {
					msg := fmt.Sprintf("[%s]%s失败, 错误: %s", carbon.Now(), mt.(mt2.MT).GetName(), err)
					results = append(results, msg)
					cntError++
					lastError = fmt.Errorf(msg)
					return
				}
				lineMatchedCnt := 0
				for _, blockTranslated := range translates {
					for contentIdx, content := range contents {
						if blockTranslated.Idx == content.Idx {
							contents[contentIdx].TextZH = blockTranslated.StrTranslated
							lineMatchedCnt++
						}
					}
				}
				results = append(results, fmt.Sprintf(
					"[%s]%s成功, 序号: %d, 字幕行数: %d, 耗时(s): %d",
					carbon.Now(), mt.(mt2.MT).GetName(), localIdx+1, lineMatchedCnt, carbon.Now().DiffAbsInSeconds(timeStartedAt),
				))
			}()
		}
	case mt2.BAIDU:
		for coroutineIdx := 0; coroutineIdx < maxCoroutine && coroutineIdx < cntBlock; coroutineIdx++ {
			wg.Add(1)
			go func(localCtx context.Context, localWG *sync.WaitGroup, localCoroutineIdx int) {
				defer localWG.Done()
				timeStart := carbon.Now()
				cntBlockTranslated := 0
				for blockIdx := 0; blockIdx < cntBlock; blockIdx++ {
					if blockIdx%10 != localCoroutineIdx {
						continue
					}
					currentBlock := contents[blockIdx]
					sourceText := currentBlock.TextZH
					if fromLanguage == lang.EN {
						sourceText = currentBlock.TextEN
					}

					args := new(bd.TextTranslateArg).New(sourceText)
					args.FromLanguage = fromLanguage.ToString()
					args.ToLanguage = toLanguage.ToString()
					var err error
					translateResp := new(mt2.TextTranslateResp)

					for failIdx := 0; failIdx < 3; failIdx++ {
						translateResp, err = mt.(mt2.MT).TextTranslate(ctx, args)
						if err != nil {
							err = fmt.Errorf("[%s]%s失败, 协程序号: %d, 字幕序号: %s, 错误: %s", carbon.Now(), mt.(mt2.MT).GetName(), localCoroutineIdx, currentBlock.Idx, err)
							time.Sleep(time.Second)
							continue
						}
						break
					}
					if err != nil {
						results = append(results, err.Error())
						cntError++
						lastError = err
						breakFlags := []string{
							"52002", "52003", "54001", "54004", "58000", "58001", "58002", "90107",
						}
						for _, flag := range breakFlags {
							if strings.Contains(err.Error(), flag) {
								runtime.Goexit()
							}
						}

						continue
					}

					if toLanguage == lang.ZH {
						contents[blockIdx].TextZH = translateResp.StrTranslated
					} else if toLanguage == lang.EN {
						contents[blockIdx].TextEN = translateResp.StrTranslated
					}
					cntBlockTranslated++
				}
				if cntBlockTranslated > 0 {
					results = append(results, fmt.Sprintf(
						"[%s]%s成功, 协程序号: %d, 字幕行数: %d, 耗时(s): %d",
						carbon.Now(), mt.(mt2.MT).GetName(), localCoroutineIdx, cntBlockTranslated, carbon.Now().DiffAbsInSeconds(timeStart),
					))
				} else {
					results = append(results, fmt.Sprintf(
						"[%s]%s, 协程序号: %d, 错误: 空转, 耗时(s): %d",
						carbon.Now(), mt.(mt2.MT).GetName(), localCoroutineIdx, carbon.Now().DiffAbsInSeconds(timeStart),
					))
				}

			}(ctx, wg, coroutineIdx)
		}
	}
	wg.Wait()
	return results, cntError, lastError
}

// preCheckBlocks 预检字幕块, 主要保证需要翻译的字幕块存在
func preCheckBlocks(contents []*Block, fromLanguage lang.StrLang) error {
	if !(fromLanguage == lang.EN || fromLanguage == lang.ZH) {
		return fmt.Errorf("暂未实现来源语言[%s]的字幕链接", fromLanguage.GetCH())
	}

	for _, content := range contents {
		if (fromLanguage == lang.EN && content.TextEN == "") || (fromLanguage == lang.ZH && content.TextZH == "") {
			return fmt.Errorf("当前字幕缺失, 索引: %s, 时间: %s", content.Idx, content.TimeStr)
		}
	}
	return nil
}

// chunkBlocksForALi 阿里云翻译的专属分包函数
func chunkBlocksForALi(contents []*Block, toLanguage lang.StrLang) ([]map[string]string, error) {
	var contentsChunked []map[string]string
	tmpMap := map[string]string{}
	tmpLen := 0
	for _, content := range contents {
		if toLanguage == lang.ZH {
			if content.TextZH != "" { // 略过已经翻译了的
				continue
			}
			if tmpLen+content.TextENCnt >= 5000 { // 单次批量翻译最大值
				contentsChunked = append(contentsChunked, tmpMap)
				tmpLen = 0
				tmpMap = map[string]string{}
			}
		} else if toLanguage == lang.EN {
			if content.TextEN != "" { // 略过已经翻译了的
				continue
			}
			if tmpLen+content.TextZHCnt >= 5000 { // 单次批量翻译最大值
				contentsChunked = append(contentsChunked, tmpMap)
				tmpLen = 0
				tmpMap = map[string]string{}
			}
		} else {
			return nil, fmt.Errorf("暂未实现目标语言[%s]的分包操作", toLanguage.GetCH())
		}

		tmpLen += content.TextENCnt
		tmpMap[content.Idx] = content.TextEN
	}
	if tmpLen != 0 {
		contentsChunked = append(contentsChunked, tmpMap)
	}
	return contentsChunked, nil
}

func hasZH(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) || (regexp.MustCompile("[\u3002\uff1b\uff0c\uff1a\u201c\u201d\uff08\uff09\u3001\uff1f\u300a\u300b]").MatchString(string(r))) {
			return true
		}
	}
	return false
}
