package port

import "context"

type Logger interface {
	InfoCtx(ctx context.Context, msg string, fields ...any)
	ErrorCtx(ctx context.Context, msg string, fields ...any)
	WarnCtx(ctx context.Context, msg string, fields ...any)
	DebugCtx(ctx context.Context, msg string, fields ...any)
}
