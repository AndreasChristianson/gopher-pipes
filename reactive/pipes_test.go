package reactive

import (
	"github.com/stretchr/testify/assert"
	"strconv"
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
	<-time.After(time.Millisecond)
	assert.True(t, called)
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

func TestBuffer_StartsCanAddWithoutObserving(t *testing.T) {
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
	<-time.After(time.Millisecond)
	// one in the source waiting to get into the chan, one in the sink waiting to sink, 10 in the buffer
	assert.Equal(t, 12, generatorCallCount)
	err := source.Cancel()
	assert.NoError(t, err)
}
