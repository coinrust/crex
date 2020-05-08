package log

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	PanicLevel = "panic"
)

var logger Logger

func SetLogger(myLogger Logger) {
	logger = myLogger
}

// Debug Using：log.Debug("test")
func Debug(args ...interface{}) {
	if logger == nil {
		return
	}

	logger.Debug(args...)
}

// Debugf Using：log.Debugf("test:%s", err)
func Debugf(template string, args ...interface{}) {
	if logger == nil {
		return
	}

	logger.Debugf(template, args...)
}

// Debugw Using：log.Debugw("test", "field1", "value1", "field2", "value2")
func Debugw(msg string, keysAndValues ...interface{}) {
	if logger == nil {
		return
	}

	logger.Debugw(msg, keysAndValues...)
}

func Info(args ...interface{}) {
	if logger == nil {
		return
	}

	logger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	if logger == nil {
		return
	}

	logger.Infof(template, args...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	if logger == nil {
		return
	}

	logger.Infow(msg, keysAndValues...)
}

func Warn(args ...interface{}) {
	if logger == nil {
		return
	}

	logger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	if logger == nil {
		return
	}

	logger.Warnf(template, args...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	if logger == nil {
		return
	}

	logger.Warnw(msg, keysAndValues...)
}

func Error(args ...interface{}) {
	if logger == nil {
		return
	}

	logger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	if logger == nil {
		return
	}

	logger.Errorf(template, args...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	if logger == nil {
		return
	}

	logger.Errorw(msg, keysAndValues...)
}

func Sync() {
	if logger == nil {
		return
	}
	logger.Sync()
}
