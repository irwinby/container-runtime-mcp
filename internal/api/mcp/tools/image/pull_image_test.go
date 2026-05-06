package image

import (
	"context"
	"errors"
	"testing"

	imagemock "github.com/irwinby/container-runtime-mcp/internal/api/mcp/tools/image/mock"
	"github.com/irwinby/container-runtime-mcp/internal/service/image"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandlerPullImage(t *testing.T) {
	type given struct {
		input PullImageInput
		err   error
	}

	type want struct {
		called   bool
		ref      string
		all      bool
		platform *ocispec.Platform
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"valid input": {
			given: given{input: PullImageInput{Ref: "nginx:latest"}},
			want:  want{called: true, ref: "nginx:latest"},
		},
		"whitespace ref": {
			given: given{input: PullImageInput{Ref: "  nginx:latest  "}},
			want:  want{called: true, ref: "nginx:latest"},
		},
		"empty ref": {
			given: given{input: PullImageInput{Ref: ""}, err: errors.New("validation error")},
			want:  want{},
		},
		"with all and platform": {
			given: given{input: PullImageInput{Ref: "nginx", All: true, Platform: &ocispec.Platform{OS: "linux", Architecture: "amd64"}}},
			want:  want{called: true, ref: "nginx", all: true, platform: &ocispec.Platform{OS: "linux", Architecture: "amd64"}},
		},
		"service error": {
			given: given{input: PullImageInput{Ref: "nginx:latest"}, err: errors.New("docker error")},
			want:  want{called: true, ref: "nginx:latest"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockService := imagemock.NewMockImageService(t)

			if test.want.called {
				mockService.On("PullImage", mock.Anything, image.PullImageParams{
					Ref:      test.want.ref,
					All:      test.want.all,
					Platform: test.want.platform,
				}).Return(test.given.err)
			}

			handler := NewToolsHandler(mockService)

			_, _, err := handler.PullImage(context.Background(), &mcp.CallToolRequest{}, test.given.input)

			if test.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.True(t, test.want.called)
		})
	}
}
