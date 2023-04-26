// Package detector 检测器
package detector

import (
	_type "anto/common"
	"anto/cron"
	"anto/cron/reader"
	"anto/lib/log"
	"anto/lib/util"
	"context"
	"fmt"
	"github.com/golang-module/carbon"
	"io/fs"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

const (
	cronName     = "SRT检测程序"
	numChanData  = 10
	numChanMsg   = 20
	numCoroutine = 5
)

var (
	apiSingleton  *SrtDetector
	onceSingleton sync.Once
)

func Singleton() *SrtDetector {
	onceSingleton.Do(func() {
		apiSingleton = new(SrtDetector)
		apiSingleton.init()
	})
	return apiSingleton
}

type SrtDetector struct {
	ctx             context.Context
	ctxFnCancel     context.CancelFunc
	chanData        chan *StrDetectorData
	chanMsg         chan string
	chanMsgRedirect chan string
	numCoroutine    int
}

func (customCron *SrtDetector) Push(data *StrDetectorData) {
	customCron.chanData <- data
}

func (customCron *SrtDetector) Run(ctx context.Context, fnCancel context.CancelFunc) {
	customCron.ctx = ctx
	customCron.ctxFnCancel = fnCancel

	customCron.jobDetector()
	customCron.jobMsg()
}

func (customCron *SrtDetector) Close() {}

func (customCron *SrtDetector) SetMsgRedirect(chanMsg chan string) {
	customCron.chanMsgRedirect = chanMsg
}

func (customCron *SrtDetector) jobDetector() {
	if customCron.numCoroutine <= 0 {
		customCron.log().WarnF("%s-%s通道的最大数量(%d)无效, 重置为5", cronName, "chanData", customCron.numCoroutine)
		customCron.numCoroutine = 5
	}

	for idx := 0; idx < customCron.numCoroutine; idx++ {
		go func(ctx context.Context, chanDetector chan *StrDetectorData, chanMsg chan string, idx int) {
			coroutineName := fmt.Sprintf("检测协程(%d)", idx)
			for true {
				select {
				case <-ctx.Done():
					customCron.log().WarnF("%s关闭(ctx.done), %s被迫退出", cronName, coroutineName)
					runtime.Goexit()
				case currentData, isOpen := <-chanDetector:
					timeStart := carbon.Now()
					if isOpen == false && currentData == nil {
						customCron.log().WarnF("%s-通道关闭, %s被迫退出", cronName, coroutineName)
						runtime.Goexit()
					}
					if currentData.SrtFile != "" {
						if len(currentData.SrtFile) > 4 && currentData.SrtFile[len(currentData.SrtFile)-4:] == ".srt" {
							reader.Singleton().Push(currentData.toReaderData(currentData.SrtFile))
							chanMsg <- fmt.Sprintf("检测到文件(%s), 耗时(s): %d", currentData.SrtFile, util.GetSecondsFromTime(timeStart))
						}
					}
					if currentData.SrtDir == "" {
						continue
					}
					_ = filepath.Walk(currentData.SrtDir, func(path string, info fs.FileInfo, err error) error {
						if info.IsDir() {
							if path != currentData.SrtDir {
								return filepath.SkipDir
							}
							return nil
						}
						if !util.IsSrtFile(path) || strings.Contains(path, _type.AppName) || info.Size() == 0 {
							return nil
						}
						reader.Singleton().Push(currentData.toReaderData(path))
						chanMsg <- fmt.Sprintf("检测到文件(%s), 耗时(s): %d", path, util.GetSecondsFromTime(timeStart))
						return nil
					})
				}
			}
		}(customCron.ctx, customCron.chanData, customCron.chanMsg, idx)
	}
}

func (customCron *SrtDetector) jobMsg() {
	cron.FuncSrtCronMsgRedirect(customCron.ctx, cronName, customCron.log(), customCron.chanMsg, customCron.chanMsgRedirect)
}

func (customCron *SrtDetector) init() {
	customCron.chanData = make(chan *StrDetectorData, numChanData)
	customCron.chanMsg = make(chan string, numChanMsg)
	customCron.numCoroutine = numCoroutine
}

func (customCron *SrtDetector) log() *log.Log {
	return log.Singleton()
}
