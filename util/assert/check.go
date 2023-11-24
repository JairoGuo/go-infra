package assert

import "regexp"

type Ruler interface {
	Rule() string
	DefaultMsg() string
}

func check(rule Ruler, value string) bool {
	return regexp.MustCompile(rule.Rule()).MatchString(value)
}
