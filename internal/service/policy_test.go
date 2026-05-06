package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPolicy(t *testing.T) {
	policy := NewPolicy(true)
	assert.True(t, policy.ReadOnly)

	policy = NewPolicy(false)
	assert.False(t, policy.ReadOnly)
}

func TestPolicyIsWriteAllowed(t *testing.T) {
	t.Run("read only blocks writes", func(t *testing.T) {
		policy := NewPolicy(true)
		err := policy.IsWriteAllowed()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "read-only")
	})

	t.Run("writable allows writes", func(t *testing.T) {
		policy := NewPolicy(false)
		err := policy.IsWriteAllowed()
		require.NoError(t, err)
	})
}
