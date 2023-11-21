package reactive

import (
	"fmt"
	"log"
	"strconv"
)

// Level represents a log level ranging from Verbose to Error.
//   - Verbose logging indicates intent and logs every item handled.
//   - Debug logging indicates state changes and logs internal fields.
//   - Info logging happensf for once per Source status changes and handled errors.
//   - Warning logging occurs when a function returns an error resulting in a message being discarded.
//   - Error logging occurs when a source shuts down unexpectedly without cleaning up.
//
// The default logger logs Warning and Error level. Use SetLogger to change this behavior.
type Level int

const (
	// Verbose logging indicates intent and logs every item handled. For example, logging the intent to send an item to a sink
	Verbose Level = iota
	// Debug logging indicates state changes and logs internal fields. For example a generator's error count increasing or an observer (Source.Observe) being registered.
	Debug
	// Info logging occurs for once per Source status changes and handled errors. For example a source successfully shuts down, or an error is encountered and handled
	Info
	// Warning logging occurs when a function returns an error resulting in a message being discarded. For example when a sink returns an error.
	Warning
	// Error logging occurs when a sink, shutdown hook or mapper panics.
	Error
)

// String() returns the name of the log Level right padded with spaces to seven characters.
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

var logger = func(level Level, source interface{}, formatString string, args ...interface{}) {
	if level < Warning {
		return
	}
	message := fmt.Sprintf(formatString, args...) // delay message formatting till level can be evaluated
	log.Printf("%s [%s]: %s\n", level, source, message)
}

// SetLogger sets the logger for all logging in the reactive package. Default logger implementation:
//
//	func(level Level, source interface{}, formatString string, args ...interface{}) {
//	  if level < Warning {
//	    return
//	  }
//	  message := fmt.Sprintf(formatString, args...) // delay message formatting till level can be evaluated
//	  log.Printf("%s [%s]: %s\n", level, source, message)
//	}
//
// Note that the formatting string should not be evaluated until after filtering by log [Level].
// Formatting Verbose and Debug logs can create superfluous cpu load.
func SetLogger(newLogger func(Level, interface{}, string, ...interface{})) {
	logger = newLogger
}
