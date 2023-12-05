package result

type ResponseBody struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Data    any    `json:"data"`
	Msg     string `json:"msg"`
	Info    string `json:"info"`
	Success bool   `json:"success"`
}

type Option func(*ResponseBody)

func WithCode(code int) Option {
	return func(body *ResponseBody) {
		body.Code = code
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

func (builder *ResponseBodyBuilder) Code(code int) *ResponseBodyBuilder {
	builder.body.Code = code
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
