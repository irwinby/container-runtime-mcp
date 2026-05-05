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

func TestServiceCreateContainer(t *testing.T) {
	type given struct {
		params CreateContainerParams
		id     string
		err    error
	}

	type want struct {
		called bool
		name   string
		image  string
		id     string
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"valid input": {
			given: given{params: NewCreateContainerParams().SetName("web").SetImage("nginx:latest"), id: "abc123"},
			want:  want{called: true, name: "web", image: "nginx:latest", id: "abc123"},
		},
		"trimmed name and image": {
			given: given{params: NewCreateContainerParams().SetName("  web  ").SetImage("  nginx:latest  "), id: "abc123"},
			want:  want{called: true, name: "web", image: "nginx:latest", id: "abc123"},
		},
		"empty name": {
			given: given{params: NewCreateContainerParams().SetName("").SetImage("nginx:latest"), err: errors.New("validation error")},
			want:  want{},
		},
		"whitespace name": {
			given: given{params: NewCreateContainerParams().SetName("   ").SetImage("nginx:latest"), err: errors.New("validation error")},
			want:  want{},
		},
		"empty image": {
			given: given{params: NewCreateContainerParams().SetName("web").SetImage(""), err: errors.New("validation error")},
			want:  want{},
		},
		"whitespace image": {
			given: given{params: NewCreateContainerParams().SetName("web").SetImage("   "), err: errors.New("validation error")},
			want:  want{},
		},
		"provider error": {
			given: given{params: NewCreateContainerParams().SetName("web").SetImage("nginx:latest"), err: errors.New("docker error")},
			want:  want{called: true, name: "web", image: "nginx:latest"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := containermock.NewMockProviderClient(t)

			if tt.want.called {
				mockClient.On("CreateContainer", mock.Anything, providers.CreateContainerParams{
					Name:  tt.want.name,
					Image: tt.want.image,
				}).Return(tt.given.id, tt.given.err)
			}

			service := NewService(mockClient, services.Policy{}, zap.NewNop())

			id, err := service.CreateContainer(context.Background(), tt.given.params)

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want.id, id)
		})
	}
}
