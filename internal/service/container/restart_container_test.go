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

func TestServiceRestartContainer(t *testing.T) {
	type given struct {
		params RestartContainerParams
		err    error
	}

	type want struct {
		called bool
		name   string
		signal string
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"valid input": {
			given: given{params: NewRestartContainerParams().SetName("web")},
			want:  want{called: true, name: "web"},
		},
		"with signal": {
			given: given{params: NewRestartContainerParams().SetName("web").SetSignal("SIGKILL")},
			want:  want{called: true, name: "web", signal: "SIGKILL"},
		},
		"trimmed name": {
			given: given{params: NewRestartContainerParams().SetName("  web  ")},
			want:  want{called: true, name: "web"},
		},
		"empty name": {
			given: given{params: NewRestartContainerParams().SetName(""), err: errors.New("validation error")},
			want:  want{},
		},
		"whitespace name": {
			given: given{params: NewRestartContainerParams().SetName("   "), err: errors.New("validation error")},
			want:  want{},
		},
		"provider error": {
			given: given{params: NewRestartContainerParams().SetName("web"), err: errors.New("docker error")},
			want:  want{called: true, name: "web"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := containermock.NewMockProviderClient(t)

			if tt.want.called {
				mockClient.On("RestartContainer", mock.Anything, providers.RestartContainerParams{
					Name:           tt.want.name,
					Signal:         tt.want.signal,
					TimeoutSeconds: tt.given.params.TimeoutSeconds,
				}).Return(tt.given.err)
			}

			service := NewService(mockClient, services.Policy{}, zap.NewNop())

			err := service.RestartContainer(context.Background(), tt.given.params)

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
