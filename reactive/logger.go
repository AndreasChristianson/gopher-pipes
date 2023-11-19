package reactive

import (
	"fmt"
	"strconv"
)

type Level int

const (
	Verbose Level = iota
	Debug
	Info
	Warning
	Error
)

func (l Level) string() string {
	switch l {
	case Verbose:
		return "Verbose"
	case Debug:
		return "Debug"
	case Info:
		return "Info"
	case Warning:
		return "Warning"
	case Error:
		return "Error"
	}
	return strconv.Itoa(int(l))
}

var logger = func(level Level, args ...interface{}) {
	if level < Warning {
		return
	}
	fmt.Print(fmt.Sprintf("%s: ", level.string()))
	fmt.Println(args...)
}

func SetLogger(newLogger func(Level, ...interface{})) {
	logger = newLogger
}
