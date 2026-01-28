package logger

import (
	"context"

	"go.uber.org/zap"

	pkglogger "github.com/thanhnamdk2710/auth-service/internal/pkg/logger"
)

type Adapter struct {
	logger *pkglogger.Logger
}

func NewAdapter(logger *pkglogger.Logger) *Adapter {
	return &Adapter{logger: logger}
}

func (a *Adapter) InfoCtx(ctx context.Context, msg string, keysAndValues ...any) {
	a.logger.InfoCtx(ctx, msg, toZapFields(keysAndValues)...)
}

func (a *Adapter) ErrorCtx(ctx context.Context, msg string, keysAndValues ...any) {
	a.logger.ErrorCtx(ctx, msg, toZapFields(keysAndValues)...)
}

func (a *Adapter) WarnCtx(ctx context.Context, msg string, keysAndValues ...any) {
	a.logger.WarnCtx(ctx, msg, toZapFields(keysAndValues)...)
}

func (a *Adapter) DebugCtx(ctx context.Context, msg string, keysAndValues ...any) {
	a.logger.DebugCtx(ctx, msg, toZapFields(keysAndValues)...)
}

func toZapFields(keysAndValues []any) []zap.Field {
	fields := make([]zap.Field, 0, len(keysAndValues)/2)

	for i := 0; i < len(keysAndValues)-1; i += 2 {
		key, ok := keysAndValues[i].(string)
		if !ok {
			continue
		}
		fields = append(fields, zap.Any(key, keysAndValues[i+1]))
	}

	return fields
}
