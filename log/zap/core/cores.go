package core

import (
	logzap "github.com/jairoguo/go-infra/log/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
	"time"
)

var ZapOption = new(option)

type option struct {
	logzap.Config
}

// GetZapCores 根据配置文件的Level获取 []zapcore.Core
func (opt *option) GetZapCores(c logzap.Config) []zapcore.Core {

	cores := make([]zapcore.Core, 0, 7)
	for level := TransportLevel(opt.Level); level <= zapcore.FatalLevel; level++ {
		cores = append(cores, opt.GetEncoderCore(level, GetLevelPriority(level)))
	}
	return cores
}

// GetEncoder 获取 zapcore.Encoder
func (opt *option) GetEncoder() zapcore.Encoder {
	if opt.Format == "json" {
		return zapcore.NewJSONEncoder(opt.GetEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(opt.GetEncoderConfig())
}

// GetEncoderConfig 获取zapcore.EncoderConfig
func (opt *option) GetEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  opt.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    opt.ZapEncodeLevel(),
		EncodeTime:     opt.CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
}

// GetEncoderCore 获取Encoder的 zapcore.Core
func (opt *option) GetEncoderCore(l zapcore.Level, level zap.LevelEnablerFunc) zapcore.Core {
	writer := opt.GetWriteSyncer(l.String()) // 日志分割
	return zapcore.NewCore(opt.GetEncoder(), writer, level)
}

// CustomTimeEncoder 自定义日志输出时间格式
func (opt *option) CustomTimeEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(opt.Prefix + " " + t.Format("2006-01-02 - 15:04:05.000"))
}

// GetWriteSyncer 获取 zapcore.WriteSyncer
func (opt *option) GetWriteSyncer(level string) zapcore.WriteSyncer {
	fileWriter := NewCutter(opt.Director, level, WithCutterFormat("2006-01-02"))
	if opt.OutConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter))
	}
	return zapcore.AddSync(fileWriter)
}

// GetLevelPriority 根据 zapcore.Level 获取 zap.LevelEnablerFunc
func GetLevelPriority(level zapcore.Level) zap.LevelEnablerFunc {
	switch level {
	case zapcore.DebugLevel:
		return func(level zapcore.Level) bool { // 调试级别
			return level == zap.DebugLevel
		}
	case zapcore.InfoLevel:
		return func(level zapcore.Level) bool { // 日志级别
			return level == zap.InfoLevel
		}
	case zapcore.WarnLevel:
		return func(level zapcore.Level) bool { // 警告级别
			return level == zap.WarnLevel
		}
	case zapcore.ErrorLevel:
		return func(level zapcore.Level) bool { // 错误级别
			return level == zap.ErrorLevel
		}
	case zapcore.DPanicLevel:
		return func(level zapcore.Level) bool { // dpanic级别
			return level == zap.DPanicLevel
		}
	case zapcore.PanicLevel:
		return func(level zapcore.Level) bool { // panic级别
			return level == zap.PanicLevel
		}
	case zapcore.FatalLevel:
		return func(level zapcore.Level) bool { // 终止级别
			return level == zap.FatalLevel
		}
	default:
		return func(level zapcore.Level) bool { // 调试级别
			return level == zap.DebugLevel
		}
	}
}

// ZapEncodeLevel 根据 EncodeLevel 返回 zapcore.LevelEncoder
func (opt *option) ZapEncodeLevel() zapcore.LevelEncoder {
	switch opt.EncodeLevel {
	case "LowercaseLevelEncoder": // 小写编码器(默认)
		return zapcore.LowercaseLevelEncoder
	case "LowercaseColorLevelEncoder": // 小写编码器带颜色
		return zapcore.LowercaseColorLevelEncoder
	case "CapitalLevelEncoder": // 大写编码器
		return zapcore.CapitalLevelEncoder
	case "CapitalColorLevelEncoder": // 大写编码器带颜色
		return zapcore.CapitalColorLevelEncoder
	default:
		return zapcore.LowercaseLevelEncoder
	}
}

// TransportLevel 根据字符串转化为 zapcore.Level
func TransportLevel(l string) zapcore.Level {
	level := strings.ToLower(l)
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.WarnLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel
	}
}
