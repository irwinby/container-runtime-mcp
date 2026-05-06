package image

import (
	"context"
	"errors"
	"testing"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	services "github.com/irwinby/container-runtime-mcp/internal/service"
	imagemock "github.com/irwinby/container-runtime-mcp/internal/service/image/mock"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestServiceRemoveImage(t *testing.T) {
	type given struct {
		params RemoveImageParams
		err    error
	}

	type want struct {
		called        bool
		ref           string
		force         bool
		pruneChildren bool
		platform      *v1.Platform
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"valid input": {
			given: given{params: NewRemoveImageParams().SetRef("nginx:latest")},
			want:  want{called: true, ref: "nginx:latest"},
		},
		"with options": {
			given: given{params: NewRemoveImageParams().SetRef("nginx").SetForce(true).SetPruneChildren(true).SetPlatform(&v1.Platform{OS: "linux", Architecture: "amd64"})},
			want:  want{called: true, ref: "nginx", force: true, pruneChildren: true, platform: &v1.Platform{OS: "linux", Architecture: "amd64"}},
		},
		"trimmed ref": {
			given: given{params: NewRemoveImageParams().SetRef("  nginx:latest  ")},
			want:  want{called: true, ref: "nginx:latest"},
		},
		"empty ref": {
			given: given{params: NewRemoveImageParams().SetRef(""), err: errors.New("validation error")},
			want:  want{},
		},
		"whitespace ref": {
			given: given{params: NewRemoveImageParams().SetRef("   "), err: errors.New("validation error")},
			want:  want{},
		},
		"provider error": {
			given: given{params: NewRemoveImageParams().SetRef("nginx:latest"), err: errors.New("docker error")},
			want:  want{called: true, ref: "nginx:latest"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := imagemock.NewMockProviderClient(t)

			if test.want.called {
				mockClient.On("RemoveImage", mock.Anything, providers.RemoveImageParams{
					Ref:           test.want.ref,
					Force:         test.want.force,
					PruneChildren: test.want.pruneChildren,
					Platform:      test.want.platform,
				}).Return(test.given.err)
			}

			service := NewService(mockClient, services.Policy{}, zap.NewNop())

			err := service.RemoveImage(context.Background(), test.given.params)

			if test.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.True(t, test.want.called)
		})
	}
}
