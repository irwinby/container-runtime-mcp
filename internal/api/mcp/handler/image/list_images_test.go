package image

import (
	"context"
	"errors"
	"testing"

	imagemock "github.com/irwinby/container-runtime-mcp/internal/api/mcp/handler/image/mock"
	"github.com/irwinby/container-runtime-mcp/internal/service/image"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandlerListImages(t *testing.T) {
	type given struct {
		input  ListImagesInput
		result []image.Image
		err    error
	}

	type want struct {
		called bool
		images []ListImagesItem
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"valid input": {
			given: given{
				input: ListImagesInput{All: true},
				result: []image.Image{
					{ID: "abc123", RepoTags: []string{"nginx:latest"}, Size: 1024 * 1024},
				},
			},
			want: want{
				called: true,
				images: []ListImagesItem{
					{ID: "abc123", RepoTags: []string{"nginx:latest"}, Size: 1024 * 1024},
				},
			},
		},
		"empty result": {
			given: given{
				input:  ListImagesInput{},
				result: []image.Image{},
			},
			want: want{
				called: true,
				images: []ListImagesItem{},
			},
		},
		"service error": {
			given: given{
				input: ListImagesInput{},
				err:   errors.New("docker error"),
			},
			want: want{
				called: true,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockService := imagemock.NewMockImageService(t)

			mockService.On("ListImages", mock.Anything, image.ListImagesParams{
				All:        test.given.input.All,
				SharedSize: test.given.input.SharedSize,
			}).Return(test.given.result, test.given.err)

			handler := NewToolsHandler(mockService)

			_, output, err := handler.ListImages(context.Background(), &mcp.CallToolRequest{}, test.given.input)

			if test.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.True(t, test.want.called)
			assert.Equal(t, test.want.images, output.Images)
		})
	}
}
