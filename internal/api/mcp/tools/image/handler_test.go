package image

import (
	"testing"

	imagemock "github.com/irwinby/container-runtime-mcp/internal/api/mcp/tools/image/mock"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestToolsHandler_Register_CanWriteTrue(t *testing.T) {
	mockSvc := imagemock.NewMockImageService(t)
	mockSvc.On("CanWrite").Return(true).Once()

	handler := NewToolsHandler(mockSvc)

	server := mcp.NewServer(&mcp.Implementation{Name: "Test", Version: "1.0.0"}, &mcp.ServerOptions{})
	handler.Register(server)
}

func TestToolsHandler_Register_CanWriteFalse(t *testing.T) {
	mockSvc := imagemock.NewMockImageService(t)
	mockSvc.On("CanWrite").Return(false).Once()

	handler := NewToolsHandler(mockSvc)

	server := mcp.NewServer(&mcp.Implementation{Name: "Test", Version: "1.0.0"}, &mcp.ServerOptions{})
	handler.Register(server)
}
