package reactive

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestMap_StartsImmediately(t *testing.T) {
	source := Just(123)
	mappedSource := Map(source, func(item int) (string, error) {
		asString := strconv.Itoa(item)
		return asString, nil
	})
	var result []string
	mappedSource.Observe(func(item string) error {
		result = append(result, item)
		return nil
	})
	source.Start()
	source.AwaitCompletion()
	assert.Equal(t, "123", result[0])
}

func TestMap_CallsUponClose(t *testing.T) {
	c := make(chan int)
	called := false
	source := FromChan(c)
	mappedSource := Map(source, func(item int) (string, error) {
		asString := strconv.Itoa(item)
		return asString, nil
	})
	mappedSource.Observe(func(s string) error {
		return nil
	})
	mappedSource.UponClose(func() {
		called = true
	})
	source.Start()
	close(c)
	source.AwaitCompletion()
	assert.True(t, called)
}

func TestMap_HandlesErrors(t *testing.T) {
	var messages []string
	SetLogger(func(level Level, source interface{}, messageFormat string, args ...interface{}) {
		message := fmt.Sprintf(messageFormat, args...)
		messages = append(messages, fmt.Sprintf("%s [%s]: %s", level, source, message))
	})
	c := make(chan string)
	source := FromChan(c)
	Map(source, func(item string) (string, error) {
		return "", errors.New("test error")
	})
	source.Start()
	c <- "test"
	close(c)
	source.AwaitCompletion()
	for _, message := range messages {
		if strings.Contains(message, "[test error]") && strings.Contains(message, "Warning") {
			return
		}
	}
	assert.Fail(t, "expected log not found")
}

func TestBuffer_StartsImmediately(t *testing.T) {
	source := Just(123)
	mappedSource := Buffer(source, 10)
	var result []int
	mappedSource.Observe(func(item int) error {
		result = append(result, item)
		return nil
	})
	source.Start()
	source.AwaitCompletion()
	assert.Equal(t, 123, result[0])
}

func TestBuffer_CanAddWithoutObserving(t *testing.T) {
	generatorCallCount := 0
	source := FromGenerator(func() (*int, error) {
		generatorCallCount++

		return &generatorCallCount, nil
	})
	buffered := Buffer[int](source, 10)
	buffered.Observe(func(item int) error {
		time.Sleep(10 * time.Millisecond)
		return nil
	})
	source.Start()
	time.Sleep(time.Millisecond)
	err := source.Cancel()
	assert.NoError(t, err)
	source.AwaitCompletion()
	// one in the source waiting to get into the chan, one in the sink waiting to sink, 10 in the buffer
	assert.Equal(t, 12, generatorCallCount)
}
