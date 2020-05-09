package backtest

import (
	"go.uber.org/zap/zapcore"
	"time"
)

const (
	LogTsKey = "log_ts"
)

// NewCore creates a Core that writes logs to a WriteSyncer.
func NewCore(enc zapcore.Encoder, ws zapcore.WriteSyncer, enab zapcore.LevelEnabler) zapcore.Core {
	return &ioCore{
		LevelEnabler: enab,
		enc:          enc,
		out:          ws,
	}
}

type ioCore struct {
	zapcore.LevelEnabler
	enc zapcore.Encoder
	out zapcore.WriteSyncer
}

func (c *ioCore) With(fields []zapcore.Field) zapcore.Core {
	clone := c.clone()
	addFields(clone.enc, fields)
	return clone
}

func (c *ioCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c *ioCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	var lIndex = -1

	for i := 0; i < len(fields); i++ {
		v := &fields[i]
		if v.Key == LogTsKey {
			//integer := ent.Time.UnixNano()
			ent.Time = time.Unix(0, v.Integer)
			//v.Integer = integer
			lIndex = i
			break
		}
	}

	if lIndex != -1 {
		fieldsRemove2(&fields, lIndex)
	}

	buf, err := c.enc.EncodeEntry(ent, fields)
	if err != nil {
		return err
	}
	_, err = c.out.Write(buf.Bytes())
	buf.Free()
	if err != nil {
		return err
	}
	if ent.Level > zapcore.ErrorLevel {
		// Since we may be crashing the program, sync the output. Ignore Sync
		// errors, pending a clean solution to issue #370.
		c.Sync()
	}
	return nil
}

func (c *ioCore) Sync() error {
	return c.out.Sync()
}

func (c *ioCore) clone() *ioCore {
	return &ioCore{
		LevelEnabler: c.LevelEnabler,
		enc:          c.enc.Clone(),
		out:          c.out,
	}
}

func addFields(enc zapcore.ObjectEncoder, fields []zapcore.Field) {
	for i := range fields {
		fields[i].AddTo(enc)
	}
}

func fieldsRemove(s []zapcore.Field, index int) []zapcore.Field {
	return append(s[:index], s[index+1:]...)
}

func fieldsRemove2(s *[]zapcore.Field, index int) {
	*s = append((*s)[:index], (*s)[index+1:]...)
}
