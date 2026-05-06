package image

import (
	"testing"

	imagemock "github.com/irwinby/container-runtime-mcp/internal/api/mcp/tools/image/mock"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestToolsHandler_Register_CanWriteTrue(t *testing.T) {
	mockService := imagemock.NewMockImageService(t)

	mockService.On("CanWrite").Return(true).Once()

	handler := NewToolsHandler(mockService)

	server := mcp.NewServer(&mcp.Implementation{Name: "Test", Version: "1.0.0"}, &mcp.ServerOptions{})
	handler.Register(server)
}

func TestToolsHandler_Register_CanWriteFalse(t *testing.T) {
	mockService := imagemock.NewMockImageService(t)

	mockService.On("CanWrite").Return(false).Once()

	handler := NewToolsHandler(mockService)

	server := mcp.NewServer(&mcp.Implementation{Name: "Test", Version: "1.0.0"}, &mcp.ServerOptions{})
	handler.Register(server)
}
