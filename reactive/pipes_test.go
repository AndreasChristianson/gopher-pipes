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
		return strconv.Itoa(item), nil
	})
	var result []string
	mappedSource.Observe(func(item string) error {
		result = append(result, item)
		return nil
	})
	<-time.After(time.Millisecond)
	assert.Equal(t, "123", result[0])
}

func TestMap_CannotCancel(t *testing.T) {
	c := make(chan int)
	defer close(c)
	source := FromChan(c)
	mappedSource := Map(source, func(item int) (string, error) {
		return strconv.Itoa(item), nil
	})
	err := mappedSource.Cancel()
	assert.Error(t, err)
}

func TestMap_CallsUponClose(t *testing.T) {
	c := make(chan int)
	called := false
	source := FromChan(c)
	mappedSource := Map(source, func(item int) (string, error) {
		return strconv.Itoa(item), nil
	})
	mappedSource.Observe(func(s string) error {
		return nil
	})
	mappedSource.UponClose(func() {
		called = true
	})
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
	<-time.After(time.Millisecond)
	assert.Equal(t, 123, result[0])
}

func TestBuffer_CannotCancel(t *testing.T) {
	c := make(chan int)
	defer close(c)
	source := FromChan(c)
	mappedSource := Buffer(source, 10)
	err := mappedSource.Cancel()
	assert.Error(t, err)
}

func TestBuffer_StartsCanAddWithoutObserving(t *testing.T) {
	generatorCallCount := 0
	source := FromGenerator(func() *GeneratorResponse[int] {
		generatorCallCount++

		return &GeneratorResponse[int]{
			Data: &generatorCallCount,
		}
	})
	Buffer(source, 10)
	<-time.After(time.Millisecond)
	assert.Equal(t, 10+1, generatorCallCount)
	err := source.Cancel()
	assert.NoError(t, err)
}
