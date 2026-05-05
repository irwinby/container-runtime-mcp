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

func TestServiceExecContainer(t *testing.T) {
	type given struct {
		params ExecContainerParams
		result providers.ExecContainerResult
		err    error
	}

	type want struct {
		called bool
		result ExecContainerResult
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"valid input": {
			given: given{
				params: NewExecContainerParams().SetName("web").SetCmd([]string{"echo", "hello"}),
				result: providers.ExecContainerResult{ExecID: "exec-123", ExitCode: 0, Stdout: "hello"},
			},
			want: want{
				called: true,
				result: ExecContainerResult{ExecID: "exec-123", ExitCode: 0, Stdout: "hello"},
			},
		},
		"trimmed name": {
			given: given{
				params: NewExecContainerParams().SetName("  web  ").SetCmd([]string{"echo"}),
				result: providers.ExecContainerResult{ExecID: "exec-123"},
			},
			want: want{called: true, result: ExecContainerResult{ExecID: "exec-123"}},
		},
		"empty name": {
			given: given{params: NewExecContainerParams().SetName("").SetCmd([]string{"echo"})},
			want:  want{},
		},
		"whitespace name": {
			given: given{params: NewExecContainerParams().SetName("   ").SetCmd([]string{"echo"})},
			want:  want{},
		},
		"empty command": {
			given: given{params: NewExecContainerParams().SetName("web").SetCmd([]string{})},
			want:  want{},
		},
		"provider error": {
			given: given{
				params: NewExecContainerParams().SetName("web").SetCmd([]string{"echo"}),
				err:    errors.New("docker error"),
			},
			want: want{called: true},
		},
		"stdin with attach_stdin": {
			given: given{
				params: NewExecContainerParams().SetName("web").SetCmd([]string{"cat"}).SetAttachStdin(true).SetStdin("hello"),
				result: providers.ExecContainerResult{ExecID: "exec-123", ExitCode: 0, Stdout: "hello"},
			},
			want: want{
				called: true,
				result: ExecContainerResult{ExecID: "exec-123", ExitCode: 0, Stdout: "hello"},
			},
		},
		"attach stdin explicit without stdin": {
			given: given{
				params: NewExecContainerParams().SetName("web").SetCmd([]string{"cat"}).SetAttachStdin(true),
				result: providers.ExecContainerResult{ExecID: "exec-123", ExitCode: 0},
			},
			want: want{called: true, result: ExecContainerResult{ExecID: "exec-123"}},
		},
		"stdin without attach_stdin fails validation": {
			given: given{
				params: NewExecContainerParams().SetName("web").SetCmd([]string{"cat"}).SetStdin("hello"),
			},
			want: want{},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := containermock.NewMockProviderClient(t)

			if tt.want.called {
				mockClient.On("ExecContainer", mock.Anything, providers.ExecContainerParams{
					Name:         tt.given.params.Name,
					Cmd:          tt.given.params.Cmd,
					Env:          tt.given.params.Env,
					WorkingDir:   tt.given.params.WorkingDir,
					User:         tt.given.params.User,
					Privileged:   tt.given.params.Privileged,
					TTY:          tt.given.params.TTY,
					AttachStdin:  tt.given.params.AttachStdin,
					AttachStdout: tt.given.params.AttachStdout,
					AttachStderr: tt.given.params.AttachStderr,
					Stdin:        tt.given.params.Stdin,
				}).Return(tt.given.result, tt.given.err)
			}

			service := NewService(mockClient, services.Policy{}, zap.NewNop())

			result, err := service.ExecContainer(context.Background(), tt.given.params)

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

func TestServiceExecContainer_ReadOnly(t *testing.T) {
	mockClient := containermock.NewMockProviderClient(t)
	policy := services.NewPolicy(true)
	service := NewService(mockClient, policy, zap.NewNop())

	_, err := service.ExecContainer(context.Background(), NewExecContainerParams().SetName("web").SetCmd([]string{"echo"}))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "read-only")
}
