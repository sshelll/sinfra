package log

import "github.com/sirupsen/logrus"

const (
	CtxLogID = "sinfra-log_id"
)

type ContextLogIDHook struct{}

func NewContextLogIDHook() *ContextLogIDHook {
	return &ContextLogIDHook{}
}

func (hook *ContextLogIDHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *ContextLogIDHook) Fire(entry *logrus.Entry) error {
	if ctx := entry.Context; ctx != nil {
		if logID := ctx.Value(CtxLogID); logID != nil {
			entry.Data[CtxLogID] = logID
		}
	}
	return nil
}
