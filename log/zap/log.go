package zap

import (
	"fmt"
	"github.com/jairoguo/go-infra/log"
	"github.com/jairoguo/go-infra/log/zap/core"
	utilos "github.com/jairoguo/go-infra/util/os"
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

func NewLogZapConfig(c Config) log.Logger {
	return &LogZap{
		log: NewZapLoggerByConfig(c),
	}
}

func NewZapLoggerByConfig(config Config) (logger *zap.Logger) {

	if ok, _ := utilos.PathExists(config.Director); !ok { // 判断是否有Director文件夹
		fmt.Printf("create %v directory\n", config.Director)
		_ = os.Mkdir(config.Director, os.ModePerm)
	}

	cores := core.ZapOption.GetZapCores(config)
	logger = zap.New(zapcore.NewTee(cores...))

	if config.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	return logger
}

func (l *LogZap) Info(msg string, args ...any) {
	l.log.Info(msg)
}

func (l *LogZap) Debug(msg string, args ...any) {
	l.log.Debug(msg)

}

func (l *LogZap) Warn(msg string, args ...any) {
	l.log.Warn(msg)

}

func (l *LogZap) Error(msg string, args ...any) {
	l.log.Error(msg)

}

func (l *LogZap) Panic(msg string, args ...any) {
	l.log.Panic(msg)

}

func (l *LogZap) Fatal(msg string, args ...any) {
	l.log.Fatal(msg)

}
