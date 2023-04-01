package writer

import (
	"context"
	"fmt"
	"github.com/golang-module/carbon"
	"go.uber.org/zap"
	"os"
	"runtime"
	"sync"
	"translator/tst/tt_log"
	"translator/tst/tt_srt"
)

var (
	apiSrtWriter  *SrtWriter
	onceSrtWriter sync.Once
)

func GetInstance() *SrtWriter {
	onceSrtWriter.Do(func() {
		apiSrtWriter = new(SrtWriter)
		apiSrtWriter.init()
	})
	return apiSrtWriter
}

type SrtWriterData struct {
	FileNameSaved string
	PrtSrt        *tt_srt.Srt
	PtrOpts       *tt_srt.EncodeOpt
}

type SrtWriter struct {
	ctx           context.Context
	chanWriter    chan *SrtWriterData
	chanMsgWriter chan string
	maxChanWriter int
}

func (customSW *SrtWriter) Run(ctx context.Context) {
	customSW.ctx = ctx
}

func (customSW *SrtWriter) Push(data *SrtWriterData) {
	customSW.chanWriter <- data
}

func (customSW *SrtWriter) jobWriter() {
	if customSW.maxChanWriter <= 0 {
		customSW.log().Warn(fmt.Sprintf("%s-%s通道的最大数量(%d)无效, 重置为5", customSW.getName(), "chanWriter", customSW.maxChanWriter))
		customSW.maxChanWriter = 5
	}

	for idx := 0; idx < customSW.maxChanWriter; idx++ {
		go func(ctx context.Context, chanWriter chan *SrtWriterData, chanMsg chan string, idx int) {
			coroutineName := fmt.Sprintf("写入协程(%d)", idx)
			chanName := "chanWriter"
			for true {
				select {
				case <-ctx.Done():
					customSW.log().Info(fmt.Sprintf("%s关闭(ctx.done), %s被迫退出", customSW.getName(), coroutineName))
					runtime.Goexit()
				case currentData, isOpen := <-chanWriter:
					timeStart := carbon.Now()
					if isOpen == false && currentData == nil {
						customSW.log().Info(fmt.Sprintf("%s-%s通道关闭, %s被迫退出", customSW.getName(), chanName, coroutineName))
						runtime.Goexit()
					}
					if currentData.FileNameSaved == "" {
						chanMsg <- fmt.Sprintf("%s未检测到保存文件名称", currentData.PrtSrt.FileName)
						continue
					}
					if currentData.PrtSrt.FlagTranslated == false {
						chanMsg <- fmt.Sprintf("%s未进行翻译, 可能选择了增量翻译模式, 源文已全部翻译", currentData.PrtSrt.FileName)
						continue
					}
					bytesEncoded, err := currentData.PrtSrt.Encode(currentData.PtrOpts)
					if err != nil {
						chanMsg <- fmt.Sprintf("%s编码失败", currentData.PrtSrt.FileName)
						continue
					}

					if err = os.WriteFile(currentData.FileNameSaved, bytesEncoded, os.ModePerm); err != nil {
						chanMsg <- fmt.Sprintf("写入文件失败(%s => %s)", currentData.PrtSrt.FileName, currentData.FileNameSaved)
						continue
					}
					chanMsg <- fmt.Sprintf(
						"写入文件成功, 源件: %s, 目标文件: %s, 写入字节数: %d, 耗时(s): %d",
						currentData.PrtSrt.FileName, currentData.FileNameSaved, len(bytesEncoded),
						carbon.Now().DiffAbsInSeconds(timeStart),
					)
				}
			}
		}(customSW.ctx, customSW.chanWriter, customSW.chanMsgWriter, idx)
	}
}

func (customSW *SrtWriter) RedirectMsgTo(targetChan chan string) {
	go func(ctx context.Context, chanMsgWriter, targetChan chan string) {
		coroutineName := "消息定向协程"
		chanName := "chanWriter"

		for true {
			select {
			case <-ctx.Done():
				customSW.log().Info(fmt.Sprintf("%s关闭(ctx.done), %s被迫退出", customSW.getName(), coroutineName))
				runtime.Goexit()
			case currentMsg, isOpen := <-chanMsgWriter:
				if isOpen == false && currentMsg == "" {
					customSW.log().Info(fmt.Sprintf("%s-%s通道关闭, %s被迫退出", customSW.getName(), chanName, coroutineName))
					runtime.Goexit()
				}
				if targetChan == nil {
					customSW.log().Info(fmt.Sprintf("%s未设置通道(%s)接管, 定向输出到日志", customSW.getName(), chanName), zap.String("msg", currentMsg))
				}
				targetChan <- fmt.Sprintf("当前时间: %s, 来源: %s, 信息: [%s]", carbon.Now().Layout(carbon.ShortDateTimeLayout), customSW.getName(), currentMsg)
			}
		}
	}(customSW.ctx, customSW.chanMsgWriter, targetChan)
}

func (customSW *SrtWriter) getName() string {
	return "SRT写入程序"
}

func (customSW *SrtWriter) init() {
	customSW.chanWriter = make(chan *SrtWriterData, 10)
	customSW.chanMsgWriter = make(chan string, 20)
	customSW.maxChanWriter = 10
}

func (customSW *SrtWriter) log() *tt_log.TTLog {
	return tt_log.GetInstance()
}
