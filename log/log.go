// @Author Jairo 2023/11/15 23:00
// @Email jairoguo@163.com

package log

type Logger interface {
	Info(msg string, fields []Field)
	Debug(msg string, fields []Field)
	Warn(msg string, fields []Field)
	Error(msg string, fields []Field)
	Panic(msg string, fields []Field)
	Fatal(msg string, fields []Field)
}

var log Logger

// Use 绑定全局日志
func Use(l Logger) {
	log = l
}

func Info(msg string, args ...any) {
	fields := handleArgs(args...)
	log.Info(msg, fields)
}

func Debug(msg string, args ...any) {
	fields := handleArgs(args...)
	log.Debug(msg, fields)
}

func Warn(msg string, args ...any) {
	fields := handleArgs(args...)
	log.Warn(msg, fields)
}

func Error(msg string, args ...any) {
	fields := handleArgs(args...)
	log.Error(msg, fields)
}

func Panic(msg string, args ...any) {
	fields := handleArgs(args...)
	log.Panic(msg, fields)
}

func Fatal(msg string, args ...any) {
	fields := handleArgs(args...)
	log.Fatal(msg, fields)
}
