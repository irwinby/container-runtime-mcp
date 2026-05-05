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

func TestHandlerPushImage(t *testing.T) {
	type given struct {
		input PushImageInput
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
			given: given{input: PushImageInput{Ref: "myapp:latest"}},
			want:  want{called: true, ref: "myapp:latest"},
		},
		"whitespace ref": {
			given: given{input: PushImageInput{Ref: "  myapp:latest  "}},
			want:  want{called: true, ref: "myapp:latest"},
		},
		"empty ref": {
			given: given{input: PushImageInput{Ref: ""}, err: errors.New("validation error")},
			want:  want{},
		},
		"with all and platform": {
			given: given{input: PushImageInput{Ref: "myapp", All: true, Platform: &ocispec.Platform{OS: "linux", Architecture: "arm64"}}},
			want:  want{called: true, ref: "myapp", all: true, platform: &ocispec.Platform{OS: "linux", Architecture: "arm64"}},
		},
		"service error": {
			given: given{input: PushImageInput{Ref: "myapp:latest"}, err: errors.New("docker error")},
			want:  want{called: true, ref: "myapp:latest"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockSvc := imagemock.NewMockImageService(t)
			if tt.want.called {
				mockSvc.On("PushImage", mock.Anything, image.PushImageParams{
					Ref:      tt.want.ref,
					All:      tt.want.all,
					Platform: tt.want.platform,
				}).Return(tt.given.err)
			}

			handler := NewToolsHandler(mockSvc)

			_, _, err := handler.PushImage(context.Background(), &mcp.CallToolRequest{}, tt.given.input)

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.True(t, tt.want.called)
		})
	}
}
