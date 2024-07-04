// @Title
// @Description
// @Author Jairo 2024/5/11 10:38
// @Email jairoguo@163.com

package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

func (r *Request) Encoder(data interface{}) (*bytes.Buffer, error) {
	buf := r.bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	err := json.NewEncoder(buf).Encode(data)
	if err != nil {
		r.bufPool.Put(buf)
		return nil, err
	}
	r.bufPool.Put(buf)
	return buf, nil
}

func (r *Request) Decoder(response *http.Response, data interface{}) error {
	if response == nil {
		return errors.New("response is nil")
	}
	defer response.Body.Close()
	err := json.NewDecoder(response.Body).Decode(&data)

	return err
}
