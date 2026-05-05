package container

import (
	"context"
	"errors"
	"testing"

	containermock "github.com/irwinby/container-runtime-mcp/internal/api/mcp/tools/container/mock"
	"github.com/irwinby/container-runtime-mcp/internal/service/container"
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
				input:  ExecContainerInput{Name: "web", Command: []string{"echo"}, AttachStdout: boolPtr(false)},
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
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockSvc := containermock.NewMockContainerService(t)
			if tt.want.called {
				mockSvc.On("ExecContainer", mock.Anything, tt.want.params).Return(tt.given.result, tt.given.err)
			}

			handler := NewToolsHandler(mockSvc)

			_, output, err := handler.ExecContainer(context.Background(), &mcp.CallToolRequest{}, tt.given.input)

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			if !tt.want.called {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want.exec, output)
		})
	}
}

func boolPtr(v bool) *bool {
	return &v
}
