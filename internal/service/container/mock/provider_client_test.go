package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockExists(t *testing.T) {
	m := NewMockProviderClient(t)
	assert.NotNil(t, m)
}
