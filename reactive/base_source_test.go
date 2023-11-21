package reactive

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestBaseSource_HandlesSinkPanic(t *testing.T) {
	underTest := FromSlice([]string{"test"})
	underTest.Observe(func(s string) error {
		panic("test panic!")
	})
	underTest.Start()
	underTest.AwaitCompletion()
}

func TestBaseSource_DoesNotSendToSinkWhenClosing(t *testing.T) {
	lock := sync.Mutex{}
	lock.Lock()
	underTest := FromGenerator(func() (*bool, error) {
		lock.Lock()
		var t = true
		return &t, nil
	})
	observeCount := 0
	underTest.Observe(func(b bool) error {
		observeCount++
		return nil
	})
	underTest.Start()
	time.Sleep(time.Millisecond)
	err := underTest.Cancel()
	assert.Nil(t, err)
	lock.Unlock()
	underTest.AwaitCompletion()
	assert.Equal(t, 0, observeCount)
}

func TestBaseSource_HandlesHookPanic(t *testing.T) {
	underTest := FromSlice([]string{"test"})
	underTest.UponClose(func() {
		panic("test panic!")
	})
	underTest.Start()
	underTest.AwaitCompletion()
}

func TestBaseSource_HandlesSinkError(t *testing.T) {
	var messages []string
	SetLogger(func(level Level, source interface{}, messageFormat string, args ...interface{}) {
		message := fmt.Sprintf(messageFormat, args...)
		messages = append(messages, fmt.Sprintf("%s [%s]: %s", level, source, message))
	})
	underTest := FromSlice([]string{"test"})
	underTest.Observe(func(string) error {
		return errors.New("test error")
	})
	underTest.Start()
	underTest.AwaitCompletion()
	for _, message := range messages {
		if strings.Contains(message, "[test error]") && strings.Contains(message, "Warning") {
			return
		}
	}
	assert.Fail(t, "expected log not found")
}
