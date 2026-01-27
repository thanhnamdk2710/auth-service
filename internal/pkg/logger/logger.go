package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxKey struct{}

type Logger struct {
	*zap.Logger
}

type Config struct {
	Level       string
	Environment string
}

func New(cfg *Config) (*Logger, error) {
	var zapCfg zap.Config

	if cfg.Environment == "production" {
		zapCfg = zap.NewProductionConfig()
		zapCfg.EncoderConfig.TimeKey = "timestamp"
		zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		zapCfg = zap.NewDevelopmentConfig()
		zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	zapCfg.Level = zap.NewAtomicLevelAt(parseLevel(cfg.Level))

	log, err := zapCfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}

	return &Logger{Logger: log}, nil
}

func parseLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	if ctx == nil {
		return l
	}

	if corrID, ok := ctx.Value(ctxKey{}).(string); ok && corrID != "" {
		return &Logger{Logger: l.With(zap.String("correlation_id", corrID))}
	}

	return l
}

func (l *Logger) InfoCtx(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithContext(ctx).Info(msg, fields...)
}

func (l *Logger) ErrorCtx(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithContext(ctx).Error(msg, fields...)
}

func (l *Logger) WarnCtx(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithContext(ctx).Warn(msg, fields...)
}

func (l *Logger) DebugCtx(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithContext(ctx).Debug(msg, fields...)
}

func (l *Logger) FatalCtx(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithContext(ctx).Fatal(msg, fields...)
}

func WithCorrelationID(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, ctxKey{}, correlationID)
}

func CorrelationIDFromContext(ctx context.Context) string {
	if corrID, ok := ctx.Value(ctxKey{}).(string); ok {
		return corrID
	}
	return ""
}
