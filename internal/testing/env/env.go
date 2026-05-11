// Package env provides test helpers for managing environment variables.
package env

import (
	"os"
	"testing"
)

// Unset unsets an environment variable and registers a cleanup that restores
// the original value after the test.
func Unset(t *testing.T, key string) {
	t.Helper()

	old, ok := os.LookupEnv(key)
	_ = os.Unsetenv(key)

	t.Cleanup(func() {
		if ok {
			_ = os.Setenv(key, old)
		} else {
			_ = os.Unsetenv(key)
		}
	})
}
