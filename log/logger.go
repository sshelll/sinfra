package log

import (
	"context"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Entry
}

func NewLogger(ctx context.Context) *Logger {
	logger := &Logger{
		Entry: logrus.NewEntry(logrus.New()).WithContext(ctx),
	}
	return logger
}

func (logger *Logger) AddHook(hook logrus.Hook) {
	logger.Logger.AddHook(hook)
}

func (logger *Logger) WithLogIDHook() *Logger {
	logger.AddHook(NewContextLogIDHook())
	return logger
}

func (logger *Logger) WithContext(ctx context.Context) *Logger {
	return &Logger{
		Entry: logger.Entry.WithContext(ctx),
	}
}

func (logger *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{
		Entry: logger.Entry.WithField(key, value),
	}
}

func (logger *Logger) WithFields(fields logrus.Fields) *Logger {
	return &Logger{
		Entry: logger.Entry.WithFields(fields),
	}
}

func (logger *Logger) Info(format string, args ...interface{}) {
	logger.Entry.Infof(format, args...)
}

func (logger *Logger) Warn(format string, args ...interface{}) {
	logger.Entry.Warnf(format, args...)
}

func (logger *Logger) Error(format string, args ...interface{}) {
	logger.Entry.Errorf(format, args...)
}

func (logger *Logger) Fatal(format string, args ...interface{}) {
	logger.Entry.Fatalf(format, args...)
}

func (logger *Logger) Debug(format string, args ...interface{}) {
	logger.Entry.Debugf(format, args...)
}

func (logger *Logger) Trace(format string, args ...interface{}) {
	logger.Entry.Tracef(format, args...)
}

func (logger *Logger) Panic(format string, args ...interface{}) {
	logger.Entry.Panicf(format, args...)
}

func (logger *Logger) LogID() (string, bool) {
	// get from data
	if v, ok := logger.Data[CtxLogID]; ok {
		if logID, ok := v.(string); ok {
			return logID, true
		}
	}
	// get from ctx
	if logger.Entry.Context == nil {
		return "", false
	}
	if v := logger.Entry.Context.Value(CtxLogID); v != nil {
		if logID, ok := v.(string); ok {
			return logID, true
		}
	}
	return "", false
}
