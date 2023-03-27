package task

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-module/carbon"
	"go.uber.org/zap"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"translator/tst/tt_log"
	"translator/tst/tt_srt"
	"translator/tst/tt_translator"
	_type "translator/type"
	"translator/util"
)

type Translate struct {
	translator      tt_translator.ITranslator
	fromLang        string
	toLang          string
	translateMode   _type.TranslateMode
	mainTrackReport _type.LangDirection
	srtFile         string
	srtDir          string

	srtFiles  []string
	resultMsg []string

	chanDecode    chan string
	chanTranslate chan *tt_srt.Srt
	chanEncode    chan *tt_srt.Srt
	chanRes       chan string

	ctx context.Context
}

func (customT *Translate) Run() ([]string, error) {
	customT.buildSrtFiles()
	if len(customT.srtFiles) == 0 {
		return nil, errors.New("未检测到srt文件")
	}

	customT.defaultCtxAndChan()

	chanWG := new(sync.WaitGroup)
	chanWG.Add(5)

	go func() {
		defer func() {
			chanWG.Done()
			close(customT.chanDecode)
		}()
		for _, srtFile := range customT.srtFiles {
			customT.chanDecode <- srtFile
		}
	}()

	{
		customT.jobDecode(chanWG)
		customT.jobTranslate(chanWG)
		customT.jobEncode(chanWG)
		customT.jobRes(chanWG)
	}

	chanWG.Wait()
	return customT.resultMsg, nil
}

func (customT *Translate) buildSrtFiles() {
	timeStart := carbon.Now()
	defer func() {
		duration := carbon.Now().DiffAbsInSeconds(timeStart)
		tt_log.GetInstance().Info(fmt.Sprintf("构建srt文件完成, 记录数量: %d, 耗时(s): %d", len(customT.srtFiles), duration))
	}()
	if customT.srtFile != "" && util.IsSrtFile(customT.srtFile) {
		if err := util.IsFileOrDirExisted(customT.srtFile); err == nil {
			tt_log.GetInstance().Info(fmt.Sprintf("构建srt文件, 检测到: %s", customT.srtFile))
			customT.srtFiles = append(customT.srtFiles, customT.srtFile)
		}
	}
	if customT.srtDir != "" {
		if err := util.IsFileOrDirExisted(customT.srtDir); err != nil {
			return
		}
		_ = filepath.Walk(customT.srtDir, func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() || !util.IsSrtFile(path) || info.Size() == 0 { // 过掉目录 或非srt文件
				return nil
			}
			for _, srtExisted := range customT.srtFiles {
				if srtExisted == path {
					return nil
				}
			}
			tt_log.GetInstance().Info(fmt.Sprintf("构建srt文件, 检测到: %s", path))
			customT.srtFiles = append(customT.srtFiles, path)
			return nil
		})
	}
	return
}

func (customT *Translate) jobTranslate(chanWG *sync.WaitGroup) {
	go func(localCtx context.Context, localChanWG *sync.WaitGroup, localChanTranslate chan *tt_srt.Srt, localChanEncode chan *tt_srt.Srt, currentEngine tt_translator.ITranslator) {
		timeStartGo := carbon.Now()
		cntTranslate := 0
		defer func() {
			tt_log.GetInstance().Info(fmt.Sprintf("翻译srt文件进程退出, 翻译数量: %d, 预期数量: %d,运行时长(s): %d", cntTranslate, len(customT.srtFiles), carbon.Now().DiffAbsInSeconds(timeStartGo)))
			localChanWG.Done()
		}()
		for true {
			select {
			case currentSrt, isOpen := <-localChanTranslate:
				if currentSrt == nil {
					if isOpen {
						close(localChanTranslate)
					}
					close(localChanEncode)
					return
				}
				tt_log.GetInstance().Info(fmt.Sprintf("翻译srt文件, 接收: %s", currentSrt.FilePath))
				timeStart := carbon.Now()
				var blockChunked []string
				tmpStr := ""
				cntBlocksIgnored := 0

				for _, currentBlock := range currentSrt.Blocks {
					if customT.translateMode != _type.ModeFull && currentBlock.SubTrack != "" { // 增量模式过滤
						cntBlocksIgnored++
						continue
					}
					if len(tmpStr) >= currentEngine.GetTextMaxLen() || len(tmpStr)+len(currentBlock.MainTrack) >= currentEngine.GetTextMaxLen() {
						blockChunked = append(blockChunked, tmpStr)
						tmpStr = ""
					}
					if len(currentBlock.MainTrack) >= currentEngine.GetTextMaxLen() {
						blockChunked = append(blockChunked, currentBlock.MainTrack)
						continue
					}
					if tmpStr == "" {
						tmpStr = currentBlock.MainTrack
						continue
					}
					tmpStr = fmt.Sprintf("%s%s%s", tmpStr, currentEngine.GetSep(), currentBlock.MainTrack)
				}
				if tmpStr != "" {
					blockChunked = append(blockChunked, tmpStr)
					tmpStr = ""
				}

				cntGoProc := currentEngine.GetProcMax()
				if cntGoProc > len(blockChunked) {
					cntGoProc = len(blockChunked)
				}
				cntBlocks := new(atomic.Int32)
				cntBlocks.Store(0)
				translateWG := new(sync.WaitGroup)
				translateWG.Add(cntGoProc)
				for i := 0; i < cntGoProc; i++ {
					go func() {
						defer translateWG.Done()
						for _, currentBlock := range blockChunked {
							res, err := currentEngine.Translate(&tt_translator.TranslateArgs{
								FromLang:    customT.fromLang,
								ToLang:      customT.toLang,
								TextContent: currentBlock,
							})
							if err != nil {
								tt_log.GetInstance().Info(fmt.Sprintf("翻译失败, 文件: %s, 错误: %s", currentSrt.FilePath, err))
								continue
							}

							for _, result := range res.Results {
								for _, srtBlock := range currentSrt.Blocks {
									if srtBlock.MainTrack == result.Id {
										srtBlock.SubTrack = result.TextTranslated
										cntBlocks.Add(1)
									}
								}
							}
						}
					}()
				}
				translateWG.Wait()
				currentSrt.FlagTranslated = cntBlocks.Load() > 0
				localChanEncode <- currentSrt
				tt_log.GetInstance().Info(fmt.Sprintf("翻译srt文件完成, 文件: %s, 语种: %s->%s, 总块数: %d, 已忽略: %d, 已翻译: %d, 耗时(s): %d", currentSrt.FilePath, customT.fromLang, customT.toLang, len(currentSrt.Blocks), cntBlocksIgnored, cntBlocks.Load(), carbon.Now().DiffAbsInSeconds(timeStart)))
				cntTranslate++
			}
		}
	}(customT.ctx, chanWG, customT.chanTranslate, customT.chanEncode, customT.translator)
}

func (customT *Translate) jobDecode(chanWG *sync.WaitGroup) {
	go func(localCtx context.Context, localChanWG *sync.WaitGroup, localChanDecode chan string, localChanTranslate chan *tt_srt.Srt, files []string) {
		timeStartGo := carbon.Now()
		cntDecode := 0
		defer func() {
			tt_log.GetInstance().Info(fmt.Sprintf("解码srt文件进程退出, 解码数量: %d, 预期数量: %d,运行时长(s): %d", cntDecode, len(customT.srtFiles), carbon.Now().DiffAbsInSeconds(timeStartGo)))
			localChanWG.Done()
		}()
		for true {
			select {
			case currentFile, isOpen := <-localChanDecode:
				if currentFile == "" {
					if isOpen {
						close(localChanDecode)
					}
					close(localChanTranslate)
					return
				}
				tt_log.GetInstance().Info(fmt.Sprintf("解码srt文件, 接收: %s", currentFile))
				timeStart := carbon.Now()

				fd, err := os.Open(currentFile)
				defer func() {
					if fd != nil {
						_ = fd.Close()
					}
				}()
				if err != nil {
					tt_log.GetInstance().Error("读取文件失败", zap.String("文件", currentFile), zap.Error(err))
					continue
				}
				tmpSrt := new(tt_srt.Srt)
				tmpSrt.FilePath = currentFile
				tmpSrt.FileName = filepath.Base(currentFile)

				if err = tmpSrt.Decode(fd); err != nil {
					tt_log.GetInstance().Error("解码srt文件失败", zap.String("文件", currentFile), zap.Error(err))
					continue
				}
				localChanTranslate <- tmpSrt
				tt_log.GetInstance().Info(fmt.Sprintf("解码srt文件完成, 文件: %s, 耗时(s): %d", currentFile, carbon.Now().DiffAbsInSeconds(timeStart)))
				cntDecode++
			}
		}
	}(customT.ctx, chanWG, customT.chanDecode, customT.chanTranslate, customT.srtFiles)
}

func (customT *Translate) jobEncode(chanWG *sync.WaitGroup) {
	go func(localCtx context.Context, localChanWG *sync.WaitGroup, localChanEncode chan *tt_srt.Srt, localChanRes chan string) {
		timeStartGo := carbon.Now()
		cntEncode := 0
		defer func() {
			tt_log.GetInstance().Info(fmt.Sprintf("编码srt文件进程退出, 编码数量: %d, 预期数量: %d,运行时长(s): %d", cntEncode, len(customT.srtFiles), carbon.Now().DiffAbsInSeconds(timeStartGo)))
			localChanWG.Done()
		}()
		for true {
			select {
			case currentSrt, isOpen := <-localChanEncode:
				if currentSrt == nil {
					if isOpen {
						close(localChanEncode)
					}
					close(localChanRes)
					return
				}
				tt_log.GetInstance().Info(fmt.Sprintf("编码srt文件, 接收: %s", currentSrt.FilePath))
				if !currentSrt.FlagTranslated {
					localChanRes <- fmt.Sprintf("翻译文件[%s]失败", currentSrt.FileName)
					continue
				}
				timeStart := carbon.Now()
				bytes, err := currentSrt.Encode(&tt_srt.EncodeOpt{FlagIsInverse: customT.mainTrackReport == _type.LangDirectionTo})
				if err != nil {
					tt_log.GetInstance().Error("编码文件失败", zap.String("文件", currentSrt.FilePath), zap.Error(err))
					continue
				}
				newFileName := currentSrt.FilePath[0 : len(currentSrt.FilePath)-4]
				newFileName = fmt.Sprintf(
					"%s.%s_%s_%s.srt", newFileName,
					customT.fromLang, customT.toLang, carbon.Now().Layout(carbon.ShortDateLayout),
				)
				if err = os.WriteFile(newFileName, bytes, os.ModePerm); err != nil {
					tt_log.GetInstance().Error("写入文件失败", zap.String("源文件", currentSrt.FilePath), zap.String("目标文件", newFileName), zap.Error(err))
					continue
				}
				localChanRes <- fmt.Sprintf("翻译文件[%s]完成", currentSrt.FileName)
				tt_log.GetInstance().Info(fmt.Sprintf("编码srt文件完成, 文件: %s, 目标文件: %s, 文件大小(byte): %d, 耗时(s): %d", currentSrt.FilePath, newFileName, len(bytes), carbon.Now().DiffAbsInSeconds(timeStart)))
				cntEncode++
			}
		}
	}(customT.ctx, chanWG, customT.chanEncode, customT.chanRes)
}

func (customT *Translate) jobRes(chanWG *sync.WaitGroup) {
	go func(localCtx context.Context, localChanWG *sync.WaitGroup, localChanRes chan string) {
		defer func() {
			localChanWG.Done()
		}()
		for true {
			select {
			case currentRes, isOpen := <-localChanRes:
				if currentRes == "" {
					if isOpen {
						close(localChanRes)
					}
					return
				}
				customT.resultMsg = append(customT.resultMsg, currentRes)
			}
		}
	}(customT.ctx, chanWG, customT.chanRes)
}

func (customT *Translate) defaultCtxAndChan() {
	customT.ctx = context.Background()

	customT.chanDecode = make(chan string, 10)
	customT.chanTranslate = make(chan *tt_srt.Srt, 2)
	customT.chanEncode = make(chan *tt_srt.Srt, 10)
	customT.chanRes = make(chan string, 20)
}

func (customT *Translate) validate() error {
	// #todo 后期完善, 这儿暂时不验证相关数据有效性
	return nil
}

func (customT *Translate) Translator() tt_translator.ITranslator {
	return customT.translator
}

func (customT *Translate) SetTranslator(translator tt_translator.ITranslator) *Translate {
	customT.translator = translator
	return customT
}

func (customT *Translate) FromLang() string {
	return customT.fromLang
}

func (customT *Translate) SetFromLang(fromLang string) *Translate {
	customT.fromLang = fromLang
	return customT
}

func (customT *Translate) ToLang() string {
	return customT.toLang
}

func (customT *Translate) SetToLang(toLang string) *Translate {
	customT.toLang = toLang
	return customT
}

func (customT *Translate) TranslateMode() _type.TranslateMode {
	return customT.translateMode
}

func (customT *Translate) SetTranslateMode(translateMode _type.TranslateMode) *Translate {
	customT.translateMode = translateMode
	return customT
}

func (customT *Translate) MainTrackReport() _type.LangDirection {
	return customT.mainTrackReport
}

func (customT *Translate) SetMainTrackReport(mainTrackReport _type.LangDirection) *Translate {
	customT.mainTrackReport = mainTrackReport
	return customT
}

func (customT *Translate) SrtFile() string {
	return customT.srtFile
}

func (customT *Translate) SetSrtFile(srtFile string) *Translate {
	customT.srtFile = srtFile
	return customT
}

func (customT *Translate) SrtDir() string {
	return customT.srtDir
}

func (customT *Translate) SetSrtDir(srtDir string) *Translate {
	customT.srtDir = srtDir
	return customT
}
