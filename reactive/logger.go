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

func (l Level) String() string {
	switch l {
	case Verbose:
		return "Verbose"
	case Debug:
		return "Debug  "
	case Info:
		return "Info   "
	case Warning:
		return "Warning"
	case Error:
		return "Error  "
	}
	return strconv.Itoa(int(l))
}

var logger = func(level Level, source string, args ...interface{}) {
	if level < Warning {
		return
	}
	fmt.Print(fmt.Sprintf("%s [%s]: ", level, source))
	fmt.Println(args...)
}

func SetLogger(newLogger func(Level, string, ...interface{})) {
	logger = newLogger
}
