// @Title
// @Description
// @Author Jairo 2023/11/16 22:33
// @Email jairoguo@163.com

package log

type Level int

const (
	INFO Level = iota
	DEBUG
	WARN
	ERROR
	PANIC
	FATAL
)

func (l Level) String() string {

	switch l {

	case INFO:
		return "info"
	case DEBUG:
		return "debug"
	case WARN:
		return "warn"
	case ERROR:
		return "error"
	case PANIC:
		return "panic"
	case FATAL:
		return "fatal"
	default:
		return "info"
	}
}
