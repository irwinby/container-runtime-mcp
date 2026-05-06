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

func TestServiceInspectImage(t *testing.T) {
	type given struct {
		params InspectImageParams
		result providers.ImageInspect
		err    error
	}

	type want struct {
		called bool
		ref    string
		result ImageInspect
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"valid input": {
			given: given{params: NewInspectImageParams().SetRef("nginx:latest"), result: providers.ImageInspect{ID: "abc", RepoTags: []string{"nginx:latest"}}},
			want:  want{called: true, ref: "nginx:latest", result: ImageInspect{ID: "abc", RepoTags: []string{"nginx:latest"}}},
		},
		"trimmed ref": {
			given: given{params: NewInspectImageParams().SetRef("  nginx:latest  "), result: providers.ImageInspect{ID: "abc"}},
			want:  want{called: true, ref: "nginx:latest", result: ImageInspect{ID: "abc"}},
		},
		"empty ref": {
			given: given{params: NewInspectImageParams().SetRef(""), err: errors.New("validation error")},
			want:  want{},
		},
		"whitespace ref": {
			given: given{params: NewInspectImageParams().SetRef("   "), err: errors.New("validation error")},
			want:  want{},
		},
		"provider error": {
			given: given{params: NewInspectImageParams().SetRef("nginx:latest"), err: errors.New("docker error")},
			want:  want{called: true, ref: "nginx:latest"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := imagemock.NewMockProviderClient(t)

			if test.want.called {
				mockClient.On("InspectImage", mock.Anything, providers.InspectImageParams{
					Ref: test.want.ref,
				}).Return(test.given.result, test.given.err)
			}

			service := NewService(mockClient, services.Policy{}, zap.NewNop())

			result, err := service.InspectImage(context.Background(), test.given.params)

			if test.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.True(t, test.want.called)
			assert.Equal(t, test.want.result, result)
		})
	}
}
