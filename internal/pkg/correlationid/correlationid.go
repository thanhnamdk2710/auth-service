package correlationid

import (
	"context"

	"github.com/google/uuid"
)

type ctxKey struct{}

const HeaderName = "X-Correlation-ID"

func New() string {
	return uuid.New().String()
}

func FromContext(ctx context.Context) string {
	if id, ok := ctx.Value(ctxKey{}).(string); ok {
		return id
	}
	return ""
}

func WithContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxKey{}, id)
}
