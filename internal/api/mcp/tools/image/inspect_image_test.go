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

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockService := imagemock.NewMockImageService(t)

			if test.want.called {
				mockService.On("InspectImage", mock.Anything, image.InspectImageParams{
					Ref: test.want.ref,
				}).Return(test.given.result, test.given.err)
			}

			handler := NewToolsHandler(mockService)

			_, output, err := handler.InspectImage(context.Background(), &mcp.CallToolRequest{}, test.given.input)

			if test.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.True(t, test.want.called)
			assert.Equal(t, test.want.image, output.Image)
		})
	}
}
