package status

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnauthorized(t *testing.T) {
	recorder := httptest.NewRecorder()
	Unauthorized(recorder, errors.New("unauthorized"))

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.Equal(t, "Bearer", recorder.Header().Get("WWW-Authenticate"))
	assert.Contains(t, recorder.Body.String(), "unauthorized")
}
