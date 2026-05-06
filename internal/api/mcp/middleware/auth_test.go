package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBearerAuth_Disabled(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := Auth("", next)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/mcp", nil)
	middleware.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestBearerAuth_Enabled(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := Auth("secret-token", next)

	tests := map[string]struct {
		authHeader string
		wantCode   int
		wantHeader string
	}{
		"no header": {
			authHeader: "",
			wantCode:   http.StatusUnauthorized,
			wantHeader: "Bearer",
		},
		"wrong scheme": {
			authHeader: "Basic YWRtaW46c2VjcmV0",
			wantCode:   http.StatusUnauthorized,
			wantHeader: "Bearer",
		},
		"missing token": {
			authHeader: "Bearer",
			wantCode:   http.StatusUnauthorized,
			wantHeader: "Bearer",
		},
		"wrong token": {
			authHeader: "Bearer wrong-token",
			wantCode:   http.StatusUnauthorized,
			wantHeader: "Bearer",
		},
		"correct token": {
			authHeader: "Bearer secret-token",
			wantCode:   http.StatusOK,
		},
		"correct token case insensitive scheme": {
			authHeader: "bearer secret-token",
			wantCode:   http.StatusOK,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/mcp", nil)

			if test.authHeader != "" {
				request.Header.Set("Authorization", test.authHeader)
			}

			middleware.ServeHTTP(recorder, request)

			require.Equal(t, test.wantCode, recorder.Code)

			if test.wantHeader != "" {
				assert.Equal(t, test.wantHeader, recorder.Header().Get("WWW-Authenticate"))
			}
		})
	}
}
