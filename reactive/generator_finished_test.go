package reactive

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFinished_String(t *testing.T) {
	finished := &GeneratorFinished{}
	assert.NotEmpty(t, finished.Error())
}
