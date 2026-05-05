package container

import (
	"context"
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
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := dockermock.NewMockDockerClient(t)

			mockClient.On("ContainerInspect", mock.Anything, tt.want.name, client.ContainerInspectOptions{}).Return(client.ContainerInspectResult{
				Container: container.InspectResponse{
					ID:   tt.want.id,
					Name: tt.want.cname,
					State: &container.State{
						Status: container.ContainerState(tt.want.state),
					},
				},
			}, tt.given.err)

			provider := NewProvider(mockClient, nopTimeout)

			result, err := provider.InspectContainer(context.Background(), tt.given.params)

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want.id, result.ID)
			require.Equal(t, tt.want.cname, result.Name)
			require.Equal(t, tt.want.state, result.State)
		})
	}
}
