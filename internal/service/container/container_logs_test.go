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

func TestServiceContainerLogs(t *testing.T) {
	type given struct {
		params ContainerLogsParams
		result providers.ContainerLogsResult
		err    error
	}

	type want struct {
		called bool
		name   string
		result ContainerLogsResult
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"valid input": {
			given: given{
				params: NewContainerLogsParams().SetName("web"),
				result: providers.NewContainerLogsResult().SetStdout("hello").SetStderr("warn"),
			},
			want: want{
				called: true,
				name:   "web",
				result: ContainerLogsResult{Stdout: "hello", Stderr: "warn"},
			},
		},
		"trimmed name": {
			given: given{
				params: NewContainerLogsParams().SetName("  web  "),
				result: providers.NewContainerLogsResult().SetStdout("hello"),
			},
			want: want{
				called: true,
				name:   "web",
				result: ContainerLogsResult{Stdout: "hello"},
			},
		},
		"empty name": {
			given: given{params: NewContainerLogsParams().SetName("")},
			want:  want{},
		},
		"whitespace name": {
			given: given{params: NewContainerLogsParams().SetName("   ")},
			want:  want{},
		},
		"both stdout and stderr false": {
			given: given{
				params: NewContainerLogsParams().SetName("web").SetStdout(false).SetStderr(false),
			},
			want: want{},
		},
		"provider error": {
			given: given{
				params: NewContainerLogsParams().SetName("web"),
				err:    errors.New("docker error"),
			},
			want: want{called: true, name: "web"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := containermock.NewMockProviderClient(t)

			if test.want.called {
				mockClient.On("ContainerLogs", mock.Anything, providers.ContainerLogsParams{
					Name:   test.want.name,
					Stdout: test.given.params.Stdout,
					Stderr: test.given.params.Stderr,
					Since:  test.given.params.Since,
					Tail:   test.given.params.Tail,
				}).Return(test.given.result, test.given.err)
			}

			service := NewService(mockClient, services.Policy{}, zap.NewNop())

			result, err := service.ContainerLogs(context.Background(), test.given.params)

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
