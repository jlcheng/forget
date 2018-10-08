package log

import "log"

//go:generate stringer -type LogLevel

type LogLevel int
const (
	LOG_NONE LogLevel = iota
	LOG_DEBUG
	LOG_WARN
)

type Logger struct {
	Level LogLevel
}

func (s Logger) Debug(args... string) {
	if s.Level >= LOG_DEBUG {
		log.Println(args)
	}
}

func (s Logger) Error(err error) {
	log.Println(err)
}