package log

import (
	"context"
	"strings"

	"github.com/google/uuid"
)

func GenCtxWithLogID() context.Context {
	logID := GenLogID()
	return context.WithValue(context.Background(), CtxLogID, logID)
}

func GenLogIDForCtx(ctx context.Context) context.Context {
	v := ctx.Value(CtxLogID)
	if v != nil {
		return ctx
	}
	return context.WithValue(ctx, CtxLogID, GenLogID())
}

func GenLogID() string {
	logID := uuid.New().String()
	logID = strings.Replace(logID, "-", "", -1)
	return logID
}
