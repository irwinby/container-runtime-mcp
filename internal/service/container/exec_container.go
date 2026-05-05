package container

import (
	"context"
	"fmt"
	"strings"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"go.uber.org/zap"
)

// ExecContainerParams holds the parameters for executing a command in a container.
type ExecContainerParams struct {
	Name         string
	Cmd          []string
	Env          []string
	WorkingDir   string
	User         string
	Privileged   bool
	TTY          bool
	AttachStdin  bool
	AttachStdout bool
	AttachStderr bool
	Stdin        string
}

func NewExecContainerParams() ExecContainerParams {
	return ExecContainerParams{
		AttachStdout: true,
		AttachStderr: true,
	}
}

func (p ExecContainerParams) SetName(name string) ExecContainerParams {
	p.Name = strings.TrimSpace(name)
	return p
}

func (p ExecContainerParams) SetCmd(cmd []string) ExecContainerParams {
	p.Cmd = cmd
	return p
}

func (p ExecContainerParams) SetEnv(env []string) ExecContainerParams {
	p.Env = env
	return p
}

func (p ExecContainerParams) SetWorkingDir(workingDir string) ExecContainerParams {
	p.WorkingDir = strings.TrimSpace(workingDir)
	return p
}

func (p ExecContainerParams) SetUser(user string) ExecContainerParams {
	p.User = strings.TrimSpace(user)
	return p
}

func (p ExecContainerParams) SetPrivileged(privileged bool) ExecContainerParams {
	p.Privileged = privileged
	return p
}

func (p ExecContainerParams) SetTTY(tty bool) ExecContainerParams {
	p.TTY = tty
	return p
}

func (p ExecContainerParams) SetAttachStdin(attachStdin bool) ExecContainerParams {
	p.AttachStdin = attachStdin
	return p
}

func (p ExecContainerParams) SetStdin(stdin string) ExecContainerParams {
	p.Stdin = stdin
	return p
}

func (p ExecContainerParams) SetAttachStdout(attachStdout bool) ExecContainerParams {
	p.AttachStdout = attachStdout
	return p
}

func (p ExecContainerParams) SetAttachStderr(attachStderr bool) ExecContainerParams {
	p.AttachStderr = attachStderr
	return p
}

// Validate checks that required fields are present and trims whitespace.
func (p ExecContainerParams) Validate() (ExecContainerParams, error) {
	if p.Name == "" {
		return ExecContainerParams{}, fmt.Errorf("name is required")
	}

	if len(p.Cmd) == 0 {
		return ExecContainerParams{}, fmt.Errorf("command is required")
	}

	if p.Stdin != "" && !p.AttachStdin {
		return ExecContainerParams{}, fmt.Errorf("attach_stdin must be true when stdin is provided")
	}

	return p, nil
}

func (s *Service) ExecContainer(ctx context.Context, params ExecContainerParams) (ExecContainerResult, error) {
	err := s.policy.IsWriteAllowed()
	if err != nil {
		s.logger.Warn("exec container blocked by policy", zap.Error(err))
		return ExecContainerResult{}, fmt.Errorf("check if write is allowed: %w", err)
	}

	params, err = params.Validate()
	if err != nil {
		s.logger.Warn("exec container validation failed", zap.Error(err))
		return ExecContainerResult{}, fmt.Errorf("validate exec container params: %w", err)
	}

	s.logger.Info("executing command in container",
		zap.String("name", params.Name),
		zap.Int("cmd_count", len(params.Cmd)),
	)

	result, err := s.providerClient.ExecContainer(ctx, providers.ExecContainerParams{
		Name:         params.Name,
		Cmd:          params.Cmd,
		Env:          params.Env,
		WorkingDir:   params.WorkingDir,
		User:         params.User,
		Privileged:   params.Privileged,
		TTY:          params.TTY,
		AttachStdin:  params.AttachStdin,
		AttachStdout: params.AttachStdout,
		AttachStderr: params.AttachStderr,
		Stdin:        params.Stdin,
	})
	if err != nil {
		s.logger.Error("exec container failed", zap.String("name", params.Name), zap.Error(err))
		return ExecContainerResult{}, fmt.Errorf("exec container: %w", err)
	}

	s.logger.Info("command executed in container",
		zap.String("name", params.Name),
		zap.Int("exit_code", result.ExitCode),
	)

	return ExecContainerResult{
		ExecID:   result.ExecID,
		ExitCode: result.ExitCode,
		Stdout:   result.Stdout,
		Stderr:   result.Stderr,
	}, nil
}
