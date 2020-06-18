package backtest

type EmptyLogger struct {
}

func (l *EmptyLogger) Debug(args ...interface{}) {

}

func (l *EmptyLogger) Debugf(template string, args ...interface{}) {

}

func (l *EmptyLogger) Debugw(msg string, keysAndValues ...interface{}) {

}

func (l *EmptyLogger) Info(args ...interface{}) {

}

func (l *EmptyLogger) Infof(template string, args ...interface{}) {

}

func (l *EmptyLogger) Infow(msg string, keysAndValues ...interface{}) {

}

func (l *EmptyLogger) Warn(args ...interface{}) {

}

func (l *EmptyLogger) Warnf(template string, args ...interface{}) {

}

func (l *EmptyLogger) Warnw(msg string, keysAndValues ...interface{}) {

}

func (l *EmptyLogger) Error(args ...interface{}) {

}

func (l *EmptyLogger) Errorf(template string, args ...interface{}) {

}

func (l *EmptyLogger) Errorw(msg string, keysAndValues ...interface{}) {

}

func (l *EmptyLogger) Sync() {

}
