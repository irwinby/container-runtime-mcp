package system

import (
	"testing"

	systemmock "github.com/irwinby/container-runtime-mcp/internal/api/mcp/tools/system/mock"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
)

func TestToolsHandler_Register(t *testing.T) {
	mockService := systemmock.NewMockSystemService(t)

	handler := NewToolsHandler(mockService)

	server := mcp.NewServer(&mcp.Implementation{Name: "Test", Version: "1.0.0"}, &mcp.ServerOptions{})
	handler.Register(server)

	assert.NotNil(t, handler)
	assert.NotNil(t, server)
}
