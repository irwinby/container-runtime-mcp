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

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockSvc := imagemock.NewMockImageService(t)
			mockSvc.On("ListImages", mock.Anything, image.ListImagesParams{
				All:        tt.given.input.All,
				SharedSize: tt.given.input.SharedSize,
			}).Return(tt.given.result, tt.given.err)

			handler := NewToolsHandler(mockSvc)

			_, output, err := handler.ListImages(context.Background(), &mcp.CallToolRequest{}, tt.given.input)

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.True(t, tt.want.called)
			assert.Equal(t, tt.want.images, output.Images)
		})
	}
}
