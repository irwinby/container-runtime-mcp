package container

import (
	"context"
	"errors"
	"testing"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	dockermock "github.com/irwinby/container-runtime-mcp/internal/provider/docker/mock"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProviderInspectContainer(t *testing.T) {
	type given struct {
		params providers.InspectContainerParams
		err    error
	}

	type want struct {
		name  string
		id    string
		cname string
		state string
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"success": {
			given: given{
				params: providers.InspectContainerParams{Name: "web"},
			},
			want: want{
				name:  "web",
				id:    "c1",
				cname: "/web",
				state: "running",
			},
		},
		"error": {
			given: given{
				params: providers.InspectContainerParams{Name: "web"},
				err:    errors.New("docker error"),
			},
			want: want{
				name: "web",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := dockermock.NewMockDockerClient(t)

			mockClient.On("ContainerInspect", mock.Anything, test.want.name, client.ContainerInspectOptions{}).Return(client.ContainerInspectResult{
				Container: container.InspectResponse{
					ID:   test.want.id,
					Name: test.want.cname,
					State: &container.State{
						Status: container.ContainerState(test.want.state),
					},
				},
			}, test.given.err)

			provider := NewProvider(mockClient, nopTimeout)

			result, err := provider.InspectContainer(context.Background(), test.given.params)

			if test.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want.id, result.ID)
			require.Equal(t, test.want.cname, result.Name)
			require.Equal(t, test.want.state, result.State)
		})
	}
}
