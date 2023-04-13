package translate

import (
	"anto/cron"
	"anto/cron/writer"
	"anto/dependency/service/translator"
	"anto/lib/log"
	"anto/lib/util"
	_type "anto/type"
	"context"
	"fmt"
	"github.com/golang-module/carbon"
	"runtime"
	"sync"
)

const (
	cronName     = "SRT翻译程序"
	numChanData  = 10
	numChanMsg   = 20
	numCoroutine = 10
)

var (
	apiSingleton  *SrtTranslator
	onceSingleton sync.Once
)

func Singleton() *SrtTranslator {
	onceSingleton.Do(func() {
		apiSingleton = new(SrtTranslator)
		apiSingleton.init()
	})
	return apiSingleton
}

type SrtTranslator struct {
	ctx             context.Context
	ctxFnCancel     context.CancelFunc
	chanData        chan *SrtTranslateData
	chanMsg         chan string
	chanMsgRedirect chan string
	numCoroutine    int
}

func (customCron *SrtTranslator) Push(data *SrtTranslateData) {
	customCron.chanData <- data
}

func (customCron *SrtTranslator) Run(ctx context.Context, fnCancel context.CancelFunc) {
	customCron.ctx = ctx
	customCron.ctxFnCancel = fnCancel

	customCron.jobTranslator()
	customCron.jobMsg()
}

func (customCron *SrtTranslator) Close() {}

func (customCron *SrtTranslator) SetMsgRedirect(chanMsg chan string) {
	customCron.chanMsgRedirect = chanMsg
}

func (customCron *SrtTranslator) jobTranslator() {
	if customCron.numCoroutine <= 0 {
		customCron.log().WarnF("%s-%s通道的最大数量(%d)无效, 重置为5",
			cronName, "chanData", customCron.numCoroutine)
		customCron.numCoroutine = 5
	}

	for idx := 0; idx < customCron.numCoroutine; idx++ {
		go func(ctx context.Context, chanTranslator chan *SrtTranslateData, chanMsg chan string, idx int) {
			coroutineName := fmt.Sprintf("翻译协程(%d)", idx)
			nextCron := writer.Singleton()
			for true {
				select {
				case <-ctx.Done():

					customCron.log().WarnF("%s关闭(ctx.done), %s被迫退出", cronName, coroutineName)
					runtime.Goexit()
				case currentData, isOpen := <-chanTranslator:
					timeStart := carbon.Now()
					if isOpen == false && currentData == nil {
						customCron.log().WarnF("%s-通道关闭, %s被迫退出", cronName, coroutineName)
						runtime.Goexit()
					}

					var blockChunked []string
					tmpBlockStr := ""
					for _, block := range currentData.PrtSrt.Blocks {
						if block.SubTrack != "" && currentData.PtrOpts.TranslateMode == _type.ModeDelta {
							continue
						}

						if len(tmpBlockStr) >= currentData.PtrOpts.Translator.GetTextMaxLen() {
							blockChunked = append(blockChunked, tmpBlockStr)
							tmpBlockStr = ""
						}

						if len(block.MainTrack) >= currentData.PtrOpts.Translator.GetTextMaxLen() {
							if tmpBlockStr != "" {
								blockChunked = append(blockChunked, tmpBlockStr)
								tmpBlockStr = ""
							}
							blockChunked = append(blockChunked, block.MainTrack)
							continue
						}

						if len(block.MainTrack)+len(tmpBlockStr) >= currentData.PtrOpts.Translator.GetTextMaxLen() {
							blockChunked = append(blockChunked, tmpBlockStr)
							tmpBlockStr = ""
						}

						if tmpBlockStr == "" {
							tmpBlockStr = block.MainTrack
						} else {
							tmpBlockStr = fmt.Sprintf("%s%s%s", tmpBlockStr, currentData.PtrOpts.Translator.GetSep(), block.MainTrack)
						}
					}
					if tmpBlockStr != "" {
						blockChunked = append(blockChunked, tmpBlockStr)
						tmpBlockStr = ""
					}

					if len(blockChunked) == 0 {
						chanMsg <- fmt.Sprintf("字幕文件(%s)未解析到需要翻译的字幕块, 疑似增量翻译模式", currentData.PrtSrt.FileName)
						continue
					}
					flagTranslated := false
					for _, currentBlock := range blockChunked {
						translateRes, err := currentData.PtrOpts.Translator.Translate(&translator.TranslateArgs{
							FromLang:    currentData.PtrOpts.FromLang,
							ToLang:      currentData.PtrOpts.ToLang,
							TextContent: currentBlock,
						})
						if err != nil {
							msg := fmt.Sprintf("翻译异常, 引擎: %s, %s", currentData.PtrOpts.Translator.GetName(), err)
							customCron.log().Error(msg)
							chanMsg <- msg
							break
						}
						for _, result := range translateRes.Results {
							for blockIdx, block := range currentData.PrtSrt.Blocks {
								if result.Id == block.MainTrack {
									if flagTranslated == false {
										flagTranslated = true
									}
									currentData.PrtSrt.Blocks[blockIdx].SubTrack = result.TextTranslated
								}
							}
						}
					}
					if flagTranslated == false {
						chanMsg <- fmt.Sprintf("字幕文件(%s)未进行翻译", currentData.PrtSrt.FileName)
						continue
					}
					currentData.PrtSrt.FlagTranslated = flagTranslated

					nextCron.Push(currentData.toSrtWriterData())

					chanMsg <- fmt.Sprintf(
						"字幕文件(%s)翻译完成, 引擎: %s, 耗时(s): %d",
						currentData.PrtSrt.FileName, currentData.PtrOpts.Translator.GetName(), util.GetSecondsFromTime(timeStart),
					)
				}
			}
		}(customCron.ctx, customCron.chanData, customCron.chanMsg, idx)
	}
}

func (customCron *SrtTranslator) jobMsg() {
	cron.FuncSrtCronMsgRedirect(customCron.ctx, cronName, customCron.log(), customCron.chanMsg, customCron.chanMsgRedirect)
}

func (customCron *SrtTranslator) init() {
	customCron.chanData = make(chan *SrtTranslateData, numChanData)
	customCron.chanMsg = make(chan string, numChanMsg)
	customCron.numCoroutine = numCoroutine
}

func (customCron *SrtTranslator) log() *log.Log {
	return log.Singleton()
}
