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

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := containermock.NewMockProviderClient(t)

			if tt.want.called {
				mockClient.On("ContainerLogs", mock.Anything, providers.ContainerLogsParams{
					Name:   tt.want.name,
					Stdout: tt.given.params.Stdout,
					Stderr: tt.given.params.Stderr,
					Since:  tt.given.params.Since,
					Tail:   tt.given.params.Tail,
				}).Return(tt.given.result, tt.given.err)
			}

			service := NewService(mockClient, services.Policy{}, zap.NewNop())

			result, err := service.ContainerLogs(context.Background(), tt.given.params)

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			if !tt.want.called {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want.result, result)
		})
	}
}
