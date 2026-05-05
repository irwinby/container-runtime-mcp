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

func TestServicePullImage(t *testing.T) {
	type given struct {
		params PullImageParams
		err    error
	}

	type want struct {
		called   bool
		ref      string
		all      bool
		platform *v1.Platform
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"valid input": {
			given: given{params: NewPullImageParams().SetRef("nginx:latest")},
			want:  want{called: true, ref: "nginx:latest"},
		},
		"trimmed ref": {
			given: given{params: NewPullImageParams().SetRef("  nginx:latest  ")},
			want:  want{called: true, ref: "nginx:latest"},
		},
		"with all and platform": {
			given: given{params: NewPullImageParams().SetRef("nginx").SetAll(true).SetPlatform(&v1.Platform{OS: "linux", Architecture: "amd64"})},
			want:  want{called: true, ref: "nginx", all: true, platform: &v1.Platform{OS: "linux", Architecture: "amd64"}},
		},
		"empty ref": {
			given: given{params: NewPullImageParams().SetRef(""), err: errors.New("validation error")},
			want:  want{},
		},
		"whitespace ref": {
			given: given{params: NewPullImageParams().SetRef("   "), err: errors.New("validation error")},
			want:  want{},
		},
		"provider error": {
			given: given{params: NewPullImageParams().SetRef("nginx:latest"), err: errors.New("docker error")},
			want:  want{called: true, ref: "nginx:latest"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := imagemock.NewMockProviderClient(t)

			if tt.want.called {
				mockClient.On("PullImage", mock.Anything, providers.PullImageParams{
					Ref:      tt.want.ref,
					All:      tt.want.all,
					Platform: tt.want.platform,
				}).Return(tt.given.err)
			}

			service := NewService(mockClient, services.Policy{}, zap.NewNop())

			err := service.PullImage(context.Background(), tt.given.params)

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.True(t, tt.want.called)
		})
	}
}
