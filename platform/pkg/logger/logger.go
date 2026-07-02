package logger

import (
	"context"
	"os"
	"strings"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Level string
	JSON  bool
}

type Logger struct {
	inner *zap.Logger
}

func New(cfg Config) (*Logger, error) {
	level := zap.NewAtomicLevelAt(parseLevel(cfg.Level))
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var encoder zapcore.Encoder
	if cfg.JSON {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)
	return &Logger{inner: zap.New(core, zap.AddCaller())}, nil
}

func NewNop() *Logger {
	return &Logger{inner: zap.NewNop()}
}

func (l *Logger) Debug(ctx context.Context, message string, fields ...zap.Field) {
	l.inner.Debug(message, append(contextFields(ctx), fields...)...)
}

func (l *Logger) Info(ctx context.Context, message string, fields ...zap.Field) {
	l.inner.Info(message, append(contextFields(ctx), fields...)...)
}

func (l *Logger) Warn(ctx context.Context, message string, fields ...zap.Field) {
	l.inner.Warn(message, append(contextFields(ctx), fields...)...)
}

func (l *Logger) Error(ctx context.Context, message string, fields ...zap.Field) {
	l.inner.Error(message, append(contextFields(ctx), fields...)...)
}

func (l *Logger) Sync() error {
	return l.inner.Sync()
}

func parseLevel(value string) zapcore.Level {
	switch strings.ToLower(value) {
	case "debug":
		return zapcore.DebugLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func contextFields(ctx context.Context) []zap.Field {
	spanContext := trace.SpanContextFromContext(ctx)
	if !spanContext.IsValid() {
		return nil
	}

	return []zap.Field{
		zap.String("trace_id", spanContext.TraceID().String()),
		zap.String("span_id", spanContext.SpanID().String()),
	}
}
