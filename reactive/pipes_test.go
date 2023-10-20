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
	mappedSource.UponClose(func() {
		called = true
	})
	close(c)
	<-time.After(time.Millisecond)
	assert.True(t, called)
}
