package reactive

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFromSlice_NoExtraData(t *testing.T) {
	underTest := FromSlice([]string{})
	underTest.Observe(func(s string) error {
		t.Fail()
		return nil
	})
	underTest.Start()
	underTest.AwaitCompletion()
}

func TestFromSlice_HappyPath(t *testing.T) {
	assertionsReached := false
	underTest := FromSlice([]string{
		"foobar",
		"pingpong",
		"fizzbuzz",
	})
	results := make([]string, 0)
	underTest.Observe(func(s string) error {
		results = append(results, s)
		return nil
	})
	underTest.UponClose(func() {
		assertionsReached = true
		assert.Equal(t, results[0], "foobar")
		assert.Equal(t, results[1], "pingpong")
		assert.Equal(t, results[2], "fizzbuzz")
	})
	underTest.Start()
	underTest.AwaitCompletion()
	assert.True(t, assertionsReached)
}

func TestJust_HappyPath(t *testing.T) {
	assertionsReached := false
	underTest := Just("foobar")
	results := make([]string, 0)
	underTest.Observe(func(s string) error {
		results = append(results, s)
		return nil
	})
	underTest.UponClose(func() {
		assertionsReached = true
		assert.Equal(t, results[0], "foobar")
	})
	underTest.Start()
	underTest.AwaitCompletion()
	assert.True(t, assertionsReached)
}
