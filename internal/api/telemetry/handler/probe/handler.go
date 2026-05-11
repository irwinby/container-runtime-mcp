package probe

import (
	"context"
	"net/http"

	"github.com/irwinby/container-runtime-mcp/internal/service/system"
)

type systemService interface {
	Ping(ctx context.Context) (system.PingResult, error)
}

type Handler struct {
	systemService systemService
}

func NewHandler(systemService systemService) *Handler {
	return &Handler{
		systemService: systemService,
	}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/livez", h.Livez)
	mux.HandleFunc("/readyz", h.Readyz)
}
