package result

import "github.com/jairoguo/go-infra/result/code"

type ResponseBody struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Data    any    `json:"data"`
	Msg     string `json:"msg"`
	Info    string `json:"info"`
	Success bool   `json:"success"`
}

type Option func(*ResponseBody)

func WithCode(code code.StatusCode) Option {
	return func(body *ResponseBody) {
		body.Code = int(code)
	}
}

func WithData(data any) Option {
	return func(body *ResponseBody) {
		body.Data = data
	}
}

func WithMsg(msg string) Option {
	return func(body *ResponseBody) {
		body.Msg = msg
	}
}

func WithStatus(status string) Option {
	return func(body *ResponseBody) {
		body.Status = status
	}
}

func WithInfo(info string) Option {
	return func(body *ResponseBody) {
		body.Info = info
	}
}

func WithSuccess(success bool) Option {
	return func(body *ResponseBody) {
		body.Success = success
	}
}

func With(opts ...Option) ResponseBody {
	body := &ResponseBody{}
	for _, opt := range opts {
		opt(body)
	}
	return *body
}

type ResponseBodyBuilder struct {
	body ResponseBody
}

func (builder *ResponseBodyBuilder) Code(code code.StatusCode) *ResponseBodyBuilder {
	builder.body.Code = int(code)
	return builder
}

func (builder *ResponseBodyBuilder) Data(data any) *ResponseBodyBuilder {
	builder.body.Data = data
	return builder
}

func (builder *ResponseBodyBuilder) Msg(msg string) *ResponseBodyBuilder {
	builder.body.Msg = msg
	return builder
}

func (builder *ResponseBodyBuilder) Status(status string) *ResponseBodyBuilder {
	builder.body.Status = status
	return builder
}

func (builder *ResponseBodyBuilder) Info(info string) *ResponseBodyBuilder {
	builder.body.Info = info
	return builder
}

func (builder *ResponseBodyBuilder) Success(success bool) *ResponseBodyBuilder {
	builder.body.Success = success
	return builder
}

func (builder *ResponseBodyBuilder) Build() ResponseBody {
	return builder.body
}

func Builder() *ResponseBodyBuilder {
	return &ResponseBodyBuilder{}
}

func OK() ResponseBody {
	return With(WithCode(code.OK), WithMsg(code.OK.Info()))
}

func OkWithData(data any) ResponseBody {
	return With(WithCode(code.OK), WithData(data), WithMsg(code.OK.Info()))
}

func Fail() ResponseBody {
	return With(WithCode(code.ERR), WithMsg(code.ERR.Info()))
}

func FailWithMessage(msg string) ResponseBody {
	return With(WithCode(code.ERR), WithMsg(msg))
}

func Code(code code.StatusCode) ResponseBody {
	return With(WithCode(code), WithMsg(code.Info()))
}
