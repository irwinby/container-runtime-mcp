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

func TestServicePushImage(t *testing.T) {
	type given struct {
		params PushImageParams
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
			given: given{params: NewPushImageParams().SetRef("myapp:latest")},
			want:  want{called: true, ref: "myapp:latest"},
		},
		"trimmed ref": {
			given: given{params: NewPushImageParams().SetRef("  myapp:latest  ")},
			want:  want{called: true, ref: "myapp:latest"},
		},
		"with all and platform": {
			given: given{params: NewPushImageParams().SetRef("myapp").SetAll(true).SetPlatform(&v1.Platform{OS: "linux", Architecture: "arm64"})},
			want:  want{called: true, ref: "myapp", all: true, platform: &v1.Platform{OS: "linux", Architecture: "arm64"}},
		},
		"empty ref": {
			given: given{params: NewPushImageParams().SetRef(""), err: errors.New("validation error")},
			want:  want{},
		},
		"whitespace ref": {
			given: given{params: NewPushImageParams().SetRef("   "), err: errors.New("validation error")},
			want:  want{},
		},
		"provider error": {
			given: given{params: NewPushImageParams().SetRef("myapp:latest"), err: errors.New("docker error")},
			want:  want{called: true, ref: "myapp:latest"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := imagemock.NewMockProviderClient(t)

			if test.want.called {
				mockClient.On("PushImage", mock.Anything, providers.PushImageParams{
					Ref:      test.want.ref,
					All:      test.want.all,
					Platform: test.want.platform,
				}).Return(test.given.err)
			}

			service := NewService(mockClient, services.Policy{}, zap.NewNop())

			err := service.PushImage(context.Background(), test.given.params)

			if test.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.True(t, test.want.called)
		})
	}
}
