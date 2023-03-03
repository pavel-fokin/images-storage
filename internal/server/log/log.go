package log

import (
	"context"

	"github.com/go-chi/httplog"
)

func Error(ctx context.Context, err error, msg string) {
	oplog := httplog.LogEntry(ctx)
	oplog.Error().Err(err).Msg(msg)
}

func Info(ctx context.Context, msg string) {
	oplog := httplog.LogEntry(ctx)
	oplog.Info().Msg(msg)
}
