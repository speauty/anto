// Package detector 检测器
package detector

import (
	"context"
	"fmt"
	"github.com/golang-module/carbon"
	"io/fs"
	"path/filepath"
	"runtime"
	"sync"
	"translator/cron/reader"
	"translator/cron/translator"
	"translator/tst/tt_log"
	"translator/tst/tt_translator"
	_type "translator/type"
	"translator/util"
)

var (
	apiSrtDetector  *SrtDetector
	onceSrtDetector sync.Once
)

func GetInstance() *SrtDetector {
	onceSrtDetector.Do(func() {
		apiSrtDetector = new(SrtDetector)
		apiSrtDetector.init()
	})
	return apiSrtDetector
}

type StrDetectorData struct {
	Translator      tt_translator.ITranslator
	FromLang        string
	ToLang          string
	TranslateMode   _type.TranslateMode
	MainTrackReport _type.LangDirection
	SrtFile         string
	SrtDir          string
	FlagTrackExport int
}

func (customSDD StrDetectorData) toReaderData(filePath string) *reader.SrtReaderData {
	return &reader.SrtReaderData{
		FilePath: filePath,
		PtrTranslatorOpts: &translator.SrtTranslatorOpts{
			Translator:      customSDD.Translator,
			FromLang:        customSDD.FromLang,
			ToLang:          customSDD.ToLang,
			TranslateMode:   customSDD.TranslateMode,
			MainTrackReport: customSDD.MainTrackReport,
			FlagTrackExport: customSDD.FlagTrackExport,
		},
	}
}

type SrtDetector struct {
	ctx             context.Context
	chanDetector    chan *StrDetectorData
	chanMsgDetector chan string
	chanMsgRedirect chan string
	maxChanDetector int
}

func (customSD *SrtDetector) SetMsgRedirect(chanMsg chan string) {
	customSD.chanMsgRedirect = chanMsg
}

func (customSD *SrtDetector) Push(data *StrDetectorData) {
	customSD.chanDetector <- data
}

func (customSD *SrtDetector) Run(ctx context.Context) {
	customSD.ctx = ctx
	customSD.jobDetector()
	customSD.jobMsg()

}

func (customSD *SrtDetector) jobDetector() {
	if customSD.maxChanDetector <= 0 {
		customSD.log().Warn(fmt.Sprintf("%s-%s通道的最大数量(%d)无效, 重置为5", customSD.getName(), "chanDetector", customSD.maxChanDetector))
		customSD.maxChanDetector = 5
	}

	for idx := 0; idx < customSD.maxChanDetector; idx++ {
		go func(ctx context.Context, chanDetector chan *StrDetectorData, chanMsg chan string, idx int) {
			coroutineName := fmt.Sprintf("检测协程(%d)", idx)
			chanName := "chanDetector"
			for true {
				select {
				case <-ctx.Done():
					customSD.log().Info(fmt.Sprintf("%s关闭(ctx.done), %s被迫退出", customSD.getName(), coroutineName))
					runtime.Goexit()
				case currentData, isOpen := <-chanDetector:
					timeStart := carbon.Now()
					if isOpen == false && currentData == nil {
						customSD.log().Info(fmt.Sprintf("%s-%s通道关闭, %s被迫退出", customSD.getName(), chanName, coroutineName))
						runtime.Goexit()
					}
					if currentData.SrtFile != "" {
						if len(currentData.SrtFile) > 4 && currentData.SrtFile[len(currentData.SrtFile)-4:] == ".srt" {
							reader.GetInstance().Push(currentData.toReaderData(currentData.SrtFile))
							chanMsg <- fmt.Sprintf("检测到文件: %s, 耗时: %d", currentData.SrtFile, carbon.Now().DiffAbsInSeconds(timeStart))
						}
					}
					if currentData.SrtDir != "" {
						_ = filepath.Walk(currentData.SrtDir, func(path string, info fs.FileInfo, err error) error {
							if info.IsDir() || !util.IsSrtFile(path) || info.Size() == 0 { // 过掉目录 或非srt文件
								return nil
							}
							if len(path) > 4 && path[len(path)-4:] == ".srt" {
								reader.GetInstance().Push(currentData.toReaderData(path))
								chanMsg <- fmt.Sprintf("检测到文件: %s, 耗时: %d", path, carbon.Now().DiffAbsInSeconds(timeStart))
							}
							return nil
						})
					}
				}
			}
		}(customSD.ctx, customSD.chanDetector, customSD.chanMsgDetector, idx)
	}
}

func (customSD *SrtDetector) jobMsg() {
	go func(ctx context.Context, chanMsgDetector, chanMsgRedirect chan string) {
		coroutineName := "消息协程"
		chanName := "chanMsgDetector"

		for true {
			select {
			case <-ctx.Done():
				customSD.log().Info(fmt.Sprintf("%s关闭(ctx.done), %s被迫退出", customSD.getName(), coroutineName))
				runtime.Goexit()
			case currentMsg, isOpen := <-chanMsgDetector:
				if isOpen == false && currentMsg == "" {
					customSD.log().Info(fmt.Sprintf("%s-%s通道关闭, %s被迫退出", customSD.getName(), chanName, coroutineName))
					runtime.Goexit()
				}
				if chanMsgRedirect != nil {
					chanMsgRedirect <- fmt.Sprintf("当前时间: %s, 来源: %s, 信息: [%s]", carbon.Now().Layout(carbon.ShortDateTimeLayout), customSD.getName(), currentMsg)
				}
				customSD.log().Info(fmt.Sprintf("来源: %s, 信息: %s", customSD.getName(), currentMsg))
			}
		}
	}(customSD.ctx, customSD.chanMsgDetector, customSD.chanMsgRedirect)
}

func (customSD *SrtDetector) getName() string {
	return "SRT检测程序"
}

func (customSD *SrtDetector) init() {
	customSD.chanDetector = make(chan *StrDetectorData, 10)
	customSD.chanMsgDetector = make(chan string, 20)
	customSD.maxChanDetector = 10
}

func (customSD *SrtDetector) log() *tt_log.TTLog {
	return tt_log.GetInstance()
}
