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

func TestHandlerTagImage(t *testing.T) {
	type given struct {
		input TagImageInput
		err   error
	}

	type want struct {
		called bool
		source string
		target string
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"valid input": {
			given: given{input: TagImageInput{Source: "nginx:latest", Target: "my-nginx:latest"}},
			want:  want{called: true, source: "nginx:latest", target: "my-nginx:latest"},
		},
		"whitespace input": {
			given: given{input: TagImageInput{Source: "  nginx:latest  ", Target: "  my-nginx:latest  "}},
			want:  want{called: true, source: "nginx:latest", target: "my-nginx:latest"},
		},
		"empty source": {
			given: given{input: TagImageInput{Source: "", Target: "my-nginx:latest"}, err: errors.New("validation error")},
			want:  want{},
		},
		"empty target": {
			given: given{input: TagImageInput{Source: "nginx:latest", Target: ""}, err: errors.New("validation error")},
			want:  want{},
		},
		"service error": {
			given: given{input: TagImageInput{Source: "nginx:latest", Target: "my-nginx:latest"}, err: errors.New("docker error")},
			want:  want{called: true, source: "nginx:latest", target: "my-nginx:latest"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockService := imagemock.NewMockImageService(t)

			if test.want.called {
				mockService.On("TagImage", mock.Anything, image.TagImageParams{
					Source: test.want.source,
					Target: test.want.target,
				}).Return(test.given.err)
			}

			handler := NewToolsHandler(mockService)

			_, _, err := handler.TagImage(context.Background(), &mcp.CallToolRequest{}, test.given.input)

			if test.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.True(t, test.want.called)
		})
	}
}
