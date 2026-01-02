package trace

import (
	"context"
	"fmt"
	"log"
)

type Logger struct {
	logger *log.Logger
}

func NewLogger() *Logger {
	return NewLoggerWithLogger(log.Default())
}

func NewLoggerWithLogger(logger *log.Logger) *Logger {
	return &Logger{
		logger: logger,
	}
}

func (l *Logger) Debug(ctx context.Context, message string, fields ...any) {
	l.log(ctx, "DEBUG", message, fields...)
}

func (l *Logger) Info(ctx context.Context, message string, fields ...any) {
	l.log(ctx, "INFO", message, fields...)
}

func (l *Logger) Warn(ctx context.Context, message string, fields ...any) {
	l.log(ctx, "WARN", message, fields...)
}

func (l *Logger) Error(ctx context.Context, message string, fields ...any) {
	l.log(ctx, "ERROR", message, fields...)
}

func (l *Logger) log(ctx context.Context, level, message string, fields ...any) {
	prefix := ""

	span, ok := FromContext(ctx)
	if ok {
		prefix = fmt.Sprintf("[traceID:%s, spanID:%s]", span.TraceID, span.SpanID)
	}

	fullMsg := fmt.Sprintf("%s[%s] %s", prefix, level, message)
	l.logger.Printf(fullMsg, fields...)
}
