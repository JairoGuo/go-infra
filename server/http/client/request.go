// @Title
// @Description
// @Author Jairo 2024/5/7 18:02
// @Email jairoguo@163.com

package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
)

type Request struct {
	*http.Request
	Host   string
	Port   int
	Prefix string
	param
	build   url.URL
	bufPool sync.Pool
}

func NewRequest() *Request {
	r := &Request{}
	r.healer = sync.Map{}
	r.query = sync.Map{}
	r.bufPool = sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
	return r
}

func (r *Request) get(url string, options ...RequestOption) (*http.Request, error) {
	r.BuildRequestConfig(options...)
	url = r.getUrl(url)
	query := r.build.Query()
	r.query.Range(func(key, value interface{}) bool {
		query.Add(key.(string), value.(string))

		return true
	})

	r.build.RawQuery = query.Encode()

	request, err := http.NewRequest("GET", r.build.String(), nil)

	if err != nil {
		return nil, err
	}

	r.healer.Range(func(key, value interface{}) bool {
		request.Header.Set(key.(string), value.(string))
		return true

	})

	return request, nil
}

func (r *Request) post(url string, data interface{}, options ...RequestOption) (*http.Request, error) {

	r.BuildRequestConfig(options...)
	var buf io.Reader

	if d, ok := data.([]byte); ok {
		buf = bytes.NewReader(d)
	} else {
		encoder, err := r.Encoder(data)
		if err != nil {
			return nil, err
		}
		buf = encoder
	}

	url = r.getUrl(url)

	request, err := http.NewRequest("POST", url, buf)

	if err != nil {
		return nil, err
	}

	r.healer.Range(func(key, value interface{}) bool {
		request.Header.Set(key.(string), value.(string))
		return true

	})

	if r.contentType != "" {
		request.Header.Set("Content-Type", r.contentType)
	}

	return request, nil
}

func (r *Request) getUrl(url string) string {
	var host string
	var port int
	var prefix string
	if r.param.host != "" {
		host = r.param.host
	} else {
		host = r.Host
	}

	if r.param.port != 0 {
		port = r.param.port
	} else {
		port = r.Port
	}

	if r.param.prefix != "" {
		prefix = r.param.prefix
	} else {
		prefix = r.Prefix
	}

	r.build.RawQuery = ""
	r.build.Scheme = "http"
	r.build.Host = fmt.Sprintf("%s:%d", host, port)

	if prefix != "" {
		url = fmt.Sprintf("/%s/%s", prefix, url)
	}
	r.build.Path = url

	return r.build.String()
}
