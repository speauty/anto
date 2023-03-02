package translate

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-module/carbon"
	"gui.subtitle/src/srv/mt"
	aliyun2 "gui.subtitle/src/srv/mt/aliyun"
	"gui.subtitle/src/srv/mt/bd"
	"gui.subtitle/src/srv/mt/youdao"
	"gui.subtitle/src/util"
	"gui.subtitle/src/util/lang"
	"io"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
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

func Translate(ctx context.Context, mtEngine interface{}, contents []*Block, fromLanguage lang.StrLang, toLanguage lang.StrLang) ([]string, int, error) {
	if err := preCheckBlocks(contents, fromLanguage); err != nil {
		return nil, 0, fmt.Errorf("预检字幕失败, 错误: %s", err.Error())
	}
	wg := &sync.WaitGroup{}
	var results []string
	cntError := new(atomic.Int32)
	var lastError error
	cntBlock := len(contents)
	coroutineCtrlCtx, coroutineCtrlCtxCancelFunc := context.WithCancel(ctx)
	defer coroutineCtrlCtxCancelFunc()

	maxCoroutine := 10
	maxRetry := 3
	cntBlockTranslated := new(atomic.Int32)
	cntBlockTranslated.Store(0)
	timeStart := carbon.Now()

	switch mtEngine.(mt.MT).GetId() {
	case mt.IdALiYun:
		blockChunked, err := chunkBlocksForALi(contents, toLanguage)
		if err != nil {
			return nil, 0, fmt.Errorf("%s字幕分包异常, 错误: %s", mtEngine.(mt.MT).GetName(), err.Error())
		}
		cntBlockChunked := len(blockChunked)
		if cntBlockChunked < maxCoroutine {
			maxCoroutine = cntBlockChunked
		}
		blockChan := make(chan string, maxCoroutine*3)
		go func() {
			for _, block := range blockChunked {
				marshal, _ := json.Marshal(block)
				blockChan <- string(marshal)
			}
			close(blockChan)
		}()

		for coroutineIdx := 0; coroutineIdx < maxCoroutine; coroutineIdx++ {
			if util.IsCtxDone(coroutineCtrlCtx) {
				results = append(results, fmt.Sprintf("[%s]%s失败, 错误: 协程出现中断信号, 停止继续创建协程", carbon.Now(), mtEngine.(mt.MT).GetName()))
				break
			}
			wg.Add(1)
			go func(localCtx context.Context, localCoroutineCtrlCtx context.Context, localCoroutineCtrlCtxCancelFunc context.CancelFunc, localWG *sync.WaitGroup, localCoroutineIdx int, localBlockChan chan string) {
				defer localWG.Done()
				localCoroutineCntTranslated := 0
				localCoroutineTimeStart := carbon.Now()
				for {
					select {
					case block, isOpen := <-blockChan:
						if !isOpen {
							results = append(results, fmt.Sprintf(
								"[%s]协程结束, 引擎: %s, 协程序号: %d, 处理字幕行数: %d, 运行时长(s): %d, 原因: 数据通道关闭, 无数据, 主动退出当前协程",
								carbon.Now(), mtEngine.(mt.MT).GetName(), localCoroutineIdx,
								localCoroutineCntTranslated, carbon.Now().DiffAbsInSeconds(localCoroutineTimeStart),
							))
							runtime.Goexit()
							return
						}
						args := new(aliyun2.TextBatchTranslateArg).New(block)
						args.FromLanguage = fromLanguage.ToString()
						args.ToLanguage = toLanguage.ToString()
						translateResp, translateErr := mtEngine.(mt.MT).TextBatchTranslate(ctx, args)
						if translateErr != nil {
							msg := fmt.Sprintf("[%s]%s翻译失败, 协程序号: %d, 错误: %s", carbon.Now(), mtEngine.(mt.MT).GetName(), localCoroutineIdx, translateErr.Error())
							results = append(results, msg)
							cntError.Add(1)
							lastError = fmt.Errorf(msg)
							localCoroutineCtrlCtxCancelFunc()
							runtime.Goexit()
							return
						}
						for _, blockTranslated := range translateResp {
							for contentIdx, content := range contents {
								if blockTranslated.Idx == content.Idx {
									if toLanguage == lang.ZH {
										contents[contentIdx].TextZH = blockTranslated.StrTranslated
										cntBlockTranslated.Add(1)
										localCoroutineCntTranslated++
									} else if toLanguage == lang.EN {
										contents[contentIdx].TextEN = blockTranslated.StrTranslated
										cntBlockTranslated.Add(1)
										localCoroutineCntTranslated++
									}
								}
							}
						}
					default:
						if util.IsCtxDone(localCoroutineCtrlCtx) {
							results = append(results, fmt.Sprintf(
								"[%s]协程结束, 引擎: %s, 协程序号: %d, 处理字幕行数: %d, 运行时长(s): %d, 错误: 协程出现中断信号, 强制退出",
								carbon.Now(), mtEngine.(mt.MT).GetName(), localCoroutineIdx,
								localCoroutineCntTranslated, carbon.Now().DiffAbsInSeconds(localCoroutineTimeStart),
							))
							runtime.Goexit()
							return
						}
					}
				}
			}(ctx, coroutineCtrlCtx, coroutineCtrlCtxCancelFunc, wg, coroutineIdx, blockChan)
		}
	case mt.IdBaiDu:
		currentCfg := mtEngine.(mt.MT).GetCfg().(*bd.Cfg)
		maxCoroutine = currentCfg.AppVersion.GetQPS() // 和接口版本绑定
		if maxCoroutine > cntBlock {                  // 和字幕块简单比较小, 取相对小值
			maxCoroutine = cntBlock
		}
		if maxCoroutine > mt.MaxCoroutine { // 上限
			maxCoroutine = mt.MaxCoroutine
		}
		blockChunked, cntBlockChunked := chunkBlocksForBaiDu(contents, int(float64(currentCfg.AppVersion.GetLenLimited())*0.8), fromLanguage)
		if maxCoroutine > cntBlockChunked {
			maxCoroutine = cntBlockChunked
		}
		blockChan := make(chan string, maxCoroutine*3)

		go func() {
			for _, block := range blockChunked {
				blockChan <- block
			}
			close(blockChan)
		}()
		for coroutineIdx := 0; coroutineIdx < maxCoroutine; coroutineIdx++ {
			if util.IsCtxDone(coroutineCtrlCtx) {
				results = append(results, fmt.Sprintf("[%s]%s失败, 错误: 协程出现中断信号, 停止继续创建协程", carbon.Now(), mtEngine.(mt.MT).GetName()))
				break
			}
			wg.Add(1)
			go func(localCtx context.Context, localCoroutineCtrlCtx context.Context, localCoroutineCtrlCtxCancelFunc context.CancelFunc, localWG *sync.WaitGroup, localCoroutineIdx int, localBlockChan chan string) {
				defer localWG.Done()
				localCoroutineCntTranslated := 0
				localCoroutineTimeStart := carbon.Now()
				for {
					select {
					case block, isOpen := <-blockChan:
						if !isOpen { // 当前通道已关闭, 并且没有残留数据
							results = append(results, fmt.Sprintf(
								"[%s]协程结束, 引擎: %s, 协程序号: %d, 处理字幕行数: %d, 运行时长(s): %d, 原因: 数据通道关闭, 无数据, 主动退出当前协程",
								carbon.Now(), mtEngine.(mt.MT).GetName(), localCoroutineIdx,
								localCoroutineCntTranslated, carbon.Now().DiffAbsInSeconds(localCoroutineTimeStart),
							))
							runtime.Goexit()
							return
						}
						args := new(bd.TextTranslateArg).New(block)
						args.FromLanguage = fromLanguage.ToString()
						args.ToLanguage = toLanguage.ToString()
						var err error
						var translateResp []mt.TextTranslateResp

						for failIdx := 0; failIdx < currentCfg.AppVersion.GetRetryLimited(); failIdx++ {
							translateResp, err = mtEngine.(mt.MT).TextTranslate(ctx, args)
							if err != nil || translateResp == nil {
								err = fmt.Errorf("[%s]%s翻译失败, 协程序号: %d, 错误: %s", carbon.Now(), mtEngine.(mt.MT).GetName(), localCoroutineIdx, err)
								time.Sleep(time.Millisecond * 100)
								continue
							}
							break
						}
						if err != nil {
							results = append(results, err.Error())
							cntError.Add(1)
							lastError = err

							if bd.ErrSign.IsExit(err) {
								localCoroutineCtrlCtxCancelFunc()
								runtime.Goexit()
								return
							}
						}
						for _, translateRes := range translateResp {
							for idx, content := range contents {
								if fromLanguage == lang.ZH {
									if content.TextZH == translateRes.Idx {
										contents[idx].TextEN = translateRes.StrTranslated
										cntBlockTranslated.Add(1)
										localCoroutineCntTranslated++
									}
								} else if fromLanguage == lang.EN {
									if content.TextEN == translateRes.Idx {
										contents[idx].TextZH = translateRes.StrTranslated
										cntBlockTranslated.Add(1)
										localCoroutineCntTranslated++
									}
								}
							}
						}

					default:
						if util.IsCtxDone(localCoroutineCtrlCtx) {
							results = append(results, fmt.Sprintf(
								"[%s]协程结束, 引擎: %s, 协程序号: %d, 处理字幕行数: %d, 运行时长(s): %d, 错误: 协程出现中断信号, 强制退出",
								carbon.Now(), mtEngine.(mt.MT).GetName(), localCoroutineIdx,
								localCoroutineCntTranslated, carbon.Now().DiffAbsInSeconds(localCoroutineTimeStart),
							))
							runtime.Goexit()
							return
						}
					}
				}
			}(ctx, coroutineCtrlCtx, coroutineCtrlCtxCancelFunc, wg, coroutineIdx, blockChan)
		}
	case mt.IdYouDao:
		blockChunked, cntBlockChunked := chunkBlocksForBaiDu(contents, 2e3, fromLanguage)
		if maxCoroutine > cntBlockChunked {
			maxCoroutine = cntBlockChunked
		}
		blockChan := make(chan string, maxCoroutine*3)
		go func() {
			for _, block := range blockChunked {
				blockChan <- block
			}
			close(blockChan)
		}()
		for coroutineIdx := 0; coroutineIdx < maxCoroutine; coroutineIdx++ {
			if util.IsCtxDone(coroutineCtrlCtx) {
				results = append(results, fmt.Sprintf("[%s]%s失败, 错误: 协程出现中断信号, 停止继续创建协程", carbon.Now(), mtEngine.(mt.MT).GetName()))
				break
			}
			wg.Add(1)
			go func(localCtx context.Context, localCoroutineCtrlCtx context.Context, localCoroutineCtrlCtxCancelFunc context.CancelFunc, localWG *sync.WaitGroup, localCoroutineIdx int, localBlockChan chan string) {
				defer localWG.Done()
				localCoroutineCntTranslated := 0
				localCoroutineTimeStart := carbon.Now()
				for {
					select {
					case block, isOpen := <-blockChan:
						if !isOpen { // 当前通道已关闭, 并且没有残留数据
							results = append(results, fmt.Sprintf(
								"[%s]协程结束, 引擎: %s, 协程序号: %d, 处理字幕行数: %d, 运行时长(s): %d, 原因: 数据通道关闭, 无数据, 主动退出当前协程",
								carbon.Now(), mtEngine.(mt.MT).GetName(), localCoroutineIdx,
								localCoroutineCntTranslated, carbon.Now().DiffAbsInSeconds(localCoroutineTimeStart),
							))
							runtime.Goexit()
							return
						}
						args := new(youdao.TextTranslateArg).New(block)
						args.FromLanguage = fromLanguage.ToString()
						args.ToLanguage = toLanguage.ToString()
						var err error
						var translateResp []mt.TextTranslateResp

						for failIdx := 0; failIdx < maxRetry; failIdx++ {
							translateResp, err = mtEngine.(mt.MT).TextTranslate(ctx, args)
							if err != nil || translateResp == nil {
								err = fmt.Errorf("[%s]%s翻译失败, 协程序号: %d, 错误: %s", carbon.Now(), mtEngine.(mt.MT).GetName(), localCoroutineIdx, err)
								continue
							}
							break
						}
						if err != nil {
							results = append(results, err.Error())
							cntError.Add(1)
							lastError = err
							runtime.Goexit()
							return
						}
						for _, translateRes := range translateResp {
							for idx, content := range contents {
								if fromLanguage == lang.ZH {
									if content.TextZH == translateRes.Idx {
										contents[idx].TextEN = translateRes.StrTranslated
										cntBlockTranslated.Add(1)
										localCoroutineCntTranslated++
									}
								} else if fromLanguage == lang.EN {
									if content.TextEN == translateRes.Idx {
										contents[idx].TextZH = translateRes.StrTranslated
										cntBlockTranslated.Add(1)
										localCoroutineCntTranslated++
									}
								}
							}
						}
					default:
						if util.IsCtxDone(localCoroutineCtrlCtx) {
							results = append(results, fmt.Sprintf(
								"[%s]协程结束, 引擎: %s, 协程序号: %d, 处理字幕行数: %d, 运行时长(s): %d, 错误: 协程出现中断信号, 强制退出",
								carbon.Now(), mtEngine.(mt.MT).GetName(), localCoroutineIdx,
								localCoroutineCntTranslated, carbon.Now().DiffAbsInSeconds(localCoroutineTimeStart),
							))
							runtime.Goexit()
							return
						}
					}
				}
			}(ctx, coroutineCtrlCtx, coroutineCtrlCtxCancelFunc, wg, coroutineIdx, blockChan)
		}
	}
	wg.Wait()
	resStr := "成功"
	if cntBlockTranslated.Load() < int32(cntBlock) {
		resStr = "失败"
	}
	results = append(results, fmt.Sprintf(
		"[%s]翻译完成, 引擎: %s, 协程数量: %d, 翻译行数: %d, 结果: %s, 耗时(s): %d",
		carbon.Now(), mtEngine.(mt.MT).GetName(), maxCoroutine, cntBlockTranslated.Load(), resStr, carbon.Now().DiffAbsInSeconds(timeStart),
	))
	return results, int(cntError.Load()), lastError
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

// chunkBlocksForBaiDu 百度翻译的专属分包函数
func chunkBlocksForBaiDu(contents []*Block, lenLimited int, fromLanguage lang.StrLang) ([]string, int) {
	var blockChunked []string
	tmpStrStack := ""
	for _, currentBlock := range contents {
		sourceText := currentBlock.TextZH
		if fromLanguage == lang.EN {
			sourceText = currentBlock.TextEN
		}
		if tmpStrStack == "" {
			tmpStrStack = sourceText
		} else {
			tmpStrStack = fmt.Sprintf("%s%s%s", tmpStrStack, mt.BlockSep, sourceText)
		}

		if len(tmpStrStack) > lenLimited { // 增量够了, 清理一波
			blockChunked = append(blockChunked, tmpStrStack)
			tmpStrStack = ""
		}
	}
	if tmpStrStack != "" { // 清理残余数据
		blockChunked = append(blockChunked, tmpStrStack)
		tmpStrStack = ""
	}
	return blockChunked, len(blockChunked)
}

func hasZH(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) || (regexp.MustCompile("[\u3002\uff1b\uff0c\uff1a\u201c\u201d\uff08\uff09\u3001\uff1f\u300a\u300b]").MatchString(string(r))) {
			return true
		}
	}
	return false
}
