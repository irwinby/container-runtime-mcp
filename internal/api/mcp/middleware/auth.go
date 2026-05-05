// Package middleware provides HTTP middleware for the MCP server.
package middleware

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"strings"

	"github.com/irwinby/container-runtime-mcp/internal/api/mcp/status"
)

// Auth returns an HTTP middleware that enforces bearer token authentication.
// If token is empty, the middleware is a no-op.
func Auth(token string, next http.Handler) http.Handler {
	if token == "" {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			status.Unauthorized(w, fmt.Errorf("authorization header is missing"))
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			status.Unauthorized(w, fmt.Errorf("authorization header is not valid"))
			return
		}

		bearerToken := strings.TrimSpace(parts[1])

		if subtle.ConstantTimeCompare([]byte(bearerToken), []byte(token)) != 1 {
			status.Unauthorized(w, fmt.Errorf("authorization token is not valid"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
