package container

import (
	"context"
	"errors"
	"testing"

	containermock "github.com/irwinby/container-runtime-mcp/internal/api/mcp/handler/container/mock"
	"github.com/irwinby/container-runtime-mcp/internal/service/container"
	"github.com/irwinby/container-runtime-mcp/pkg/ptr"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandlerExecContainer(t *testing.T) {
	type given struct {
		input  ExecContainerInput
		result container.ExecContainerResult
		err    error
	}

	type want struct {
		called bool
		name   string
		cmd    []string
		params container.ExecContainerParams
		exec   ExecContainerOutput
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"valid input": {
			given: given{
				input:  ExecContainerInput{Name: "web", Command: []string{"echo", "hello"}},
				result: container.ExecContainerResult{ExecID: "exec-123", ExitCode: 0, Stdout: "hello"},
			},
			want: want{
				called: true,
				name:   "web",
				cmd:    []string{"echo", "hello"},
				params: container.ExecContainerParams{
					Name:         "web",
					Cmd:          []string{"echo", "hello"},
					AttachStdout: true,
					AttachStderr: true,
				},
				exec: ExecContainerOutput{ExecID: "exec-123", ExitCode: 0, Stdout: "hello"},
			},
		},
		"empty name": {
			given: given{
				input: ExecContainerInput{Name: ""},
				err:   errors.New("validation error"),
			},
			want: want{},
		},
		"empty command": {
			given: given{
				input: ExecContainerInput{Name: "web"},
				err:   errors.New("validation error"),
			},
			want: want{},
		},
		"service error": {
			given: given{
				input: ExecContainerInput{Name: "web", Command: []string{"echo"}},
				err:   errors.New("docker error"),
			},
			want: want{
				called: true,
				name:   "web",
				cmd:    []string{"echo"},
				params: container.ExecContainerParams{
					Name:         "web",
					Cmd:          []string{"echo"},
					AttachStdout: true,
					AttachStderr: true,
				},
			},
		},
		"stdin auto-enables attach_stdin": {
			given: given{
				input:  ExecContainerInput{Name: "web", Command: []string{"cat"}, Stdin: "hello"},
				result: container.ExecContainerResult{ExecID: "exec-123", ExitCode: 0, Stdout: "hello"},
			},
			want: want{
				called: true,
				name:   "web",
				cmd:    []string{"cat"},
				params: container.ExecContainerParams{
					Name:         "web",
					Cmd:          []string{"cat"},
					AttachStdin:  true,
					AttachStdout: true,
					AttachStderr: true,
					Stdin:        "hello",
				},
				exec: ExecContainerOutput{ExecID: "exec-123", ExitCode: 0, Stdout: "hello"},
			},
		},
		"explicit attach_stdout false": {
			given: given{
				input:  ExecContainerInput{Name: "web", Command: []string{"echo"}, AttachStdout: ptr.Bool(false)},
				result: container.ExecContainerResult{ExecID: "exec-123", ExitCode: 0, Stderr: "warn"},
			},
			want: want{
				called: true,
				name:   "web",
				cmd:    []string{"echo"},
				params: container.ExecContainerParams{
					Name:         "web",
					Cmd:          []string{"echo"},
					AttachStdout: false,
					AttachStderr: true,
				},
				exec: ExecContainerOutput{ExecID: "exec-123", ExitCode: 0, Stderr: "warn"},
			},
		},
		"explicit attach_stdin and attach_stderr": {
			given: given{
				input:  ExecContainerInput{Name: "web", Command: []string{"echo"}, AttachStdin: ptr.Bool(true), AttachStderr: ptr.Bool(false)},
				result: container.ExecContainerResult{ExecID: "exec-123", ExitCode: 0},
			},
			want: want{
				called: true,
				name:   "web",
				cmd:    []string{"echo"},
				params: container.ExecContainerParams{
					Name:         "web",
					Cmd:          []string{"echo"},
					AttachStdin:  true,
					AttachStdout: true,
					AttachStderr: false,
				},
				exec: ExecContainerOutput{ExecID: "exec-123", ExitCode: 0},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockService := containermock.NewMockContainerService(t)

			if test.want.called {
				mockService.On("ExecContainer", mock.Anything, test.want.params).Return(test.given.result, test.given.err)
			}

			handler := NewToolsHandler(mockService)

			_, output, err := handler.ExecContainer(context.Background(), &mcp.CallToolRequest{}, test.given.input)

			if test.given.err != nil {
				require.Error(t, err)
				return
			}

			if !test.want.called {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.want.exec, output)
		})
	}
}
