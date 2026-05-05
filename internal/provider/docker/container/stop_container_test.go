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

func TestProviderStopContainer(t *testing.T) {
	type given struct {
		params providers.StopContainerParams
		err    error
	}

	type want struct {
		name    string
		signal  string
		timeout *int
	}

	timeout := 30

	tests := map[string]struct {
		given given
		want  want
	}{
		"success": {
			given: given{
				params: providers.StopContainerParams{
					Name:           "web",
					Signal:         "SIGKILL",
					TimeoutSeconds: &timeout,
				},
			},
			want: want{
				name:    "web",
				signal:  "SIGKILL",
				timeout: &timeout,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := dockermock.NewMockDockerClient(t)

			mockClient.On("ContainerStop", mock.Anything, tt.want.name, client.ContainerStopOptions{
				Signal:  tt.want.signal,
				Timeout: tt.want.timeout,
			}).Return(client.ContainerStopResult{}, tt.given.err)

			provider := NewProvider(mockClient, nopTimeout)

			err := provider.StopContainer(context.Background(), tt.given.params)

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
