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

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := containermock.NewMockProviderClient(t)

			if test.want.called {
				mockClient.On("ExecContainer", mock.Anything, providers.ExecContainerParams{
					Name:         test.given.params.Name,
					Cmd:          test.given.params.Cmd,
					Env:          test.given.params.Env,
					WorkingDir:   test.given.params.WorkingDir,
					User:         test.given.params.User,
					Privileged:   test.given.params.Privileged,
					TTY:          test.given.params.TTY,
					AttachStdin:  test.given.params.AttachStdin,
					AttachStdout: test.given.params.AttachStdout,
					AttachStderr: test.given.params.AttachStderr,
					Stdin:        test.given.params.Stdin,
				}).Return(test.given.result, test.given.err)
			}

			service := NewService(mockClient, services.Policy{}, zap.NewNop())

			result, err := service.ExecContainer(context.Background(), test.given.params)

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

func TestServiceExecContainer_ReadOnly(t *testing.T) {
	mockClient := containermock.NewMockProviderClient(t)
	policy := services.NewPolicy(true)
	service := NewService(mockClient, policy, zap.NewNop())

	_, err := service.ExecContainer(context.Background(), NewExecContainerParams().SetName("web").SetCmd([]string{"echo"}))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "read-only")
}
