package container

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetters(t *testing.T) {
	t.Run("create container params", func(t *testing.T) {
		params := NewCreateContainerParams().
			SetName("  web  ").
			SetImage("  nginx  ")

		assert.Equal(t, "web", params.Name)
		assert.Equal(t, "nginx", params.Image)
	})

	t.Run("remove container params", func(t *testing.T) {
		params := NewRemoveContainerParams().
			SetName("  web  ").
			SetForce(true).
			SetRemoveVolumes(true).
			SetRemoveLinks(true)

		assert.Equal(t, "web", params.Name)
		assert.True(t, params.Force)
		assert.True(t, params.RemoveVolumes)
		assert.True(t, params.RemoveLinks)
	})

	t.Run("list containers params", func(t *testing.T) {
		params := NewListContainersParams().
			SetAll(true).
			SetLimit(10).
			SetSize(true).
			SetLatest(true)

		assert.True(t, params.All)
		assert.Equal(t, 10, params.Limit)
		assert.True(t, params.Size)
		assert.True(t, params.Latest)
	})

	t.Run("inspect container params", func(t *testing.T) {
		params := NewInspectContainerParams().SetName("  web  ")
		assert.Equal(t, "web", params.Name)
	})

	t.Run("start container params", func(t *testing.T) {
		params := NewStartContainerParams().SetName("  web  ")
		assert.Equal(t, "web", params.Name)
	})

	t.Run("stop container params", func(t *testing.T) {
		timeout := 30
		params := NewStopContainerParams().
			SetName("  web  ").
			SetSignal("SIGTERM").
			SetTimeoutSeconds(&timeout)

		assert.Equal(t, "web", params.Name)
		assert.Equal(t, "SIGTERM", params.Signal)
		assert.Equal(t, &timeout, params.TimeoutSeconds)
	})

	t.Run("restart container params", func(t *testing.T) {
		timeout := 30
		params := NewRestartContainerParams().
			SetName("  web  ").
			SetSignal("SIGTERM").
			SetTimeoutSeconds(&timeout)

		assert.Equal(t, "web", params.Name)
		assert.Equal(t, "SIGTERM", params.Signal)
		assert.Equal(t, &timeout, params.TimeoutSeconds)
	})

	t.Run("container logs params", func(t *testing.T) {
		params := NewContainerLogsParams().
			SetName("  web  ").
			SetStdout(false).
			SetStderr(false).
			SetSince("1h").
			SetTimestamps(true).
			SetTail("100")

		assert.Equal(t, "web", params.Name)
		assert.False(t, params.Stdout)
		assert.False(t, params.Stderr)
		assert.Equal(t, "1h", params.Since)
		assert.True(t, params.Timestamps)
		assert.Equal(t, "100", params.Tail)
	})

	t.Run("exec container params", func(t *testing.T) {
		params := NewExecContainerParams().
			SetName("  web  ").
			SetCmd([]string{"echo"}).
			SetEnv([]string{"FOO=bar"}).
			SetWorkingDir("/tmp").
			SetUser("root").
			SetPrivileged(true).
			SetTTY(true).
			SetAttachStdin(true).
			SetStdin("hello").
			SetAttachStdout(false).
			SetAttachStderr(false)

		assert.Equal(t, "web", params.Name)
		assert.Equal(t, []string{"echo"}, params.Cmd)
		assert.Equal(t, []string{"FOO=bar"}, params.Env)
		assert.Equal(t, "/tmp", params.WorkingDir)
		assert.Equal(t, "root", params.User)
		assert.True(t, params.Privileged)
		assert.True(t, params.TTY)
		assert.True(t, params.AttachStdin)
		assert.Equal(t, "hello", params.Stdin)
		assert.False(t, params.AttachStdout)
		assert.False(t, params.AttachStderr)
	})
}

func TestServiceCanWrite(t *testing.T) {
	t.Run("read only", func(t *testing.T) {
		service := &Service{policy: struct{ ReadOnly bool }{ReadOnly: true}}
		assert.False(t, service.CanWrite())
	})

	t.Run("writable", func(t *testing.T) {
		service := &Service{policy: struct{ ReadOnly bool }{ReadOnly: false}}
		assert.True(t, service.CanWrite())
	})
}
