package debug

import "fmt"

//go:generate stringer -type LogLevel

type LogLevel int
const (
	LOG_NONE LogLevel = iota
	LOG_WARN
	LOG_DEBUG
)
var Level = LOG_NONE

func Debug(args ...interface{}) {
	if Level >= LOG_DEBUG {
		fmt.Println(args...)
	}
}

func Warn(args ...interface{}) {
	if Level >= LOG_WARN {
		fmt.Println(args...)
	}
}

func OnError(err error) {
	fmt.Println(err)
}