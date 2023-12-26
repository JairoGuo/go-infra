// @Title
// @Description
// @Author Jairo 2023/12/26 10:20
// @Email jairoguo@163.com

package errors

type ErrorI18n interface {
	Tran(err *defaultError) string
}

var errI18n ErrorI18n
var enableErrI18n bool

func BindErrorI18n(i18n ErrorI18n) {
	enableErrI18n = true
	errI18n = i18n
}

type defaultErrorI18n struct {
}

func (i defaultErrorI18n) Tran(err *defaultError) string {
	return "Tran" + err.Key + err.error.Error()
}
