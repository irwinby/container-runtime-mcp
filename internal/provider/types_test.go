package providers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContainerLogsResult(t *testing.T) {
	result := NewContainerLogsResult()
	assert.Equal(t, ContainerLogsResult{}, result)

	result = result.SetStdout("stdout").SetStderr("stderr")
	assert.Equal(t, "stdout", result.Stdout)
	assert.Equal(t, "stderr", result.Stderr)
}

func TestContainerBuilder(t *testing.T) {
	container := NewContainer().
		SetID("id-1").
		SetNames([]string{"name1"}).
		SetImage("nginx").
		SetState("running").
		SetStatus("Up").
		SetCreated(123)

	assert.Equal(t, "id-1", container.ID)
	assert.Equal(t, []string{"name1"}, container.Names)
	assert.Equal(t, "nginx", container.Image)
	assert.Equal(t, "running", container.State)
	assert.Equal(t, "Up", container.Status)
	assert.Equal(t, int64(123), container.Created)
}

func TestContainerInspectBuilder(t *testing.T) {
	inspect := NewContainerInspect().
		SetID("id-1").
		SetName("name1").
		SetImage("nginx").
		SetState("running").
		SetStatus("Up").
		SetCreated("2024-01-01").
		SetPath("/bin/sh").
		SetArgs([]string{"-c"}).
		SetRestartCount(1)

	assert.Equal(t, "id-1", inspect.ID)
	assert.Equal(t, "name1", inspect.Name)
	assert.Equal(t, "nginx", inspect.Image)
	assert.Equal(t, "running", inspect.State)
	assert.Equal(t, "Up", inspect.Status)
	assert.Equal(t, "2024-01-01", inspect.Created)
	assert.Equal(t, "/bin/sh", inspect.Path)
	assert.Equal(t, []string{"-c"}, inspect.Args)
	assert.Equal(t, 1, inspect.RestartCount)
}

func TestImageBuilder(t *testing.T) {
	image := NewImage().
		SetID("img-1").
		SetRepoTags([]string{"nginx:latest"}).
		SetSize(1000).
		SetCreated(123).
		SetContainers(5)

	assert.Equal(t, "img-1", image.ID)
	assert.Equal(t, []string{"nginx:latest"}, image.RepoTags)
	assert.Equal(t, int64(1000), image.Size)
	assert.Equal(t, int64(123), image.Created)
	assert.Equal(t, int64(5), image.Containers)
}

func TestImageInspectBuilder(t *testing.T) {
	inspect := NewImageInspect().
		SetID("img-1").
		SetRepoTags([]string{"nginx:latest"}).
		SetSize(1000).
		SetCreated("2024-01-01").
		SetArchitecture("amd64").
		SetOS("linux")

	assert.Equal(t, "img-1", inspect.ID)
	assert.Equal(t, []string{"nginx:latest"}, inspect.RepoTags)
	assert.Equal(t, int64(1000), inspect.Size)
	assert.Equal(t, "2024-01-01", inspect.Created)
	assert.Equal(t, "amd64", inspect.Architecture)
	assert.Equal(t, "linux", inspect.OS)
}

func TestSystemInfoBuilder(t *testing.T) {
	info := NewSystemInfo().
		SetID("sys-1").
		SetContainers(10).
		SetContainersRunning(5).
		SetContainersPaused(1).
		SetContainersStopped(4).
		SetImages(20).
		SetDriver("overlay2").
		SetDriverStatus([][2]string{{"key", "val"}}).
		SetKernelVersion("5.4").
		SetOperatingSystem("Ubuntu").
		SetOSType("linux").
		SetArchitecture("amd64").
		SetNCPU(8).
		SetMemTotal(16000).
		SetServerVersion("24.0")

	assert.Equal(t, "sys-1", info.ID)
	assert.Equal(t, 10, info.Containers)
	assert.Equal(t, 5, info.ContainersRunning)
	assert.Equal(t, 1, info.ContainersPaused)
	assert.Equal(t, 4, info.ContainersStopped)
	assert.Equal(t, 20, info.Images)
	assert.Equal(t, "overlay2", info.Driver)
	assert.Equal(t, [][2]string{{"key", "val"}}, info.DriverStatus)
	assert.Equal(t, "5.4", info.KernelVersion)
	assert.Equal(t, "Ubuntu", info.OperatingSystem)
	assert.Equal(t, "linux", info.OSType)
	assert.Equal(t, "amd64", info.Architecture)
	assert.Equal(t, 8, info.NCPU)
	assert.Equal(t, int64(16000), info.MemTotal)
	assert.Equal(t, "24.0", info.ServerVersion)
}

func TestSystemVersionBuilder(t *testing.T) {
	version := NewSystemVersion().
		SetVersion("24.0").
		SetAPIVersion("1.43").
		SetMinAPIVersion("1.12").
		SetOs("linux").
		SetArch("amd64").
		SetPlatformName("Docker")

	assert.Equal(t, "24.0", version.Version)
	assert.Equal(t, "1.43", version.APIVersion)
	assert.Equal(t, "1.12", version.MinAPIVersion)
	assert.Equal(t, "linux", version.Os)
	assert.Equal(t, "amd64", version.Arch)
	assert.Equal(t, "Docker", version.PlatformName)
}

func TestPingResultBuilder(t *testing.T) {
	result := NewPingResult().
		SetAPIVersion("1.43").
		SetOSType("linux").
		SetExperimental(true).
		SetBuilderVersion("2")

	assert.Equal(t, "1.43", result.APIVersion)
	assert.Equal(t, "linux", result.OSType)
	assert.True(t, result.Experimental)
	assert.Equal(t, "2", result.BuilderVersion)
}

func TestVolumeBuilder(t *testing.T) {
	volume := NewVolume().
		SetName("vol-1").
		SetDriver("local").
		SetMountpoint("/mnt").
		SetLabels(map[string]string{"a": "b"}).
		SetScope("local")

	assert.Equal(t, "vol-1", volume.Name)
	assert.Equal(t, "local", volume.Driver)
	assert.Equal(t, "/mnt", volume.Mountpoint)
	assert.Equal(t, map[string]string{"a": "b"}, volume.Labels)
	assert.Equal(t, "local", volume.Scope)
}

func TestVolumeInspectBuilder(t *testing.T) {
	inspect := NewVolumeInspect().
		SetName("vol-1").
		SetDriver("local").
		SetMountpoint("/mnt").
		SetLabels(map[string]string{"a": "b"}).
		SetScope("local")

	assert.Equal(t, "vol-1", inspect.Name)
	assert.Equal(t, "local", inspect.Driver)
	assert.Equal(t, "/mnt", inspect.Mountpoint)
	assert.Equal(t, map[string]string{"a": "b"}, inspect.Labels)
	assert.Equal(t, "local", inspect.Scope)
}
