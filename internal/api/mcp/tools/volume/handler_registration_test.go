package volume

import (
	"testing"

	volumemock "github.com/irwinby/container-runtime-mcp/internal/api/mcp/tools/volume/mock"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestToolsHandler_Register_CanWriteTrue(t *testing.T) {
	mockSvc := volumemock.NewMockVolumeService(t)
	mockSvc.On("CanWrite").Return(true).Once()

	handler := NewToolsHandler(mockSvc)

	server := mcp.NewServer(&mcp.Implementation{Name: "Test", Version: "1.0.0"}, &mcp.ServerOptions{})
	handler.Register(server)
}

func TestToolsHandler_Register_CanWriteFalse(t *testing.T) {
	mockSvc := volumemock.NewMockVolumeService(t)
	mockSvc.On("CanWrite").Return(false).Once()

	handler := NewToolsHandler(mockSvc)

	server := mcp.NewServer(&mcp.Implementation{Name: "Test", Version: "1.0.0"}, &mcp.ServerOptions{})
	handler.Register(server)
}
