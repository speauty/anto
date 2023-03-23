package tt_log

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
	apiTTLog  *TTLog
	onceTTLog sync.Once
)

func GetInstance() *TTLog {
	onceTTLog.Do(func() {
		apiTTLog = new(TTLog)
		apiTTLog.genEncoder()
		apiTTLog.genLogger()
	})
	return apiTTLog
}

type TTLog struct { // 接入部分API
	logger  *zap.SugaredLogger
	encoder zapcore.Encoder
}

func (customTL *TTLog) Debug(args ...interface{}) {
	customTL.logger.Debug(args...)
}
func (customTL *TTLog) Info(args ...interface{}) {
	customTL.logger.Info(args...)
}
func (customTL *TTLog) Warn(args ...interface{}) {
	customTL.logger.Warn(args...)
}
func (customTL *TTLog) Error(args ...interface{}) {
	customTL.logger.Error(args...)
}
func (customTL *TTLog) Panic(args ...interface{}) {
	customTL.logger.Panic(args...)
}
func (customTL *TTLog) Fatal(args ...interface{}) {
	customTL.logger.Fatal(args...)
}

func (customTL *TTLog) genLogger() {
	infoLevel := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv >= zapcore.InfoLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv >= zapcore.ErrorLevel
	})

	ioInfo := customTL.genWriter("./logs/info.log")
	ioErr := customTL.genWriter("./logs/error.log")

	core := zapcore.NewTee(
		zapcore.NewCore(customTL.encoder, zapcore.AddSync(ioInfo), infoLevel),
		zapcore.NewCore(customTL.encoder, zapcore.AddSync(ioErr), errorLevel),
	)
	customTL.logger = zap.New(core, zap.AddCaller()).Sugar()
}

func (customTL *TTLog) genWriter(filename string) io.Writer {
	fd, err := rotatelogs.New(strings.Replace(filename, ".log", "", -1) + ".%Y%m%d.log")
	if err != nil {
		panic(err)
	}
	return fd
}

func (customTL *TTLog) genEncoder() {
	customTL.encoder = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey: "msg", LevelKey: "lv", TimeKey: "ts", StacktraceKey: "trace",
		EncodeLevel: zapcore.LowercaseLevelEncoder,
		EncodeTime: func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(carbon.FromStdTime(time).Layout(carbon.ShortDateTimeLayout))
		},
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})
}
