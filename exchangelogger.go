package crex

// EmptyExchangeLogger 交易所撮合日志
type EmptyExchangeLogger struct {
}

func (l *EmptyExchangeLogger) Debug(args ...interface{}) {

}

func (l *EmptyExchangeLogger) Debugf(template string, args ...interface{}) {

}

func (l *EmptyExchangeLogger) Debugw(msg string, keysAndValues ...interface{}) {

}

func (l *EmptyExchangeLogger) Info(args ...interface{}) {

}

func (l *EmptyExchangeLogger) Infof(template string, args ...interface{}) {

}

func (l *EmptyExchangeLogger) Infow(msg string, keysAndValues ...interface{}) {

}

func (l *EmptyExchangeLogger) Warn(args ...interface{}) {

}

func (l *EmptyExchangeLogger) Warnf(template string, args ...interface{}) {

}

func (l *EmptyExchangeLogger) Warnw(msg string, keysAndValues ...interface{}) {

}

func (l *EmptyExchangeLogger) Error(args ...interface{}) {

}

func (l *EmptyExchangeLogger) Errorf(template string, args ...interface{}) {

}

func (l *EmptyExchangeLogger) Errorw(msg string, keysAndValues ...interface{}) {

}

func (l *EmptyExchangeLogger) Sync() {

}
