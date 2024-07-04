// @Title
// @Description
// @Author Jairo 2024/5/13 11:20
// @Email jairoguo@163.com

package http

import (
	"errors"
	"net/http"
)

type Http struct {
	Config
	client *http.Client
	*Request
}

func New(options ...Option) *Http {
	h := &Http{}
	h.BuildHttpOption(options...)
	h.client = &http.Client{
		Timeout: h.Timeout,
	}
	h.Request = NewRequest()
	return h
}

func NewByConfig(config Config) *Http {
	h := &Http{}
	h.client = &http.Client{
		Timeout: config.Timeout,
	}
	h.Request = NewRequest()
	h.Config = config
	return h
}

func (h *Http) BindDestination(host string, port int) {
	h.Request.Host = host
	h.Request.Port = port
}

func (h *Http) GET(url string, options ...RequestOption) (*http.Response, error) {
	request, err := h.get(url, options...)

	if err != nil {
		return nil, err
	}

	response, err := h.client.Do(request)
	if err != nil {
		if h.EnableRetry {
			response, err = h.retryHandler(request, h.Retries, h.backoffStrategy)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err

		}

	}
	if response != nil && response.StatusCode != http.StatusOK {
		return response, errors.New(response.Status)
	}

	return response, nil
}

func (h *Http) POST(url string, data interface{}, options ...RequestOption) (*http.Response, error) {

	request, err := h.post(url, data, options...)
	if err != nil {
		return nil, err
	}

	response, err := h.client.Do(request)
	if err != nil {
		if h.EnableRetry {
			response, err = h.retryHandler(request, h.Retries, h.backoffStrategy)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err

		}
	}
	if response.StatusCode != http.StatusOK {
		return response, errors.New(response.Status)
	}
	return response, nil
}
