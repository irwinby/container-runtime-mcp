package container

import (
	"context"
	"errors"
	"testing"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	services "github.com/irwinby/container-runtime-mcp/internal/service"
	containermock "github.com/irwinby/container-runtime-mcp/internal/service/container/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestServiceInspectContainer(t *testing.T) {
	type given struct {
		params InspectContainerParams
		result providers.ContainerInspect
		err    error
	}

	type want struct {
		called bool
		name   string
		result ContainerInspect
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"valid input": {
			given: given{params: NewInspectContainerParams().SetName("web"), result: providers.ContainerInspect{ID: "abc", Name: "web", State: "running"}},
			want:  want{called: true, name: "web", result: ContainerInspect{ID: "abc", Name: "web", State: "running"}},
		},
		"trimmed name": {
			given: given{params: NewInspectContainerParams().SetName("  web  "), result: providers.ContainerInspect{ID: "abc", Name: "web"}},
			want:  want{called: true, name: "web", result: ContainerInspect{ID: "abc", Name: "web"}},
		},
		"empty name": {
			given: given{params: NewInspectContainerParams().SetName(""), err: errors.New("validation error")},
			want:  want{},
		},
		"whitespace name": {
			given: given{params: NewInspectContainerParams().SetName("   "), err: errors.New("validation error")},
			want:  want{},
		},
		"provider error": {
			given: given{params: NewInspectContainerParams().SetName("web"), err: errors.New("docker error")},
			want:  want{called: true, name: "web"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := containermock.NewMockProviderClient(t)

			if tt.want.called {
				mockClient.On("InspectContainer", mock.Anything, providers.InspectContainerParams{
					Name: tt.want.name,
				}).Return(tt.given.result, tt.given.err)
			}

			service := NewService(mockClient, services.Policy{}, zap.NewNop())

			result, err := service.InspectContainer(context.Background(), tt.given.params)

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want.result, result)
		})
	}
}
