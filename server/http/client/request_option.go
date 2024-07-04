// @Title
// @Description
// @Author Jairo 2024/5/13 11:30
// @Email jairoguo@163.com

package http

import "sync"

type param struct {
	host        string
	port        int
	prefix      string
	healer      sync.Map // TODO 并发安全
	contentType string
	query       sync.Map
}

type RequestOption func(*param)

func WithHost(host string) RequestOption {
	return func(o *param) {
		o.host = host
	}
}

func WithPort(port int) RequestOption {
	return func(o *param) {
		o.port = port
	}
}

func WithPrefix(prefix string) RequestOption {
	return func(o *param) {
		o.prefix = prefix
	}

}

func WithHealer(key, value string) RequestOption {
	return func(o *param) {
		o.healer.Store(key, value)
	}
}

func WithContentType(contentType string) RequestOption {
	return func(o *param) {
		o.contentType = contentType
	}
}

func WithQuery(key, value string) RequestOption {
	return func(o *param) {
		o.query.Store(key, value)
	}
}

func (p *param) BuildRequestConfig(options ...RequestOption) *param {
	p.reset()
	option := p

	for _, opt := range options {
		opt(option)
	}

	if option.host != "" {
		p.host = option.host
	}

	if option.port != 0 {
		p.port = option.port
	}

	if option.prefix != "" {
		p.prefix = option.prefix
	}

	option.healer.Range(func(key, value interface{}) bool {
		p.healer.Store(key, value)
		return true

	})

	if option.contentType != "" {
		p.contentType = option.contentType
	} else {
		p.contentType = "application/json"

	}

	option.query.Range(func(key, value interface{}) bool {
		p.query.Store(key, value)
		return true
	})

	return option
}

func (p *param) reset() {

	p.host = ""
	p.port = 0
	p.prefix = ""
	p.healer.Range(func(key, value interface{}) bool {
		p.healer.Delete(key)
		return true
	})
	p.contentType = ""
	p.query.Range(func(key, value interface{}) bool {
		p.query.Delete(key)
		return true
	})
}
