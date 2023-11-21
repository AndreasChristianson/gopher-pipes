package reactive

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestLogging_EmitsMessages(t *testing.T) {
	var messages []string
	SetLogger(func(level Level, source interface{}, messageFormat string, args ...interface{}) {
		message := fmt.Sprintf(messageFormat, args...)
		messages = append(messages, fmt.Sprintf("%s [%s]: %s", level, source, message))
	})
	underTest := FromSlice([]string{"test"})
	underTest.Observe(func(string) error {
		return errors.New("expected error")
	})
	underTest.Observe(func(string) error {
		panic("expected panic")
	})
	underTest.UponClose(func() {

	})
	underTest.Start()
	underTest.AwaitCompletion()
	checkForLog(t, messages, Warning, "expected error")
	checkForLog(t, messages, Error, "expected panic")
	checkForLog(t, messages, Info, "Source is closed")
	checkForLog(t, messages, Debug, "Registering sink")
	checkForLog(t, messages, Verbose, "Beginning to send item (test)")
}

func checkForLog(t *testing.T, messages []string, level Level, s string) {
	for _, message := range messages {
		if strings.Contains(message, s) && strings.Contains(message, level.String()) {
			return
		}
	}
	assert.Fail(t, fmt.Sprintf("%s:*%s* not found in messages", level, s))
}
