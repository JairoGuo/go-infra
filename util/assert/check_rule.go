package assert

type Rule string

const (
	URL               Rule = "^(((ht|f)tps?):\\/\\/)?([^!@#$%^&*?.\\s-]([^!@#$%^&*?.\\s]{0,63}[^!@#$%^&*?.\\s])?\\.)+[a-z]{2,6}\\/?"
	EMAIL             Rule = "^[A-Za-z0-9\\u4e00-\\u9fa5]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$"
	NUMBER            Rule = "^\\d+$"
	CHINESE_CHARACTER      = "^(?:[\\u3400-\\u4DB5\\u4E00-\\u9FEA\\uFA0E\\uFA0F\\uFA11\\uFA13\\uFA14\\uFA1F\\uFA21\\uFA23\\uFA24\\uFA27-\\uFA29]|[\\uD840-\\uD868\\uD86A-\\uD86C\\uD86F-\\uD872\\uD874-\\uD879][\\uDC00-\\uDFFF]|\\uD869[\\uDC00-\\uDED6\\uDF00-\\uDFFF]|\\uD86D[\\uDC00-\\uDF34\\uDF40-\\uDFFF]|\\uD86E[\\uDC00-\\uDC1D\\uDC20-\\uDFFF]|\\uD873[\\uDC00-\\uDEA1\\uDEB0-\\uDFFF]|\\uD87A[\\uDC00-\\uDFE0])+$"
	FILENAME          Rule = "\\.[^.\\\\/:*?\"<>|\\r\\n]+$"
)

func (r Rule) Rule() string {
	return string(r)
}

func (r Rule) DefaultMsg() string {
	switch r {
	case URL:
		return "不是URL"
	case EMAIL:
		return "非法邮箱格式"
	case NUMBER:
		return "不是数字"
	case CHINESE_CHARACTER:
		return "不是汉字"
	case FILENAME:
		return "不是带扩展名的文件名称"
	default:
		return "断言失败"
	}
}
