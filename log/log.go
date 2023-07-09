package log

import (
	"context"

	"github.com/sirupsen/logrus"
)

var StdLogger *Logger

func WithField(ctx context.Context, key string, value interface{}) *Logger {
	return &Logger{
		Entry: StdLogger.Entry.WithContext(ctx).WithField(key, value),
	}
}

func WithFields(ctx context.Context, fields logrus.Fields) *Logger {
	return &Logger{
		Entry: StdLogger.Entry.WithContext(ctx).WithFields(fields),
	}
}

func Info(ctx context.Context, format string, args ...interface{}) {
	StdLogger.WithContext(ctx).Infof(format, args...)
}

func Warn(ctx context.Context, format string, args ...interface{}) {
	StdLogger.WithContext(ctx).Warnf(format, args...)
}

func Error(ctx context.Context, format string, args ...interface{}) {
	StdLogger.WithContext(ctx).Errorf(format, args...)
}

func Fatal(ctx context.Context, format string, args ...interface{}) {
	StdLogger.WithContext(ctx).Fatalf(format, args...)
}

func Debug(ctx context.Context, format string, args ...interface{}) {
	StdLogger.WithContext(ctx).Debugf(format, args...)
}

func Trace(ctx context.Context, format string, args ...interface{}) {
	StdLogger.WithContext(ctx).Tracef(format, args...)
}

func Panic(ctx context.Context, format string, args ...interface{}) {
	StdLogger.WithContext(ctx).Panicf(format, args...)
}

func init() {
	logger := logrus.New()
	logger.AddHook(NewContextLogIDHook())
	StdLogger = &Logger{
		Entry: logrus.NewEntry(logger),
	}
}
