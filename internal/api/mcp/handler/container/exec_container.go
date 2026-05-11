package container

import (
	"context"
	"fmt"

	"github.com/irwinby/container-runtime-mcp/internal/service/container"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ExecContainerInput struct {
	Name         string   `json:"name" jsonschema:"container name or ID"`
	Command      []string `json:"command" jsonschema:"command and arguments to execute"`
	Env          []string `json:"env,omitempty" jsonschema:"environment variables in KEY=VALUE format"`
	WorkingDir   string   `json:"working_dir,omitempty" jsonschema:"working directory for the command"`
	User         string   `json:"user,omitempty" jsonschema:"user to run the command as"`
	Privileged   bool     `json:"privileged,omitempty" jsonschema:"run in privileged mode"`
	TTY          bool     `json:"tty,omitempty" jsonschema:"allocate a pseudo-TTY"`
	AttachStdin  *bool    `json:"attach_stdin,omitempty" jsonschema:"attach standard input"`
	AttachStdout *bool    `json:"attach_stdout,omitempty" jsonschema:"attach standard output"`
	AttachStderr *bool    `json:"attach_stderr,omitempty" jsonschema:"attach standard error"`
	Stdin        string   `json:"stdin,omitempty" jsonschema:"data to write to standard input"`
}

func (i *ExecContainerInput) Validate() error {
	if i == nil {
		return fmt.Errorf("exec container input is nil")
	}

	if i.Name == "" {
		return fmt.Errorf("name is required")
	}

	if len(i.Command) == 0 {
		return fmt.Errorf("command is required")
	}

	return nil
}

type ExecContainerOutput struct {
	ExecID   string `json:"exec_id"`
	ExitCode int    `json:"exit_code"`
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
}

func (h *ToolsHandler) ExecContainer(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input ExecContainerInput,
) (*mcp.CallToolResult, ExecContainerOutput, error) {
	err := input.Validate()
	if err != nil {
		return nil, ExecContainerOutput{}, fmt.Errorf("validate exec container input: %w", err)
	}

	params := container.NewExecContainerParams().
		SetName(input.Name).
		SetCmd(input.Command).
		SetEnv(input.Env).
		SetWorkingDir(input.WorkingDir).
		SetUser(input.User).
		SetPrivileged(input.Privileged).
		SetTTY(input.TTY).
		SetStdin(input.Stdin)

	if input.AttachStdin != nil {
		params = params.SetAttachStdin(*input.AttachStdin)
	} else if input.Stdin != "" {
		params = params.SetAttachStdin(true)
	}

	if input.AttachStdout != nil {
		params = params.SetAttachStdout(*input.AttachStdout)
	}

	if input.AttachStderr != nil {
		params = params.SetAttachStderr(*input.AttachStderr)
	}

	result, err := h.containerService.ExecContainer(ctx, params)
	if err != nil {
		return nil, ExecContainerOutput{}, fmt.Errorf("exec container: %w", err)
	}

	return nil, ExecContainerOutput{
		ExecID:   result.ExecID,
		ExitCode: result.ExitCode,
		Stdout:   result.Stdout,
		Stderr:   result.Stderr,
	}, nil
}
