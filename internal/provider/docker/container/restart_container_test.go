package container

import (
	"context"
	"testing"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	dockermock "github.com/irwinby/container-runtime-mcp/internal/provider/docker/mock"
	"github.com/moby/moby/client"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProviderRestartContainer(t *testing.T) {
	type given struct {
		params providers.RestartContainerParams
		err    error
	}

	type want struct {
		name    string
		signal  string
		timeout *int
	}

	timeout := 10

	tests := map[string]struct {
		given given
		want  want
	}{
		"success": {
			given: given{
				params: providers.RestartContainerParams{
					Name:           "web",
					Signal:         "SIGTERM",
					TimeoutSeconds: &timeout,
				},
			},
			want: want{
				name:    "web",
				signal:  "SIGTERM",
				timeout: &timeout,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := dockermock.NewMockDockerClient(t)

			mockClient.On("ContainerRestart", mock.Anything, tt.want.name, client.ContainerRestartOptions{
				Signal:  tt.want.signal,
				Timeout: tt.want.timeout,
			}).Return(client.ContainerRestartResult{}, tt.given.err)

			provider := NewProvider(mockClient, nopTimeout)

			err := provider.RestartContainer(context.Background(), tt.given.params)

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
