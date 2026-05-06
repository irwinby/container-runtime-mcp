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

func TestServiceListContainers(t *testing.T) {
	type given struct {
		params ListContainersParams
		result []providers.Container
		err    error
	}

	type want struct {
		called bool
		result []Container
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"valid input": {
			given: given{params: NewListContainersParams().SetAll(true).SetLimit(10), result: []providers.Container{{ID: "abc", Names: []string{"web"}, Image: "nginx", State: "running", Status: "Up 1 hour"}}},
			want:  want{called: true, result: []Container{{ID: "abc", Names: []string{"web"}, Image: "nginx", State: "running", Status: "Up 1 hour"}}},
		},
		"provider error": {
			given: given{params: NewListContainersParams(), err: errors.New("docker error")},
			want:  want{called: true},
		},
		"negative limit": {
			given: given{params: NewListContainersParams().SetLimit(-1)},
			want:  want{called: false},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := containermock.NewMockProviderClient(t)

			if test.want.called {
				mockClient.On("ListContainers", mock.Anything, providers.ListContainersParams{
					All:    test.given.params.All,
					Limit:  test.given.params.Limit,
					Size:   test.given.params.Size,
					Latest: test.given.params.Latest,
				}).Return(test.given.result, test.given.err)
			}

			service := NewService(mockClient, services.Policy{}, zap.NewNop())

			result, err := service.ListContainers(context.Background(), test.given.params)

			if test.given.err != nil {
				require.Error(t, err)
				return
			}

			if !test.want.called {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.want.result, result)
		})
	}
}
