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

func TestServiceRemoveContainer(t *testing.T) {
	type given struct {
		params RemoveContainerParams
		err    error
	}

	type want struct {
		called        bool
		name          string
		force         bool
		removeVolumes bool
		removeLinks   bool
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"valid input": {
			given: given{params: NewRemoveContainerParams().SetName("web")},
			want:  want{called: true, name: "web"},
		},
		"trimmed name": {
			given: given{params: NewRemoveContainerParams().SetName("  web  ")},
			want:  want{called: true, name: "web"},
		},
		"with options": {
			given: given{params: NewRemoveContainerParams().SetName("web").SetForce(true).SetRemoveVolumes(true).SetRemoveLinks(true)},
			want:  want{called: true, name: "web", force: true, removeVolumes: true, removeLinks: true},
		},
		"empty name": {
			given: given{params: NewRemoveContainerParams().SetName(""), err: errors.New("validation error")},
			want:  want{},
		},
		"whitespace name": {
			given: given{params: NewRemoveContainerParams().SetName("   "), err: errors.New("validation error")},
			want:  want{},
		},
		"provider error": {
			given: given{params: NewRemoveContainerParams().SetName("web"), err: errors.New("docker error")},
			want:  want{called: true, name: "web"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := containermock.NewMockProviderClient(t)

			if test.want.called {
				mockClient.On("RemoveContainer", mock.Anything, providers.RemoveContainerParams{
					Name:          test.want.name,
					Force:         test.want.force,
					RemoveVolumes: test.want.removeVolumes,
					RemoveLinks:   test.want.removeLinks,
				}).Return(test.given.err)
			}

			service := NewService(mockClient, services.Policy{}, zap.NewNop())

			err := service.RemoveContainer(context.Background(), test.given.params)

			if test.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
