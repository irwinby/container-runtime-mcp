package container

import (
	"context"
	"errors"
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
		"success with signal and timeout": {
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
		"success without signal and timeout": {
			given: given{
				params: providers.RestartContainerParams{
					Name: "web",
				},
			},
			want: want{
				name: "web",
			},
		},
		"error": {
			given: given{
				params: providers.RestartContainerParams{
					Name: "web",
				},
				err: errors.New("docker error"),
			},
			want: want{
				name: "web",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := dockermock.NewMockDockerClient(t)

			mockClient.On("ContainerRestart", mock.Anything, test.want.name, client.ContainerRestartOptions{
				Signal:  test.want.signal,
				Timeout: test.want.timeout,
			}).Return(client.ContainerRestartResult{}, test.given.err)

			provider := NewProvider(mockClient, nopTimeout)

			err := provider.RestartContainer(context.Background(), test.given.params)

			if test.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
