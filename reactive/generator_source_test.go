package reactive

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFromGenerator_HappyPath(t *testing.T) {
	responses := []string{"foobar", "test", "fizzbuzz"}
	pos := 0
	underTest := FromGeneratorWithDefaultBackoff(func() (*string, error) {
		if pos >= len(responses) {
			return nil, GeneratorFinished{}
		}

		ret := responses[pos]
		pos++
		return &ret, nil
	})
	results := make([]string, 0)
	underTest.Observe(func(item string) error {
		results = append(results, item)
		return nil
	})
	underTest.Start()
	underTest.AwaitCompletion()
	assert.Equal(t, "foobar", results[0])
	assert.Equal(t, "test", results[1])
	assert.Equal(t, "fizzbuzz", results[2])
}

func TestFromGenerator_Cancelable(t *testing.T) {
	underTest := FromGeneratorWithDefaultBackoff(func() (*string, error) {
		return nil, nil
	})

	underTest.Start()
	err := underTest.Cancel()
	assert.NoError(t, err)
	underTest.AwaitCompletion()
}

func TestFromGenerator_Backoff(t *testing.T) {
	callCount := 0
	underTest := FromGeneratorWithExponentialBackoff(func() (*string, error) {
		callCount++
		return nil, errors.New("")
	}, 100, 10)

	underTest.Start()
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, 3, callCount)
	err := underTest.Cancel()
	assert.NoError(t, err)
	underTest.AwaitCompletion()
}

func TestFromGenerator_NoBackoff(t *testing.T) {
	callCount := 0
	underTest := FromGenerator(func() (*string, error) {
		callCount++
		return nil, errors.New("")
	})

	underTest.Start()
	time.Sleep(100 * time.Millisecond)
	assert.Greater(t, callCount, 300)
	err := underTest.Cancel()
	assert.NoError(t, err)
	underTest.AwaitCompletion()
}

func TestFromGenerator_ClosingPreventsFurtherObservations(t *testing.T) {
	underTest := FromGenerator(func() (*string, error) {
		s := "fizzbuzz"
		return &s, nil
	})
	observeCount := 0
	underTest.Observe(func(item string) error {
		observeCount++
		time.Sleep(100 * time.Millisecond)
		return nil
	})

	underTest.Start()
	time.Sleep(time.Millisecond)
	err := underTest.Cancel()
	assert.NoError(t, err)
	underTest.AwaitCompletion()
	assert.Equal(t, 1, observeCount)
}
