// @Title
// @Description
// @Author Jairo 2024/5/11 10:07
// @Email jairoguo@163.com

package http

import (
	"errors"
	"net/http"
)

func IdentifyState(response *http.Response) error {
	//log.Debug("response", "URL", response.Request.URL, "Method", response.Request.Method, "Status", response.Status, "Body", response.Body)

	if response.StatusCode != http.StatusOK {
		return errors.New("response status is not 200")
	}

	return nil
}
