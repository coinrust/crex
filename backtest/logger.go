package backtest

import (
	"fmt"
	"github.com/coinrust/crex/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

type GetCurrentTime interface {
	GetTime() time.Time
}

type BtLogger struct {
	Path        string // 文件路径，如：./app.log
	Level       string // 日志输出的级别
	MaxFileSize int    // 日志文件大小的最大值，单位(M)
	MaxBackups  int    // 最多保留备份数
	MaxAge      int    // 日志文件保存的时间，单位(天)
	Compress    bool   // 是否压缩
	Caller      bool   // 日志是否需要显示调用位置
	JsonFormat  bool   // 是否以Json格式输出
	Stdout      bool   // 是否输出到控制台

	gct    GetCurrentTime
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

func NewBtLogger(gct GetCurrentTime, path string, level string, jsonFormat bool, stdout bool) *BtLogger {
	logger := &BtLogger{
		Path:        path,
		Level:       level,
		MaxFileSize: 512,
		MaxBackups:  10,
		MaxAge:      60,
		Compress:    false,
		Caller:      true,
		JsonFormat:  jsonFormat,
		Stdout:      stdout,
		gct:         gct,
		logger:      nil,
		sugar:       nil,
	}
	logger.build()
	return logger
}

func (l *BtLogger) build() {
	writeSyncer := []zapcore.WriteSyncer{
		zapcore.AddSync(l.createLumberjackHook()),
	}

	if l.Stdout {
		writeSyncer = append(writeSyncer, zapcore.Lock(os.Stdout))
	}

	var level zapcore.Level
	switch l.Level {
	case log.DebugLevel:
		level = zap.DebugLevel
	case log.InfoLevel:
		level = zap.InfoLevel
	case log.WarnLevel:
		level = zap.WarnLevel
	case log.ErrorLevel:
		level = zap.ErrorLevel
	case log.PanicLevel:
		level = zap.PanicLevel
	default:
		level = zap.InfoLevel
	}

	conf := zap.NewProductionEncoderConfig() // "2006-01-02 15:04:05.000"
	conf.EncodeTime = zapcore.ISO8601TimeEncoder
	var cnf zapcore.Encoder
	if l.JsonFormat {
		cnf = zapcore.NewJSONEncoder(conf)
	} else {
		cnf = zapcore.NewConsoleEncoder(conf)
	}
	core := NewCore(cnf,
		zapcore.NewMultiWriteSyncer(writeSyncer...),
		level)

	l.logger = zap.New(core)
	if l.Caller {
		l.logger = l.logger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(2))
	}
	l.sugar = l.logger.Sugar()
}

// createLumberjackHook 创建LumberjackHook，其作用是为了将日志文件切割，压缩
func (l *BtLogger) createLumberjackHook() *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   l.Path,
		MaxSize:    l.MaxFileSize,
		MaxBackups: l.MaxBackups,
		MaxAge:     l.MaxAge,
		Compress:   l.Compress,
	}
}

func (l *BtLogger) Debug(args ...interface{}) {
	l.sugar.Debugw(formatMessage("", args...), zap.Time(LogTsKey, l.getCurrentTs()))
}

func (l *BtLogger) Debugf(template string, args ...interface{}) {
	l.sugar.Debugw(formatMessage(template, args...), zap.Time(LogTsKey, l.getCurrentTs()))
}

func (l *BtLogger) Debugw(msg string, keysAndValues ...interface{}) {
	kvs := append([]interface{}{zap.Time(LogTsKey, l.getCurrentTs())}, keysAndValues...)
	l.sugar.Debugw(msg, kvs...)
}

func (l *BtLogger) Info(args ...interface{}) {
	l.sugar.Infow(formatMessage("", args...), zap.Time(LogTsKey, l.getCurrentTs()))
}

func (l *BtLogger) Infof(template string, args ...interface{}) {
	l.sugar.Infow(formatMessage(template, args...), zap.Time(LogTsKey, l.getCurrentTs()))
}

func (l *BtLogger) Infow(msg string, keysAndValues ...interface{}) {
	kvs := append([]interface{}{zap.Time(LogTsKey, l.getCurrentTs())}, keysAndValues...)
	l.sugar.Infow(msg, kvs...)
}

func (l *BtLogger) Warn(args ...interface{}) {
	l.sugar.Warnw(formatMessage("", args...), zap.Time(LogTsKey, l.getCurrentTs()))
}

func (l *BtLogger) Warnf(template string, args ...interface{}) {
	l.sugar.Warnw(formatMessage(template, args...), zap.Time(LogTsKey, l.getCurrentTs()))
}

func (l *BtLogger) Warnw(msg string, keysAndValues ...interface{}) {
	kvs := append([]interface{}{zap.Time(LogTsKey, l.getCurrentTs())}, keysAndValues...)
	l.sugar.Warnw(msg, kvs...)
}

func (l *BtLogger) Error(args ...interface{}) {
	l.sugar.Errorw(formatMessage("", args...), zap.Time(LogTsKey, l.getCurrentTs()))
}

func (l *BtLogger) Errorf(template string, args ...interface{}) {
	l.sugar.Errorw(formatMessage(template, args...), zap.Time(LogTsKey, l.getCurrentTs()))
}

func (l *BtLogger) Errorw(msg string, keysAndValues ...interface{}) {
	kvs := append([]interface{}{zap.Time(LogTsKey, l.getCurrentTs())}, keysAndValues...)
	l.sugar.Errorw(msg, kvs...)
}

func formatMessage(template string, fmtArgs ...interface{}) string {
	// Format with Sprint, Sprintf, or neither.
	msg := template
	if msg == "" && len(fmtArgs) > 0 {
		msg = fmt.Sprint(fmtArgs...)
	} else if msg != "" && len(fmtArgs) > 0 {
		msg = fmt.Sprintf(template, fmtArgs...)
	}
	return msg
}

func (l *BtLogger) getCurrentTs() (ts time.Time) {
	if l.gct == nil {
		return time.Now()
	}
	return l.gct.GetTime()
}

func (l *BtLogger) Sync() {
	l.sugar.Sync()
}
