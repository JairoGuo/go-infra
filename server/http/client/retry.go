// @Title
// @Description
// @Author Jairo 2024/5/13 11:15
// @Email jairoguo@163.com

package http

import (
	"bytes"
	"edge-side-client/core/log"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type BackoffStrategy func(retry int) time.Duration

func (h *Http) retryHandler(req *http.Request, maxRetries int, backoffStrategy BackoffStrategy) (*http.Response, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	var (
		originalBody []byte
		err          error
	)
	if req.Body != nil {
		originalBody, err = copyBody(req.Body)
		resetBody(req, originalBody)
	}
	if err != nil {
		return nil, err
	}
	attemptLimit := maxRetries
	if attemptLimit <= 0 {
		attemptLimit = 1
	}

	var resp *http.Response
	//重试次数
	for i := 1; i <= attemptLimit; i++ {
		resp, err = h.client.Do(req)
		if err != nil {
			log.Debug("retry error: ", err)
		}
		// 重试 500 以上的错误码
		if err == nil && resp.StatusCode < 500 {
			return resp, err
		}
		// 如果正在重试，那么释放fd
		if resp != nil {
			resp.Body.Close()
		}
		// 重置body
		if req.Body != nil {
			resetBody(req, originalBody)
		}
		time.Sleep(backoffStrategy(i) + time.Millisecond)
	}

	return resp, err
}

func (h *Http) BindBackoffStrategy(strategy BackoffStrategy) {
	h.backoffStrategy = strategy
}

func copyBody(src io.ReadCloser) ([]byte, error) {
	b, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, err
	}
	err = src.Close()
	if err != nil {
		return nil, err
	}
	return b, nil
}

func resetBody(request *http.Request, originalBody []byte) {
	request.Body = io.NopCloser(bytes.NewBuffer(originalBody))
	request.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewBuffer(originalBody)), nil
	}
}

func (h *Http) retryHedged(req *http.Request, maxRetries int, backoffStrategy BackoffStrategy) (*http.Response, error) {
	var (
		originalBody []byte
		err          error
	)
	if req != nil && req.Body != nil {
		originalBody, err = copyBody(req.Body)
	}
	if err != nil {
		return nil, err
	}

	AttemptLimit := maxRetries
	if AttemptLimit <= 0 {
		AttemptLimit = 1
	}

	// 每次请求copy新的request
	copyRequest := func() (request *http.Request) {
		request = req.Clone(req.Context())
		if request.Body != nil {
			resetBody(request, originalBody)
		}
		return
	}

	multiplexCh := make(chan struct {
		resp  *http.Response
		err   error
		retry int
	})

	totalSentRequests := &sync.WaitGroup{}
	allRequestsBackCh := make(chan struct{})
	go func() {
		totalSentRequests.Wait()
		close(allRequestsBackCh)
	}()
	var resp *http.Response
	for i := 1; i <= AttemptLimit; i++ {
		totalSentRequests.Add(1)
		go func() {
			// 标记已经执行完
			defer totalSentRequests.Done()
			req = copyRequest()
			resp, err = h.client.Do(req)
			if err != nil {
				log.Debug("retry hedged error : ", err)
			}
			// 重试 500 以上的错误码
			if err == nil && resp.StatusCode < 500 {
				multiplexCh <- struct {
					resp  *http.Response
					err   error
					retry int
				}{resp: resp, err: err, retry: i}
				return
			}
			// 如果正在重试，那么释放fd
			if resp != nil {
				resp.Body.Close()
			}
			// 重置body
			if req.Body != nil {
				resetBody(req, originalBody)
			}
			time.Sleep(backoffStrategy(i) + 1*time.Microsecond)
		}()
	}

	select {
	case res := <-multiplexCh:
		return res.resp, res.err
	case <-allRequestsBackCh:
		// 到这里，说明全部的 goroutine 都执行完毕，但是都请求失败了
		return nil, errors.New("all req finish，but all fail")
	}
}
