package trace

import "fmt"

//go:generate stringer -type LogLevel -output loglevel_string.go

type LogLevel int

const (
	LOG_NONE LogLevel = iota
	LOG_WARN
	LOG_INFO
	LOG_DEBUG
)

var Level = LOG_NONE

func Debug(args ...interface{}) {
	if Level >= LOG_DEBUG {
		fmt.Printf("%v|", LOG_DEBUG)
		fmt.Println(args...)
	}
}

func Info(args ...interface{}) {
	if Level >= LOG_INFO {
		fmt.Printf("%v|", LOG_INFO)
		fmt.Println(args...)
	}
}

func Warn(args ...interface{}) {
	if Level >= LOG_WARN {
		fmt.Printf("%v|", LOG_WARN)
		fmt.Println(args...)
	}
}

func OnError(err error) {
	fmt.Println(err)
}
