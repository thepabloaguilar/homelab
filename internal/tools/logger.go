package tools

import (
	"context"
	"errors"
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logContextKey struct{}

func NewLogger() *zap.Logger {
	cfg := zap.NewProductionConfig()

	cfg.Encoding = "console"
	cfg.DisableStacktrace = true
	cfg.DisableCaller = true

	cfg.EncoderConfig.TimeKey = "time"
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	return zap.Must(cfg.Build())
}

func LoggerToContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, logContextKey{}, logger)
}

func LoggerFromContext(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(logContextKey{}).(*zap.Logger)
	if !ok {
		return zap.NewNop()
	}

	return logger
}

func LogCloser(ctx context.Context, closer io.Closer) {
	if err := closer.Close(); err != nil {
		if !errors.Is(err, io.EOF) {
			LoggerFromContext(ctx).Error("closer close error", zap.Error(err))
		}
	}
}
