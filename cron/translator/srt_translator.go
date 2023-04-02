package translator

import (
	"context"
	"fmt"
	"github.com/golang-module/carbon"
	"runtime"
	"sync"
	"translator/cron/writer"
	"translator/tst/tt_log"
	"translator/tst/tt_srt"
	"translator/tst/tt_translator"
	_type "translator/type"
)

var (
	apiSrtTranslator  *SrtTranslator
	onceSrtTranslator sync.Once
)

func GetInstance() *SrtTranslator {
	onceSrtTranslator.Do(func() {
		apiSrtTranslator = new(SrtTranslator)
		apiSrtTranslator.init()
	})
	return apiSrtTranslator
}

type SrtTranslatorData struct {
	PrtSrt  *tt_srt.Srt
	PtrOpts *SrtTranslatorOpts
}

func (customSTD *SrtTranslatorData) toSrtWriterData() *writer.SrtWriterData {
	tmpData := &writer.SrtWriterData{
		FileNameSaved: customSTD.fileNameSavedBuilder(),
		PrtSrt:        customSTD.PrtSrt,
		PtrOpts:       &tt_srt.EncodeOpt{FlagIsInverse: customSTD.PtrOpts.MainTrackReport == _type.LangDirectionTo},
	}
	return tmpData
}

func (customSTD *SrtTranslatorData) fileNameSavedBuilder() string {
	newFileName := customSTD.PrtSrt.FilePath[0 : len(customSTD.PrtSrt.FilePath)-4]
	newFileName = fmt.Sprintf(
		"%s.%s2%s.srt", newFileName, customSTD.PtrOpts.FromLang, customSTD.PtrOpts.ToLang,
	)
	return newFileName
}

type SrtTranslatorOpts struct {
	Translator      tt_translator.ITranslator
	FromLang        string
	ToLang          string
	TranslateMode   _type.TranslateMode
	MainTrackReport _type.LangDirection
}

type SrtTranslator struct {
	ctx               context.Context
	chanTranslator    chan *SrtTranslatorData
	chanMsgTranslator chan string
	chanMsgRedirect   chan string
	maxChanTranslator int
}

func (customST *SrtTranslator) SetMsgRedirect(chanMsg chan string) {
	customST.chanMsgRedirect = chanMsg
}

func (customST *SrtTranslator) Run(ctx context.Context) {
	customST.ctx = ctx
	customST.jobTranslator()
	customST.jobMsg()
}

func (customST *SrtTranslator) Push(data *SrtTranslatorData) {
	customST.chanTranslator <- data
}

func (customST *SrtTranslator) jobTranslator() {
	if customST.maxChanTranslator <= 0 {
		customST.log().Warn(fmt.Sprintf("%s-%s通道的最大数量(%d)无效, 重置为5",
			customST.getName(), "chanTranslator", customST.maxChanTranslator))
		customST.maxChanTranslator = 5
	}

	for idx := 0; idx < customST.maxChanTranslator; idx++ {
		go func(ctx context.Context, chanTranslator chan *SrtTranslatorData, chanMsg chan string, idx int) {
			coroutineName := fmt.Sprintf("翻译协程(%d)", idx)
			chanName := "chanTranslator"
			nextCron := writer.GetInstance()
			for true {
				select {
				case <-ctx.Done():
					customST.log().Info(fmt.Sprintf("%s关闭(ctx.done), %s被迫退出", customST.getName(), coroutineName))
					runtime.Goexit()
				case currentData, isOpen := <-chanTranslator:
					timeStart := carbon.Now()
					if isOpen == false && currentData == nil {
						customST.log().Info(fmt.Sprintf("%s-%s通道关闭, %s被迫退出", customST.getName(), chanName, coroutineName))
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
						translateRes, err := currentData.PtrOpts.Translator.Translate(&tt_translator.TranslateArgs{
							FromLang:    currentData.PtrOpts.FromLang,
							ToLang:      currentData.PtrOpts.ToLang,
							TextContent: currentBlock,
						})
						if err != nil {
							chanMsg <- fmt.Sprintf("翻译异常, 引擎: %s", currentData.PtrOpts.Translator.GetName())
							continue
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
					currentData.PrtSrt.FlagTranslated = flagTranslated
					if flagTranslated == false {
						chanMsg <- fmt.Sprintf("字幕文件(%s)未进行翻译", currentData.PrtSrt.FileName)
						continue
					}
					chanMsg <- fmt.Sprintf(
						"字幕文件(%s)翻译完成, 引擎: %s, 耗时(s): %d",
						currentData.PrtSrt.FileName, currentData.PtrOpts.Translator.GetName(), carbon.Now().DiffAbsInSeconds(timeStart),
					)
					nextCron.Push(currentData.toSrtWriterData())
				}
			}
		}(customST.ctx, customST.chanTranslator, customST.chanMsgTranslator, idx)
	}
}

func (customST *SrtTranslator) jobMsg() {
	go func(ctx context.Context, chanMsgTranslator, chanMsgRedirect chan string) {
		coroutineName := "消息协程"
		chanName := "chanMsgTranslator"

		for true {
			select {
			case <-ctx.Done():
				customST.log().Info(fmt.Sprintf("%s关闭(ctx.done), %s被迫退出", customST.getName(), coroutineName))
				runtime.Goexit()
			case currentMsg, isOpen := <-chanMsgTranslator:
				if isOpen == false && currentMsg == "" {
					customST.log().Info(fmt.Sprintf("%s-%s通道关闭, %s被迫退出", customST.getName(), chanName, coroutineName))
					runtime.Goexit()
				}
				if chanMsgRedirect != nil {
					chanMsgRedirect <- fmt.Sprintf("当前时间: %s, 来源: %s, 信息: [%s]", carbon.Now().Layout(carbon.ShortDateTimeLayout), customST.getName(), currentMsg)
				}
				customST.log().Info(fmt.Sprintf("来源: %s, 信息: %s", customST.getName(), currentMsg))
			}
		}
	}(customST.ctx, customST.chanMsgTranslator, customST.chanMsgRedirect)
}

func (customST *SrtTranslator) getName() string {
	return "SRT翻译程序"
}

func (customST *SrtTranslator) init() {
	customST.chanTranslator = make(chan *SrtTranslatorData, 10)
	customST.chanMsgTranslator = make(chan string, 20)
	customST.maxChanTranslator = 10
}

func (customST *SrtTranslator) log() *tt_log.TTLog {
	return tt_log.GetInstance()
}
