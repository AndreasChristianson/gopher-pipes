package reactive

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChanSource_HappyPath(t *testing.T) {
	c := make(chan string)
	underTest := FromChan(c)
	results := make([]string, 0)
	underTest.Observe(func(s string) error {
		results = append(results, s)
		return nil
	})
	underTest.UponClose(func() {
		assert.Equal(t, results[0], "foobar")
		assert.Equal(t, results[1], "test")
		assert.Equal(t, results[2], "fizzbuzz")
	})
	underTest.Start()
	c <- "foobar"
	c <- "test"
	c <- "fizzbuzz"
	close(c)
	underTest.AwaitCompletion()
}
func TestChanSource_CallsUponClose(t *testing.T) {
	c := make(chan string)
	underTest := FromChan(c)
	uponCloseCalled := false
	underTest.UponClose(func() {
		uponCloseCalled = true
	})
	underTest.Start()
	close(c)
	underTest.AwaitCompletion()
	assert.True(t, uponCloseCalled)
}

func TestChanSource_Realtime(t *testing.T) {
	c := make(chan string)
	underTest := FromChan(c)
	results := make([]string, 0)
	underTest.Observe(func(s string) error {
		results = append(results, s)
		return nil
	})
	underTest.Start()
	c <- "foobar"
	c <- "test"
	c <- "fizzbuzz"
	close(c)
	underTest.AwaitCompletion()
	assert.Equal(t, results[0], "foobar")
	assert.Equal(t, results[1], "test")
	assert.Equal(t, results[2], "fizzbuzz")

}

func TestChanSource_BufferedRealtime(t *testing.T) {
	c := make(chan string, 3)
	underTest := FromChan(c)
	results := make([]string, 0)
	underTest.Observe(func(s string) error {
		results = append(results, s)
		return nil
	})
	underTest.Start()
	c <- "foobar"
	c <- "test"
	c <- "fizzbuzz"
	close(c)
	underTest.AwaitCompletion()
	assert.Equal(t, results[0], "foobar")
	assert.Equal(t, results[1], "test")
	assert.Equal(t, results[2], "fizzbuzz")
}
func TestChanSource_HandlesPreStartedChannels(t *testing.T) {
	c := make(chan string, 10)
	c <- "foobar"
	c <- "test"
	c <- "fizzbuzz"

	underTest := FromChan(c)

	var results []string
	underTest.Observe(func(s string) error {
		results = append(results, s)
		return nil
	})
	underTest.Start()
	close(c)
	underTest.AwaitCompletion()
	assert.Equal(t, results[0], "foobar")
	assert.Equal(t, results[1], "test")
	assert.Equal(t, results[2], "fizzbuzz")
}
