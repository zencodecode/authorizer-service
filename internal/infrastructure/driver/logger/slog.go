package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/zencodecode/authorizer-service/internal/domain/service"
)

type Logger struct {
	l *slog.Logger
}

func New(level slog.Level) service.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	})
	return &Logger{l: slog.New(handler)}
}

func (s *Logger) Debug(ctx context.Context, msg string, args ...any) {
	s.l.DebugContext(ctx, msg, args...)
}

func (s *Logger) Info(ctx context.Context, msg string, args ...any) {
	s.l.InfoContext(ctx, msg, args...)
}

func (s *Logger) Warn(ctx context.Context, msg string, args ...any) {
	s.l.WarnContext(ctx, msg, args...)
}

func (s *Logger) Error(ctx context.Context, msg string, args ...any) {
	s.l.ErrorContext(ctx, msg, args...)
}

func (s *Logger) With(args ...any) service.Logger {
	return &Logger{l: s.l.With(args...)}
}
