package ptr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBool(t *testing.T) {
	trueValue := Bool(true)
	assert.NotNil(t, trueValue)
	assert.True(t, *trueValue)

	falseValue := Bool(false)
	assert.NotNil(t, falseValue)
	assert.False(t, *falseValue)
}
