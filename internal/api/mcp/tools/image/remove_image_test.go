package image

import (
	"context"
	"errors"
	"testing"

	imagemock "github.com/irwinby/container-runtime-mcp/internal/api/mcp/tools/image/mock"
	"github.com/irwinby/container-runtime-mcp/internal/service/image"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandlerRemoveImage(t *testing.T) {
	type given struct {
		input RemoveImageInput
		err   error
	}

	type want struct {
		called bool
		ref    string
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"valid input": {
			given: given{input: RemoveImageInput{Ref: "nginx:latest"}},
			want:  want{called: true, ref: "nginx:latest"},
		},
		"whitespace ref": {
			given: given{input: RemoveImageInput{Ref: "  nginx:latest  "}},
			want:  want{called: true, ref: "nginx:latest"},
		},
		"empty ref": {
			given: given{input: RemoveImageInput{Ref: ""}, err: errors.New("validation error")},
			want:  want{},
		},
		"service error": {
			given: given{input: RemoveImageInput{Ref: "nginx:latest"}, err: errors.New("docker error")},
			want:  want{called: true, ref: "nginx:latest"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockService := imagemock.NewMockImageService(t)

			if test.want.called {
				mockService.On("RemoveImage", mock.Anything, image.RemoveImageParams{
					Ref: test.want.ref,
				}).Return(test.given.err)
			}

			handler := NewToolsHandler(mockService)

			_, _, err := handler.RemoveImage(context.Background(), &mcp.CallToolRequest{}, test.given.input)

			if test.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.True(t, test.want.called)
		})
	}
}
