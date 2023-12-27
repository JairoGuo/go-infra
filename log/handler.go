// @Title
// @Description
// @Author Jairo 2023/11/15 23:41
// @Email jairoguo@163.com

package log

import (
	"fmt"
	"strings"
)

type Handler struct {
}

func handle(log Logger, level Level, msg string, args ...any) {

	switch level {

	}
}

func handleMassage(msg string, args ...any) string {

	format := strings.Replace(msg, "{}", "%v", 0)
	message := fmt.Sprintf(format, args...)
	return message
}

func handleArgs(args ...any) []Field {
	var argsList []Field

	var preValue any

	for _, v := range args {
		switch v.(type) {
		case Field:
			field := v.(Field)
			if preValue != nil {
				argsList = append(argsList, Field{Value: preValue})
				preValue = nil
			}
			argsList = append(argsList, field)

		case string:
			if preValue != nil {
				argsList = append(argsList, Field{Key: preValue.(string), Value: v})
				preValue = nil
			} else {
				preValue = v
			}
		case any:
			if preValue != nil {
				argsList = append(argsList, Field{Key: preValue.(string), Value: v})
				preValue = nil
			} else {
				argsList = append(argsList, Field{Value: v})
			}

		}
	}

	return argsList
}
