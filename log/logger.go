package log

type Logger interface {
	// Debug Using：log.Debug("test")
	Debug(args ...interface{})

	// Debugf Using：log.Debugf("test:%s", err)
	Debugf(template string, args ...interface{})

	// Debugw Using：log.Debugw("test", "field1", "value1", "field2", "value2")
	Debugw(msg string, keysAndValues ...interface{})

	Info(args ...interface{})

	Infof(template string, args ...interface{})

	Infow(msg string, keysAndValues ...interface{})

	Warn(args ...interface{})

	Warnf(template string, args ...interface{})

	Warnw(msg string, keysAndValues ...interface{})

	Error(args ...interface{})

	Errorf(template string, args ...interface{})

	Errorw(msg string, keysAndValues ...interface{})

	Sync()
}
