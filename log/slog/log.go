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

func (l *LogSlog) Info(msg string, fields []log.Field) {
	slog.Info(msg, getFields(fields)...)
}

func (l *LogSlog) Debug(msg string, fields []log.Field) {
	slog.Debug(msg, getFields(fields)...)
}

func (l *LogSlog) Warn(msg string, fields []log.Field) {
	slog.Warn(msg, getFields(fields)...)
}

func (l *LogSlog) Error(msg string, fields []log.Field) {
	slog.Error(msg, getFields(fields)...)
}

func (l *LogSlog) Panic(msg string, fields []log.Field) {

}

func (l *LogSlog) Fatal(msg string, fields []log.Field) {

}

func getFields(fields []log.Field) []any {
	var argsFields []any

	for _, field := range fields {
		if field.Key != "" {
			argsFields = append(argsFields, field.Key)
			argsFields = append(argsFields, field.Value)
		} else {
			argsFields = append(argsFields, field.Value)
		}
	}

	return argsFields
}
