package container

import (
	"context"
	"errors"
	"testing"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	services "github.com/irwinby/container-runtime-mcp/internal/service"
	containermock "github.com/irwinby/container-runtime-mcp/internal/service/container/mock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestServiceStartContainer(t *testing.T) {
	type given struct {
		params StartContainerParams
		err    error
	}

	type want struct {
		called bool
		name   string
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"valid input": {
			given: given{params: NewStartContainerParams().SetName("web")},
			want:  want{called: true, name: "web"},
		},
		"trimmed name": {
			given: given{params: NewStartContainerParams().SetName("  web  ")},
			want:  want{called: true, name: "web"},
		},
		"empty name": {
			given: given{params: NewStartContainerParams().SetName(""), err: errors.New("validation error")},
			want:  want{},
		},
		"whitespace name": {
			given: given{params: NewStartContainerParams().SetName("   "), err: errors.New("validation error")},
			want:  want{},
		},
		"provider error": {
			given: given{params: NewStartContainerParams().SetName("web"), err: errors.New("docker error")},
			want:  want{called: true, name: "web"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := containermock.NewMockProviderClient(t)

			if test.want.called {
				mockClient.On("StartContainer", mock.Anything, providers.StartContainerParams{
					Name: test.want.name,
				}).Return(test.given.err)
			}

			service := NewService(mockClient, services.Policy{}, zap.NewNop())

			err := service.StartContainer(context.Background(), test.given.params)

			if test.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
