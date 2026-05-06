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

func TestServiceTagImage(t *testing.T) {
	type given struct {
		params TagImageParams
		err    error
	}

	type want struct {
		called bool
		source string
		target string
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"valid input": {
			given: given{params: NewTagImageParams().SetSource("nginx:latest").SetTarget("my-nginx:latest")},
			want:  want{called: true, source: "nginx:latest", target: "my-nginx:latest"},
		},
		"trimmed source and target": {
			given: given{params: NewTagImageParams().SetSource("  nginx:latest  ").SetTarget("  my-nginx:latest  ")},
			want:  want{called: true, source: "nginx:latest", target: "my-nginx:latest"},
		},
		"empty source": {
			given: given{params: NewTagImageParams().SetSource("").SetTarget("my-nginx:latest"), err: errors.New("validation error")},
			want:  want{},
		},
		"empty target": {
			given: given{params: NewTagImageParams().SetSource("nginx:latest").SetTarget(""), err: errors.New("validation error")},
			want:  want{},
		},
		"provider error": {
			given: given{params: NewTagImageParams().SetSource("nginx:latest").SetTarget("my-nginx:latest"), err: errors.New("docker error")},
			want:  want{called: true, source: "nginx:latest", target: "my-nginx:latest"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := imagemock.NewMockProviderClient(t)

			if test.want.called {
				mockClient.On("TagImage", mock.Anything, providers.TagImageParams{
					Source: test.want.source,
					Target: test.want.target,
				}).Return(test.given.err)
			}

			service := NewService(mockClient, services.Policy{}, zap.NewNop())

			err := service.TagImage(context.Background(), test.given.params)

			if test.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.True(t, test.want.called)
		})
	}
}
