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

	type given struct {
		authHeader string
	}

	type want struct {
		code   int
		header string
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"no header": {
			given: given{authHeader: ""},
			want:  want{code: http.StatusUnauthorized, header: "Bearer"},
		},
		"wrong scheme": {
			given: given{authHeader: "Basic YWRtaW46c2VjcmV0"},
			want:  want{code: http.StatusUnauthorized, header: "Bearer"},
		},
		"missing token": {
			given: given{authHeader: "Bearer"},
			want:  want{code: http.StatusUnauthorized, header: "Bearer"},
		},
		"wrong token": {
			given: given{authHeader: "Bearer wrong-token"},
			want:  want{code: http.StatusUnauthorized, header: "Bearer"},
		},
		"correct token": {
			given: given{authHeader: "Bearer secret-token"},
			want:  want{code: http.StatusOK},
		},
		"correct token case insensitive scheme": {
			given: given{authHeader: "bearer secret-token"},
			want:  want{code: http.StatusOK},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/mcp", nil)

			if test.given.authHeader != "" {
				request.Header.Set("Authorization", test.given.authHeader)
			}

			middleware.ServeHTTP(recorder, request)

			require.Equal(t, test.want.code, recorder.Code)

			if test.want.header != "" {
				assert.Equal(t, test.want.header, recorder.Header().Get("WWW-Authenticate"))
			}
		})
	}
}
