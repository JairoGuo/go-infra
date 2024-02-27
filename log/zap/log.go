package zap

import (
	"github.com/jairoguo/go-infra/log"
	"github.com/jairoguo/go-infra/log/zap/config"
	"github.com/jairoguo/go-infra/log/zap/core"
	utilos "github.com/jairoguo/go-infra/util/path"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type LogZap struct {
	log *zap.Logger
}

func NewLogZap(l *zap.Logger) log.Logger {
	return &LogZap{
		log: l,
	}
}

func NewLogZapByConfig(c config.Config) log.Logger {
	return &LogZap{
		log: buildZapLoggerByConfig(c),
	}
}

func buildZapLoggerByConfig(c config.Config) (logger *zap.Logger) {

	core.BindConfig(c)

	if ok, _ := utilos.PathExists(c.Director); !ok { // 判断是否有Director文件夹
		_ = os.Mkdir(c.Director, os.ModePerm)
	}

	cores := core.ZapOption.GetZapCores(c)
	logger = zap.New(zapcore.NewTee(cores...))

	if c.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	return logger
}

func (l *LogZap) Info(msg string, fields []log.Field) {
	l.log.Info(msg, getFields(fields)...)
}

func (l *LogZap) Debug(msg string, fields []log.Field) {
	l.log.Debug(msg, getFields(fields)...)

}

func (l *LogZap) Warn(msg string, fields []log.Field) {
	l.log.Warn(msg, getFields(fields)...)

}

func (l *LogZap) Error(msg string, fields []log.Field) {
	l.log.Error(msg, getFields(fields)...)

}

func (l *LogZap) Panic(msg string, fields []log.Field) {
	l.log.Panic(msg, getFields(fields)...)

}

func (l *LogZap) Fatal(msg string, fields []log.Field) {
	l.log.Fatal(msg, getFields(fields)...)

}

func getFields(fields []log.Field) []zap.Field {
	var zapFields []zap.Field

	for _, field := range fields {
		if field.Key == "" {
			zapFields = append(zapFields, zap.Any("!BADKEY", field.Value))
		} else {
			zapFields = append(zapFields, zap.Any(field.Key, field.Value))
		}
	}

	return zapFields
}
