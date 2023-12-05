package code

type StatusCode int

type StatusInfo struct {
	Code string
	Info string
}

var statusCodeMap = make(map[int]StatusInfo)

const (
	ERR StatusCode = iota - 1
	OK
)

func init() {
	with(ERR, "E0000", "错误")
	with(OK, "S0000", "成功")
}

func with(cs StatusCode, code string, info string) {
	statusCodeMap[int(cs)] = StatusInfo{code, info}
}

func (code StatusCode) Value() StatusInfo {
	return statusCodeMap[int(code)]
}

func (code StatusCode) Code() string {
	return statusCodeMap[int(code)].Code
}

func (code StatusCode) Info() string {
	return statusCodeMap[int(code)].Info
}

func (code StatusCode) String() string {
	return statusCodeMap[int(code)].Code
}
