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

var logger = func(level Level, source string, formatString string, args ...interface{}) {
	if level < Warning {
		return
	}
	message := fmt.Sprintf(formatString, args...) // delay message formatting till level can be evaluated
	fmt.Print(fmt.Sprintf("%s [%s]: %s", level, source, message))
}

func SetLogger(newLogger func(Level, string, string, ...interface{})) {
	logger = newLogger
}
