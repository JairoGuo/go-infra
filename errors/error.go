// @Title
// @Description
// @Author Jairo 2023/12/25 20:05
// @Email jairoguo@163.com

package errors

import (
	"errors"
	"fmt"
)

type Error interface {
	error
	WithKey(key string) Error
}

type defaultError struct {
	error      error
	Key        string
	EnableI18n bool
	i18n       ErrorI18n
}

type Option func(*defaultError)

func WithKey(key string) Option {
	return func(e *defaultError) {
		e.Key = key
	}
}

func WithEnableI18n() Option {
	return func(e *defaultError) {
		e.EnableI18n = true
	}
}

func WithBindI18n(i18n ErrorI18n) Option {
	return func(e *defaultError) {
		e.EnableI18n = true
		e.i18n = i18n
	}
}

func New(text string, opts ...Option) Error {
	e := &defaultError{
		error: errors.New(text),
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func (e *defaultError) Error() string {
	if e.EnableI18n || enableErrI18n {
		var businessError BusinessError
		switch {
		case errors.As(e, &businessError):
			if e.i18n != nil {
				return e.i18n.Tran(e)
			} else {
				if errI18n != nil {
					return errI18n.Tran(e)
				} else {
					_ = fmt.Errorf("not bound to internationalization")
					return e.error.Error()
				}
			}

		default:
			return e.error.Error()
		}

	}
	return e.error.Error()
}

func (e *defaultError) WithKey(key string) Error {
	e.Key = key
	return e
}

func Wrap(err error, text string) error {
	return fmt.Errorf(text+"error: %w", err)
}

func (e *defaultError) Unwrap() error {
	return e.error
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func (e *defaultError) Is(target error) bool {
	return errors.Is(e, target)
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func (e *defaultError) As(target any) bool {
	return errors.As(e, &target)
}

func As(err error, target any) bool {
	return errors.As(err, &target)
}
