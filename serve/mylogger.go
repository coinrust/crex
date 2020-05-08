package serve

import (
	"github.com/coinrust/crex/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

type MyLogger struct {
	Path        string // 文件路径，如：./app.log
	Level       string // 日志输出的级别
	MaxFileSize int    // 日志文件大小的最大值，单位(M)
	MaxBackups  int    // 最多保留备份数
	MaxAge      int    // 日志文件保存的时间，单位(天)
	Compress    bool   // 是否压缩
	Caller      bool   // 日志是否需要显示调用位置
	JsonFormat  bool   // 是否以Json格式输出
	Stdout      bool   // 是否输出到控制台

	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

func NewMyLogger(path string, level string, jsonFormat bool) *MyLogger {
	logger := &MyLogger{
		Path:        path,
		Level:       level,
		MaxFileSize: 512,
		MaxBackups:  10,
		MaxAge:      60,
		Compress:    false,
		Caller:      true,
		JsonFormat:  jsonFormat,
		Stdout:      true,
		logger:      nil,
		sugar:       nil,
	}
	logger.build()
	return logger
}

func (l *MyLogger) build() {
	writeSyncer := []zapcore.WriteSyncer{
		zapcore.AddSync(l.createLumberjackHook()),
	}

	//if l.Stdout {
	writeSyncer = append(writeSyncer, zapcore.Lock(os.Stdout))
	//}

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
	core := zapcore.NewCore(cnf,
		zapcore.NewMultiWriteSyncer(writeSyncer...),
		level)

	l.logger = zap.New(core)
	if l.Caller {
		l.logger = l.logger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(2))
	}
	l.sugar = l.logger.Sugar()
}

// createLumberjackHook 创建LumberjackHook，其作用是为了将日志文件切割，压缩
func (l *MyLogger) createLumberjackHook() *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   l.Path,
		MaxSize:    l.MaxFileSize,
		MaxBackups: l.MaxBackups,
		MaxAge:     l.MaxAge,
		Compress:   l.Compress,
	}
}

func (l *MyLogger) Debug(args ...interface{}) {
	l.sugar.Debug(args...)
}

func (l *MyLogger) Debugf(template string, args ...interface{}) {
	l.sugar.Debugf(template, args...)
}

func (l *MyLogger) Debugw(msg string, keysAndValues ...interface{}) {
	l.sugar.Debugw(msg, keysAndValues...)
}

func (l *MyLogger) Info(args ...interface{}) {
	l.sugar.Info(args...)
}

func (l *MyLogger) Infof(template string, args ...interface{}) {
	l.sugar.Infof(template, args...)
}

func (l *MyLogger) Infow(msg string, keysAndValues ...interface{}) {
	l.sugar.Infow(msg, keysAndValues...)
}

func (l *MyLogger) Warn(args ...interface{}) {
	l.sugar.Warn(args...)
}

func (l *MyLogger) Warnf(template string, args ...interface{}) {
	l.sugar.Warnf(template, args...)
}

func (l *MyLogger) Warnw(msg string, keysAndValues ...interface{}) {
	l.sugar.Warnw(msg, keysAndValues...)
}

func (l *MyLogger) Error(args ...interface{}) {
	l.sugar.Error(args...)
}

func (l *MyLogger) Errorf(template string, args ...interface{}) {
	l.sugar.Errorw(template, args...)
}

func (l *MyLogger) Errorw(msg string, keysAndValues ...interface{}) {
	l.sugar.Errorw(msg, keysAndValues...)
}

func (l *MyLogger) Sync() {
	log.Sync()
}
