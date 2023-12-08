package config

import "strings"

type Mode int

const (
	DEFAULT Mode = iota
	TEST
	DEBUG
	RELEASE
)

func ParseMode(value string) Mode {

	switch strings.ToUpper(value) {
	case "TEST":
		return TEST
	case "DEBUG":
		return DEBUG
	case "RELEASE":
		return RELEASE
	default:
		return DEFAULT
	}

}

func (m Mode) String() string {
	switch m {
	case DEFAULT:
		return ""
	case TEST:
		return "test"
	case DEBUG:
		return "debug"
	case RELEASE:
		return "release"
	default:
		return ""
	}

}
