package assert

import "errors"

func Is(value string, check Ruler) (bool, error) {
	return IsWithMsg(value, check, "")
}

func IsWithMsg(value string, ruler Ruler, msg string) (bool, error) {

	if check(ruler, value) == false {
		if msg == "" {
			return false, errors.New(ruler.DefaultMsg())
		}
		return false, errors.New(msg)
	}

	return true, nil

}
