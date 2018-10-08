package log

import "log"

//go:generate stringer -type LogLevel

type LogLevel int
const (
	LOG_NONE LogLevel = iota
	LOG_DEBUG
	LOG_WARN
)
var Level = LOG_NONE

func Debug(args ...interface{}) {
	if Level >= LOG_DEBUG {
		log.Println(args...)
	}
}

func Warn(args ...interface{}) {
	if Level >= LOG_WARN {
		log.Println(args...)
	}
}

func OnError(err error) {
	log.Println(err)
}