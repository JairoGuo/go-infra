// @Author Jairo 2023/11/15 23:00
// @Email jairoguo@163.com

package log

type Logger interface {
	Info(value string, args ...any)
	Debug(value string, args ...any)
	Warn(value string, args ...any)
	Error(value string, args ...any)
	Panic(value string, args ...any)
	Fatal(value string, args ...any)
}

var log Logger

// Use 绑定全局日志
func Use(l Logger) {
	log = l
}
