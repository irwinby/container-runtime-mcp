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

func TestHandlerInspectImage(t *testing.T) {
	type given struct {
		input  InspectImageInput
		result image.ImageInspect
		err    error
	}

	type want struct {
		called bool
		ref    string
		image  InspectImageDetails
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"valid input": {
			given: given{
				input:  InspectImageInput{Ref: "nginx:latest"},
				result: image.ImageInspect{ID: "abc", RepoTags: []string{"nginx:latest"}},
			},
			want: want{
				called: true,
				ref:    "nginx:latest",
				image:  InspectImageDetails{ID: "abc", RepoTags: []string{"nginx:latest"}},
			},
		},
		"empty ref": {
			given: given{input: InspectImageInput{Ref: ""}, err: errors.New("validation error")},
			want:  want{},
		},
		"service error": {
			given: given{
				input: InspectImageInput{Ref: "nginx:latest"},
				err:   errors.New("docker error"),
			},
			want: want{
				called: true,
				ref:    "nginx:latest",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockSvc := imagemock.NewMockImageService(t)
			if tt.want.called {
				mockSvc.On("InspectImage", mock.Anything, image.InspectImageParams{
					Ref: tt.want.ref,
				}).Return(tt.given.result, tt.given.err)
			}

			handler := NewToolsHandler(mockSvc)

			_, output, err := handler.InspectImage(context.Background(), &mcp.CallToolRequest{}, tt.given.input)

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.True(t, tt.want.called)
			assert.Equal(t, tt.want.image, output.Image)
		})
	}
}
