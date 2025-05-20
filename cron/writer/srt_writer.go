package writer

import (
	"anto/cron"
	"anto/lib/log"
	"anto/lib/util"
	"context"
	"fmt"
	"github.com/golang-module/carbon"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

const (
	cronName     = "SRT写入程序"
	numChanData  = 10
	numChanMsg   = 20
	numCoroutine = 10
)

var (
	apiSingleton  *SrtWriter
	onceSingleton sync.Once
)

func Singleton() *SrtWriter {
	onceSingleton.Do(func() {
		apiSingleton = new(SrtWriter)
		apiSingleton.init()
	})
	return apiSingleton
}

type SrtWriter struct {
	ctx             context.Context
	ctxFnCancel     context.CancelFunc
	chanData        chan *SrtWriterData
	chanMsg         chan string
	chanMsgRedirect chan string
	numCoroutine    int
}

func (customCron *SrtWriter) Push(data *SrtWriterData) {
	customCron.chanData <- data
}

func (customCron *SrtWriter) Run(ctx context.Context, fnCancel context.CancelFunc) {
	customCron.ctx = ctx
	customCron.ctxFnCancel = fnCancel

	customCron.jobWriter()
	customCron.jobMsg()
}

func (customCron *SrtWriter) Close() {}

func (customCron *SrtWriter) SetMsgRedirect(chanMsg chan string) {
	customCron.chanMsgRedirect = chanMsg
}

func (customCron *SrtWriter) jobWriter() {
	if customCron.numCoroutine <= 0 {
		customCron.log().WarnF("%s-%s通道的最大数量(%d)无效, 重置为5", cronName, "chanData", customCron.numCoroutine)
		customCron.numCoroutine = 5
	}

	for idx := 0; idx < customCron.numCoroutine; idx++ {
		go func(ctx context.Context, chanWriter chan *SrtWriterData, chanMsg chan string, idx int) {
			coroutineName := fmt.Sprintf("写入协程(%d)", idx)
			for true {
				select {
				case <-ctx.Done():
					customCron.log().WarnF("%s关闭(ctx.done), %s被迫退出", cronName, coroutineName)
					runtime.Goexit()
				case currentData, isOpen := <-chanWriter:
					timeStart := carbon.Now()
					if isOpen == false && currentData == nil {
						customCron.log().WarnF("%s-通道关闭, %s被迫退出", cronName, coroutineName)
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
						customCron.log().ErrorF("%s编码失败, 错误: %s", currentData.PrtSrt.FileName, err)
						chanMsg <- fmt.Sprintf("%s编码失败", currentData.PrtSrt.FileName)
						continue
					}

					dirSrt := filepath.Dir(currentData.FileNameSaved)
					if err = os.MkdirAll(dirSrt, os.ModePerm); err != nil {
						customCron.log().ErrorF("创建目录[%s]失败, 错误: %s", dirSrt, err)
						chanMsg <- fmt.Sprintf("%s创建失败", dirSrt)
						continue
					}

					fd, err := os.OpenFile(currentData.FileNameSaved, os.O_CREATE|os.O_WRONLY, os.ModePerm)
					if err != nil {
						customCron.log().ErrorF("创建文件[%s]失败, 错误: %s", currentData.FileNameSaved, err)
						chanMsg <- fmt.Sprintf("%s创建失败", currentData.FileNameSaved)
						_ = fd.Close()
						continue
					}
					if _, err = fd.Write(bytesEncoded); err != nil {
						customCron.log().ErrorF("写入文件失败(%s => %s), 错误: %s", currentData.PrtSrt.FileName, currentData.FileNameSaved, err)
						chanMsg <- fmt.Sprintf("写入文件失败(%s => %s)", currentData.PrtSrt.FileName, currentData.FileNameSaved)
						_ = fd.Close()
						continue
					}
					_ = fd.Close()
					chanMsg <- fmt.Sprintf(
						"写入文件成功, 源件: %s, 目标文件: %s, 写入字节数: %d, 耗时(s): %d",
						currentData.PrtSrt.FileName, currentData.FileNameSaved, len(bytesEncoded),
						util.GetSecondsFromTime(timeStart),
					)
				}
			}
		}(customCron.ctx, customCron.chanData, customCron.chanMsg, idx)
	}
}

func (customCron *SrtWriter) jobMsg() {
	cron.FuncSrtCronMsgRedirect(customCron.ctx, cronName, customCron.log(), customCron.chanMsg, customCron.chanMsgRedirect)
}

func (customCron *SrtWriter) init() {
	customCron.chanData = make(chan *SrtWriterData, numChanData)
	customCron.chanMsg = make(chan string, numChanMsg)
	customCron.numCoroutine = numCoroutine
}

func (customCron *SrtWriter) log() *log.Log {
	return log.Singleton()
}
