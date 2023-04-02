package reader

import (
	"bytes"
	"context"
	"fmt"
	"github.com/golang-module/carbon"
	"io"
	"os"
	"runtime"
	"sync"
	"translator/cron/translator"
	"translator/tst/tt_log"
	"translator/tst/tt_srt"
)

var (
	apiSrtReader  *SrtReader
	onceSrtReader sync.Once
)

func GetInstance() *SrtReader {
	onceSrtReader.Do(func() {
		apiSrtReader = new(SrtReader)
		apiSrtReader.init()
	})
	return apiSrtReader
}

type SrtReaderData struct {
	FilePath          string
	PrtSrt            *tt_srt.Srt
	PtrTranslatorOpts *translator.SrtTranslatorOpts
}

func (customSRD *SrtReaderData) toTranslatorData() *translator.SrtTranslatorData {
	return &translator.SrtTranslatorData{
		PrtSrt:  customSRD.PrtSrt,
		PtrOpts: customSRD.PtrTranslatorOpts,
	}
}

type SrtReader struct {
	ctx             context.Context
	chanReader      chan *SrtReaderData
	chanMsgReader   chan string
	chanMsgRedirect chan string
	maxChanReader   int
}

func (customSR *SrtReader) SetMsgRedirect(chanMsg chan string) {
	customSR.chanMsgRedirect = chanMsg
}

func (customSR *SrtReader) Run(ctx context.Context) {
	customSR.ctx = ctx
	customSR.jobReader()
	customSR.jobMsg()
}

func (customSR *SrtReader) Push(data *SrtReaderData) {
	customSR.chanReader <- data
}

func (customSR *SrtReader) jobReader() {
	if customSR.maxChanReader <= 0 {
		customSR.log().Warn(fmt.Sprintf("%s-%s通道的最大数量(%d)无效, 重置为5", customSR.getName(), "chanReader", customSR.maxChanReader))
		customSR.maxChanReader = 5
	}

	for idx := 0; idx < customSR.maxChanReader; idx++ {
		go func(ctx context.Context, chanReader chan *SrtReaderData, chanMsg chan string, idx int) {
			coroutineName := fmt.Sprintf("读取协程(%d)", idx)
			chanName := "chanReader"
			for true {
				select {
				case <-ctx.Done():
					customSR.log().Info(fmt.Sprintf("%s关闭(ctx.done), %s被迫退出", customSR.getName(), coroutineName))
					runtime.Goexit()
				case currentData, isOpen := <-chanReader:
					timeStart := carbon.Now()
					if isOpen == false && currentData == nil {
						customSR.log().Info(fmt.Sprintf("%s-%s通道关闭, %s被迫退出", customSR.getName(), chanName, coroutineName))
						runtime.Goexit()
					}

					if currentData.FilePath == "" || currentData.PtrTranslatorOpts == nil {
						chanMsg <- fmt.Sprintf("当前数据包无效, 即将丢弃")
						continue
					}
					fileFd, err := os.Open(currentData.FilePath)
					defer func() { // @todo wait to fix
						if fileFd != nil {
							_ = fileFd.Close()
						}
					}()
					if err != nil {
						chanMsg <- fmt.Sprintf("打开文件(%s)异常, 错误: %s, 即将丢弃", currentData.FilePath, err)
						continue
					}
					bytesRead, err := io.ReadAll(fileFd)
					if err != nil {
						chanMsg <- fmt.Sprintf("读取文件(%s)异常, 错误: %s, 即将丢弃", currentData.FilePath, err)
						continue
					}

					currentData.PrtSrt = new(tt_srt.Srt)
					currentData.PrtSrt.FilePath = currentData.FilePath
					currentData.PrtSrt.FileSize = len(bytesRead)
					currentData.PrtSrt.FileNameSync()

					if err := currentData.PrtSrt.Decode(bytes.NewReader(bytesRead)); err != nil {
						chanMsg <- fmt.Sprintf("解析文件(%s)异常, 错误: %s, 即将丢弃", currentData.FilePath, err)
						continue
					}
					currentData.PrtSrt.CntBlock = len(currentData.PrtSrt.Blocks)
					translator.GetInstance().Push(currentData.toTranslatorData())
					chanMsg <- fmt.Sprintf(
						"读取文件(%s)成功, 文件名: %s, 字幕块: %d, 文件大小: %d, 耗时: %d",
						currentData.FilePath, currentData.PrtSrt.FileName, currentData.PrtSrt.CntBlock,
						currentData.PrtSrt.FileSize, carbon.Now().DiffAbsInSeconds(timeStart),
					)
				}
			}
		}(customSR.ctx, customSR.chanReader, customSR.chanMsgReader, idx)
	}
}

func (customSR *SrtReader) jobMsg() {
	go func(ctx context.Context, chanMsgReader, chanMsgRedirect chan string) {
		coroutineName := "消息协程"
		chanName := "chanMsgReader"

		for true {
			select {
			case <-ctx.Done():
				customSR.log().Info(fmt.Sprintf("%s关闭(ctx.done), %s被迫退出", customSR.getName(), coroutineName))
				runtime.Goexit()
			case currentMsg, isOpen := <-chanMsgReader:
				if isOpen == false && currentMsg == "" {
					customSR.log().Info(fmt.Sprintf("%s-%s通道关闭, %s被迫退出", customSR.getName(), chanName, coroutineName))
					runtime.Goexit()
				}
				if chanMsgRedirect != nil {
					chanMsgRedirect <- fmt.Sprintf("当前时间: %s, 来源: %s, 信息: [%s]", carbon.Now().Layout(carbon.ShortDateTimeLayout), customSR.getName(), currentMsg)
				}
				customSR.log().Info(fmt.Sprintf("来源: %s, 信息: %s", customSR.getName(), currentMsg))
			}
		}
	}(customSR.ctx, customSR.chanMsgReader, customSR.chanMsgRedirect)
}

func (customSR *SrtReader) getName() string {
	return "SRT读取程序"
}

func (customSR *SrtReader) init() {
	customSR.chanReader = make(chan *SrtReaderData, 10)
	customSR.chanMsgReader = make(chan string, 20)
	customSR.maxChanReader = 10
}

func (customSR *SrtReader) log() *tt_log.TTLog {
	return tt_log.GetInstance()
}
