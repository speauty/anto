package log

import (
	"github.com/golang-module/carbon"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"strings"
	"sync"
	"time"
)

var (
	apiLog  *Log
	onceLog sync.Once
)

func Singleton() *Log {
	onceLog.Do(func() {
		apiLog = new(Log)
		apiLog.initEncoder()
		apiLog.initLogger()
	})
	return apiLog
}

type Log struct {
	logger  *zap.SugaredLogger
	encoder zapcore.Encoder
}

func (customLog *Log) Debug(args ...interface{}) {
	customLog.logger.Debug(args...)
}
func (customLog *Log) Info(args ...interface{}) {
	customLog.logger.Info(args...)
}
func (customLog *Log) InfoF(tpl string, args ...interface{}) {
	customLog.logger.Infof(tpl, args...)
}
func (customLog *Log) Warn(args ...interface{}) {
	customLog.logger.Warn(args...)
}
func (customLog *Log) WarnF(tpl string, args ...interface{}) {
	customLog.logger.Warnf(tpl, args...)
}
func (customLog *Log) Error(args ...interface{}) {
	customLog.logger.Error(args...)
}
func (customLog *Log) ErrorF(tpl string, args ...interface{}) {
	customLog.logger.Errorf(tpl, args...)
}
func (customLog *Log) Panic(args ...interface{}) {
	customLog.logger.Panic(args...)
}
func (customLog *Log) Fatal(args ...interface{}) {
	customLog.logger.Fatal(args...)
}

func (customLog *Log) initLogger() {
	infoLevel := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv >= zapcore.InfoLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv >= zapcore.ErrorLevel
	})

	ioInfo := customLog.initWriter("./logs/info.log")
	ioErr := customLog.initWriter("./logs/error.log")

	core := zapcore.NewTee(
		zapcore.NewCore(customLog.encoder, zapcore.AddSync(ioInfo), infoLevel),
		zapcore.NewCore(customLog.encoder, zapcore.AddSync(ioErr), errorLevel),
	)
	customLog.logger = zap.New(core, zap.AddCaller()).Sugar()
}

func (customLog *Log) initWriter(filename string) io.Writer {
	fd, err := rotatelogs.New(strings.Replace(filename, ".log", "", -1) + ".%Y%m%d.log")
	if err != nil {
		panic(err)
	}
	return fd
}

func (customLog *Log) initEncoder() {
	customLog.encoder = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey: "msg", LevelKey: "level", TimeKey: "ts", StacktraceKey: "trace",
		EncodeLevel: zapcore.LowercaseLevelEncoder,
		EncodeTime: func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(carbon.FromStdTime(time).Layout(carbon.ShortDateTimeLayout))
		},
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})
}
