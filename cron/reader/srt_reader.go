package reader

import (
	"anto/cron"
	"anto/cron/translate"
	"anto/lib/log"
	"anto/lib/srt"
	"anto/lib/util"
	"bytes"
	"context"
	"fmt"
	"github.com/golang-module/carbon"
	"io"
	"os"
	"runtime"
	"sync"
)

const (
	cronName     = "SRT读取程序"
	numChanData  = 10
	numChanMsg   = 20
	numCoroutine = 10
)

var (
	apiSingleton  *SrtReader
	onceSingleton sync.Once
)

func Singleton() *SrtReader {
	onceSingleton.Do(func() {
		apiSingleton = new(SrtReader)
		apiSingleton.init()
	})
	return apiSingleton
}

type SrtReader struct {
	ctx             context.Context
	ctxFnCancel     context.CancelFunc
	chanData        chan *SrtReaderData
	chanMsg         chan string
	chanMsgRedirect chan string
	numCoroutine    int
}

func (customCron *SrtReader) Push(data *SrtReaderData) {
	customCron.chanData <- data
}

func (customCron *SrtReader) Run(ctx context.Context, fnCancel context.CancelFunc) {
	customCron.ctx = ctx
	customCron.ctxFnCancel = fnCancel

	customCron.jobReader()
	customCron.jobMsg()
}

func (customCron *SrtReader) Close() {}

func (customCron *SrtReader) SetMsgRedirect(chanMsg chan string) {
	customCron.chanMsgRedirect = chanMsg
}

func (customCron *SrtReader) jobReader() {
	if customCron.numCoroutine <= 0 {
		customCron.log().WarnF("%s-%s通道的最大数量(%d)无效, 重置为5", cronName, "chanData", customCron.numCoroutine)
		customCron.numCoroutine = 5
	}

	for idx := 0; idx < customCron.numCoroutine; idx++ {
		go func(ctx context.Context, chanReader chan *SrtReaderData, chanMsg chan string, idx int) {
			coroutineName := fmt.Sprintf("读取协程(%d)", idx)
			for true {
				select {
				case <-ctx.Done():
					customCron.log().WarnF("%s关闭(ctx.done), %s被迫退出", cronName, coroutineName)
					runtime.Goexit()
				case currentData, isOpen := <-chanReader:
					timeStart := carbon.Now()
					if isOpen == false && currentData == nil {
						customCron.log().WarnF("%s-通道关闭, %s被迫退出", cronName, coroutineName)
						runtime.Goexit()
					}

					if currentData.FilePath == "" || currentData.PtrTranslatorOpts == nil {
						chanMsg <- "当前数据包无效, 即将丢弃"
						continue
					}
					fileFd, err := os.Open(currentData.FilePath)
					if err != nil {
						msg := fmt.Sprintf("打开文件(%s)异常, 错误: %s, 即将丢弃", currentData.FilePath, err)
						customCron.log().Error(msg)
						chanMsg <- msg
						continue
					}
					bytesRead, err := io.ReadAll(fileFd)
					if err != nil {
						_ = fileFd.Close()
						msg := fmt.Sprintf("读取文件(%s)异常, 错误: %s, 即将丢弃", currentData.FilePath, err)
						customCron.log().Error(msg)
						chanMsg <- msg
						continue
					}
					_ = fileFd.Close()

					currentData.PrtSrt = new(srt.Srt)
					currentData.PrtSrt.FilePath = currentData.FilePath
					currentData.PrtSrt.FileSize = len(bytesRead)
					currentData.PrtSrt.FileNameSync()

					if err = currentData.PrtSrt.Decode(bytes.NewReader(bytesRead)); err != nil {
						msg := fmt.Sprintf("解析文件(%s)异常, 错误: %s, 即将丢弃", currentData.FilePath, err)
						customCron.log().Error(msg)
						chanMsg <- msg
						continue
					}

					currentData.PrtSrt.CntBlock = len(currentData.PrtSrt.Blocks)
					translate.Singleton().Push(currentData.toTranslateData())

					chanMsg <- fmt.Sprintf(
						"读取文件(%s)成功, 文件名: %s, 字幕块: %d, 文件大小: %d, 耗时: %d",
						currentData.FilePath, currentData.PrtSrt.FileName, currentData.PrtSrt.CntBlock,
						currentData.PrtSrt.FileSize, util.GetSecondsFromTime(timeStart),
					)
				}
			}
		}(customCron.ctx, customCron.chanData, customCron.chanMsg, idx)
	}
}

func (customCron *SrtReader) jobMsg() {
	cron.FuncSrtCronMsgRedirect(customCron.ctx, cronName, customCron.log(), customCron.chanMsg, customCron.chanMsgRedirect)
}

func (customCron *SrtReader) init() {
	customCron.chanData = make(chan *SrtReaderData, numChanData)
	customCron.chanMsg = make(chan string, numChanMsg)
	customCron.numCoroutine = numCoroutine
}

func (customCron *SrtReader) log() *log.Log {
	return log.Singleton()
}
