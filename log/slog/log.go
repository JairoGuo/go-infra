package slog

import (
	"github.com/jairoguo/go-infra/log"
	"log/slog"
)

type LogSlog struct {
}

func NewLogSlog(handler slog.Handler) log.Logger {
	slog.SetDefault(slog.New(handler))
	return &LogSlog{}
}

func (l *LogSlog) Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

func (l *LogSlog) Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

func (l *LogSlog) Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

func (l *LogSlog) Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

func (l *LogSlog) Panic(msg string, args ...any) {

}

func (l *LogSlog) Fatal(msg string, args ...any) {

}
