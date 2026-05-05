package image

import (
	"context"
	"errors"
	"testing"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	services "github.com/irwinby/container-runtime-mcp/internal/service"
	imagemock "github.com/irwinby/container-runtime-mcp/internal/service/image/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestServiceListImages(t *testing.T) {
	type given struct {
		params ListImagesParams
		result []providers.Image
		err    error
	}

	type want struct {
		called bool
		result []Image
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"valid input": {
			given: given{params: NewListImagesParams().SetAll(true), result: []providers.Image{{ID: "abc", RepoTags: []string{"nginx:latest"}, Size: 1000}}},
			want:  want{called: true, result: []Image{{ID: "abc", RepoTags: []string{"nginx:latest"}, Size: 1000}}},
		},
		"provider error": {
			given: given{params: NewListImagesParams(), err: errors.New("docker error")},
			want:  want{called: true},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := imagemock.NewMockProviderClient(t)

			mockClient.On("ListImages", mock.Anything, providers.ListImagesParams{
				All:        tt.given.params.All,
				SharedSize: tt.given.params.SharedSize,
			}).Return(tt.given.result, tt.given.err)

			service := NewService(mockClient, services.Policy{}, zap.NewNop())

			result, err := service.ListImages(context.Background(), tt.given.params)

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.True(t, tt.want.called)
			assert.Equal(t, tt.want.result, result)
		})
	}
}
