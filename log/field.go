// @Title
// @Description
// @Author Jairo 2023/12/19 10:22
// @Email jairoguo@163.com

package log

type Field struct {
	Key   string
	Value any
}

func Value(key string, value string) Field {
	return Field{key, value}
}

func Filed(key string, value any) {

}
