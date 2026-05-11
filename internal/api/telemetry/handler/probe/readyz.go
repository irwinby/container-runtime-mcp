package probe

import (
	"fmt"
	"net/http"

	"github.com/irwinby/container-runtime-mcp/pkg/status"
)

func (h *Handler) Readyz(w http.ResponseWriter, r *http.Request) {
	_, err := h.systemService.Ping(r.Context())
	if err != nil {
		status.ServiceUnavailable(w, fmt.Errorf("service unavailable: %w", err))
		return
	}

	w.WriteHeader(http.StatusOK)
}
