// @Author Jairo 2023/11/15 23:00
// @Email jairoguo@163.com

package log

type Logger interface {
	Info(msg string, args ...any)
	Debug(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Panic(msg string, args ...any)
	Fatal(msg string, args ...any)
}

var log Logger

// Use 绑定全局日志
func Use(l Logger) {
	log = l
}

func Info(msg string, args ...any) {
	log.Info(msg, args)
}

func Debug(msg string, args ...any) {
	log.Debug(msg, args)
}

func Warn(msg string, args ...any) {
	log.Warn(msg, args)
}

func Error(msg string, args ...any) {
	log.Error(msg, args)
}

func Panic(msg string, args ...any) {
	log.Panic(msg, args)
}

func Fatal(msg string, args ...any) {
	log.Fatal(msg, args)
}
