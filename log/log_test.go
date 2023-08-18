package log

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLog(t *testing.T) {
	ctx := GenCtxWithLogID()
	Info(ctx, "info, %d", 1)
	Warn(ctx, "warn, %d", 2)
	Error(ctx, "error, %d", 3)
	Debug(ctx, "debug, %d", 4)
	Trace(ctx, "trace, %d", 5)
	fields := logrus.Fields{
		"key1": "value1",
		"key2": "value2",
	}
	fields2 := logrus.Fields{
		"key3": "value3",
		"key4": "value4",
	}
	WithFields(ctx, fields).Info("info, %d", 1)
	WithFields(ctx, fields2).Warn("warn, %d", 2)
	WithFields(ctx, fields).Error("error, %d", 3)
	WithFields(ctx, fields2).Debug("debug, %d", 4)
	WithFields(ctx, fields).Trace("trace, %d", 5)
}

func TestLogger(t *testing.T) {
	ctx := GenCtxWithLogID()
	logger := NewLogger(ctx).WithLogIDHook()
	logger.Info("info, %d", 1)
	logger.Warn("warn, %d", 2)
	logger.Error("error, %d", 3)
	logger.Debug("debug, %d", 4)
	logger.Trace("trace, %d", 5)
	fields := logrus.Fields{
		"key1": "value1",
		"key2": "value2",
	}
	fields2 := logrus.Fields{
		"key3": "value3",
		"key4": "value4",
	}
	logger.WithFields(fields).Info("info, %d", 1)
	logger.WithFields(fields2).Warn("warn, %d", 2)
	logger.WithFields(fields).Error("error, %d", 3)
	logger.WithFields(fields2).Debug("debug, %d", 4)
	logger.WithFields(fields).Trace("trace, %d", 5)
}
