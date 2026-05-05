package providers

import ocispec "github.com/opencontainers/image-spec/specs-go/v1"

// CreateContainerParams holds the parameters for creating a container.
type CreateContainerParams struct {
	Name  string
	Image string
}

// RemoveContainerParams holds the parameters for removing a container.
type RemoveContainerParams struct {
	Name          string
	Force         bool
	RemoveVolumes bool
	RemoveLinks   bool
}

// ListContainersParams holds the parameters for listing containers.
type ListContainersParams struct {
	All    bool
	Limit  int
	Size   bool
	Latest bool
}

// InspectContainerParams holds the parameters for inspecting a container.
type InspectContainerParams struct {
	Name string
}

// ContainerLogsParams holds the parameters for retrieving container logs.
type ContainerLogsParams struct {
	Name       string
	Stdout     bool
	Stderr     bool
	Since      string
	Timestamps bool
	Tail       string
}

// ContainerLogsResult holds the container log output.
type ContainerLogsResult struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

// NewContainerLogsResult creates a new empty ContainerLogsResult.
func NewContainerLogsResult() ContainerLogsResult {
	return ContainerLogsResult{}
}

// SetStdout sets the stdout content.
func (r ContainerLogsResult) SetStdout(stdout string) ContainerLogsResult {
	r.Stdout = stdout
	return r
}

// SetStderr sets the stderr content.
func (r ContainerLogsResult) SetStderr(stderr string) ContainerLogsResult {
	r.Stderr = stderr
	return r
}

// StartContainerParams holds the parameters for starting a container.
type StartContainerParams struct {
	Name string
}

// StopContainerParams holds the parameters for stopping a container.
type StopContainerParams struct {
	Name           string
	Signal         string
	TimeoutSeconds *int
}

// RestartContainerParams holds the parameters for restarting a container.
type RestartContainerParams struct {
	Name           string
	Signal         string
	TimeoutSeconds *int
}

// Container represents a summarized container for listing.
type Container struct {
	ID      string   `json:"id"`
	Names   []string `json:"names"`
	Image   string   `json:"image"`
	State   string   `json:"state"`
	Status  string   `json:"status"`
	Created int64    `json:"created"`
}

// NewContainer creates a new empty Container.
func NewContainer() Container {
	return Container{}
}

// SetID sets the container ID.
func (c Container) SetID(id string) Container {
	c.ID = id
	return c
}

// SetNames sets the container names.
func (c Container) SetNames(names []string) Container {
	c.Names = names
	return c
}

// SetImage sets the container image.
func (c Container) SetImage(image string) Container {
	c.Image = image
	return c
}

// SetState sets the container state.
func (c Container) SetState(state string) Container {
	c.State = state
	return c
}

// SetStatus sets the container status.
func (c Container) SetStatus(status string) Container {
	c.Status = status
	return c
}

// SetCreated sets the container creation timestamp.
func (c Container) SetCreated(created int64) Container {
	c.Created = created
	return c
}

// ContainerInspect represents a detailed container for inspection.
type ContainerInspect struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	State        string   `json:"state"`
	Status       string   `json:"status"`
	Created      string   `json:"created"`
	Path         string   `json:"path"`
	Args         []string `json:"args"`
	RestartCount int      `json:"restart_count"`
}

// NewContainerInspect creates a new empty ContainerInspect.
func NewContainerInspect() ContainerInspect {
	return ContainerInspect{}
}

// SetID sets the container ID.
func (c ContainerInspect) SetID(id string) ContainerInspect {
	c.ID = id
	return c
}

// SetName sets the container name.
func (c ContainerInspect) SetName(name string) ContainerInspect {
	c.Name = name
	return c
}

// SetImage sets the container image.
func (c ContainerInspect) SetImage(image string) ContainerInspect {
	c.Image = image
	return c
}

// SetState sets the container state.
func (c ContainerInspect) SetState(state string) ContainerInspect {
	c.State = state
	return c
}

// SetStatus sets the container status.
func (c ContainerInspect) SetStatus(status string) ContainerInspect {
	c.Status = status
	return c
}

// SetCreated sets the container creation time.
func (c ContainerInspect) SetCreated(created string) ContainerInspect {
	c.Created = created
	return c
}

// SetPath sets the container path.
func (c ContainerInspect) SetPath(path string) ContainerInspect {
	c.Path = path
	return c
}

// SetArgs sets the container arguments.
func (c ContainerInspect) SetArgs(args []string) ContainerInspect {
	c.Args = args
	return c
}

// SetRestartCount sets the container restart count.
func (c ContainerInspect) SetRestartCount(restartCount int) ContainerInspect {
	c.RestartCount = restartCount
	return c
}

// PullImageParams holds the parameters for pulling an image.
type PullImageParams struct {
	Ref      string
	All      bool
	Platform *ocispec.Platform
}

// PushImageParams holds the parameters for pushing an image.
type PushImageParams struct {
	Ref      string
	All      bool
	Platform *ocispec.Platform
}

// ListImagesParams holds the parameters for listing images.
type ListImagesParams struct {
	All        bool
	SharedSize bool
}

// InspectImageParams holds the parameters for inspecting an image.
type InspectImageParams struct {
	Ref string
}

// RemoveImageParams holds the parameters for removing an image.
type RemoveImageParams struct {
	Ref           string
	Force         bool
	PruneChildren bool
	Platform      *ocispec.Platform
}

// TagImageParams holds the parameters for tagging an image.
type TagImageParams struct {
	Source string
	Target string
}

// Image represents a summarized image for listing.
type Image struct {
	ID         string   `json:"id"`
	RepoTags   []string `json:"repo_tags"`
	Size       int64    `json:"size"`
	Created    int64    `json:"created"`
	Containers int64    `json:"containers"`
}

// NewImage creates a new empty Image.
func NewImage() Image {
	return Image{}
}

// SetID sets the image ID.
func (i Image) SetID(id string) Image {
	i.ID = id
	return i
}

// SetRepoTags sets the image repository tags.
func (i Image) SetRepoTags(repoTags []string) Image {
	i.RepoTags = repoTags
	return i
}

// SetSize sets the image size.
func (i Image) SetSize(size int64) Image {
	i.Size = size
	return i
}

// SetCreated sets the image creation timestamp.
func (i Image) SetCreated(created int64) Image {
	i.Created = created
	return i
}

// SetContainers sets the number of containers using this image.
func (i Image) SetContainers(containers int64) Image {
	i.Containers = containers
	return i
}

// ImageInspect represents a detailed image for inspection.
type ImageInspect struct {
	ID           string   `json:"id"`
	RepoTags     []string `json:"repo_tags"`
	Size         int64    `json:"size"`
	Created      string   `json:"created"`
	Architecture string   `json:"architecture"`
	OS           string   `json:"os"`
}

// NewImageInspect creates a new empty ImageInspect.
func NewImageInspect() ImageInspect {
	return ImageInspect{}
}

// SetID sets the image ID.
func (i ImageInspect) SetID(id string) ImageInspect {
	i.ID = id
	return i
}

// SetRepoTags sets the image repository tags.
func (i ImageInspect) SetRepoTags(repoTags []string) ImageInspect {
	i.RepoTags = repoTags
	return i
}

// SetSize sets the image size.
func (i ImageInspect) SetSize(size int64) ImageInspect {
	i.Size = size
	return i
}

// SetCreated sets the image creation time.
func (i ImageInspect) SetCreated(created string) ImageInspect {
	i.Created = created
	return i
}

// SetArchitecture sets the image architecture.
func (i ImageInspect) SetArchitecture(architecture string) ImageInspect {
	i.Architecture = architecture
	return i
}

// SetOS sets the image operating system.
func (i ImageInspect) SetOS(os string) ImageInspect {
	i.OS = os
	return i
}

// SystemInfo represents Docker system information.
type SystemInfo struct {
	ID                string      `json:"id"`
	Containers        int         `json:"containers"`
	ContainersRunning int         `json:"containers_running"`
	ContainersPaused  int         `json:"containers_paused"`
	ContainersStopped int         `json:"containers_stopped"`
	Images            int         `json:"images"`
	Driver            string      `json:"driver"`
	DriverStatus      [][2]string `json:"driver_status"`
	KernelVersion     string      `json:"kernel_version"`
	OperatingSystem   string      `json:"operating_system"`
	OSType            string      `json:"os_type"`
	Architecture      string      `json:"architecture"`
	NCPU              int         `json:"ncpu"`
	MemTotal          int64       `json:"mem_total"`
	ServerVersion     string      `json:"server_version"`
}

// NewSystemInfo creates a new empty SystemInfo.
func NewSystemInfo() SystemInfo {
	return SystemInfo{}
}

// SetID sets the system ID.
func (s SystemInfo) SetID(id string) SystemInfo {
	s.ID = id
	return s
}

// SetContainers sets the total container count.
func (s SystemInfo) SetContainers(containers int) SystemInfo {
	s.Containers = containers
	return s
}

// SetContainersRunning sets the running container count.
func (s SystemInfo) SetContainersRunning(containersRunning int) SystemInfo {
	s.ContainersRunning = containersRunning
	return s
}

// SetContainersPaused sets the paused container count.
func (s SystemInfo) SetContainersPaused(containersPaused int) SystemInfo {
	s.ContainersPaused = containersPaused
	return s
}

// SetContainersStopped sets the stopped container count.
func (s SystemInfo) SetContainersStopped(containersStopped int) SystemInfo {
	s.ContainersStopped = containersStopped
	return s
}

// SetImages sets the image count.
func (s SystemInfo) SetImages(images int) SystemInfo {
	s.Images = images
	return s
}

// SetDriver sets the storage driver.
func (s SystemInfo) SetDriver(driver string) SystemInfo {
	s.Driver = driver
	return s
}

// SetDriverStatus sets the driver status.
func (s SystemInfo) SetDriverStatus(driverStatus [][2]string) SystemInfo {
	s.DriverStatus = driverStatus
	return s
}

// SetKernelVersion sets the kernel version.
func (s SystemInfo) SetKernelVersion(kernelVersion string) SystemInfo {
	s.KernelVersion = kernelVersion
	return s
}

// SetOperatingSystem sets the operating system.
func (s SystemInfo) SetOperatingSystem(operatingSystem string) SystemInfo {
	s.OperatingSystem = operatingSystem
	return s
}

// SetOSType sets the OS type.
func (s SystemInfo) SetOSType(osType string) SystemInfo {
	s.OSType = osType
	return s
}

// SetArchitecture sets the architecture.
func (s SystemInfo) SetArchitecture(architecture string) SystemInfo {
	s.Architecture = architecture
	return s
}

// SetNCPU sets the number of CPUs.
func (s SystemInfo) SetNCPU(ncpu int) SystemInfo {
	s.NCPU = ncpu
	return s
}

// SetMemTotal sets the total memory.
func (s SystemInfo) SetMemTotal(memTotal int64) SystemInfo {
	s.MemTotal = memTotal
	return s
}

// SetServerVersion sets the server version.
func (s SystemInfo) SetServerVersion(serverVersion string) SystemInfo {
	s.ServerVersion = serverVersion
	return s
}

// SystemVersion represents Docker version information.
type SystemVersion struct {
	Version       string `json:"version"`
	APIVersion    string `json:"api_version"`
	MinAPIVersion string `json:"min_api_version"`
	Os            string `json:"os"`
	Arch          string `json:"arch"`
	PlatformName  string `json:"platform_name"`
}

// NewSystemVersion creates a new empty SystemVersion.
func NewSystemVersion() SystemVersion {
	return SystemVersion{}
}

// SetVersion sets the version.
func (s SystemVersion) SetVersion(version string) SystemVersion {
	s.Version = version
	return s
}

// SetAPIVersion sets the API version.
func (s SystemVersion) SetAPIVersion(apiVersion string) SystemVersion {
	s.APIVersion = apiVersion
	return s
}

// SetMinAPIVersion sets the minimum API version.
func (s SystemVersion) SetMinAPIVersion(minAPIVersion string) SystemVersion {
	s.MinAPIVersion = minAPIVersion
	return s
}

// SetOs sets the OS.
func (s SystemVersion) SetOs(os string) SystemVersion {
	s.Os = os
	return s
}

// SetArch sets the architecture.
func (s SystemVersion) SetArch(arch string) SystemVersion {
	s.Arch = arch
	return s
}

// SetPlatformName sets the platform name.
func (s SystemVersion) SetPlatformName(platformName string) SystemVersion {
	s.PlatformName = platformName
	return s
}

// PingResult represents a Docker ping response.
type PingResult struct {
	APIVersion     string `json:"api_version"`
	OSType         string `json:"os_type"`
	Experimental   bool   `json:"experimental"`
	BuilderVersion string `json:"builder_version"`
}

// NewPingResult creates a new empty PingResult.
func NewPingResult() PingResult {
	return PingResult{}
}

// SetAPIVersion sets the API version.
func (p PingResult) SetAPIVersion(apiVersion string) PingResult {
	p.APIVersion = apiVersion
	return p
}

// SetOSType sets the OS type.
func (p PingResult) SetOSType(osType string) PingResult {
	p.OSType = osType
	return p
}

// SetExperimental sets the experimental flag.
func (p PingResult) SetExperimental(experimental bool) PingResult {
	p.Experimental = experimental
	return p
}

// SetBuilderVersion sets the builder version.
func (p PingResult) SetBuilderVersion(builderVersion string) PingResult {
	p.BuilderVersion = builderVersion
	return p
}

// ListVolumesParams holds the parameters for listing volumes.
type ListVolumesParams struct {
	Dangling bool
}

// InspectVolumeParams holds the parameters for inspecting a volume.
type InspectVolumeParams struct {
	Name string
}

// CreateVolumeParams holds the parameters for creating a volume.
type CreateVolumeParams struct {
	Name       string
	Driver     string
	DriverOpts map[string]string
	Labels     map[string]string
}

// RemoveVolumeParams holds the parameters for removing a volume.
type RemoveVolumeParams struct {
	Name  string
	Force bool
}

// Volume represents a summarized volume for listing.
type Volume struct {
	Name       string            `json:"name"`
	Driver     string            `json:"driver"`
	Mountpoint string            `json:"mountpoint"`
	Labels     map[string]string `json:"labels"`
	Scope      string            `json:"scope"`
}

// NewVolume creates a new empty Volume.
func NewVolume() Volume {
	return Volume{}
}

// SetName sets the volume name.
func (v Volume) SetName(name string) Volume {
	v.Name = name
	return v
}

// SetDriver sets the volume driver.
func (v Volume) SetDriver(driver string) Volume {
	v.Driver = driver
	return v
}

// SetMountpoint sets the volume mountpoint.
func (v Volume) SetMountpoint(mountpoint string) Volume {
	v.Mountpoint = mountpoint
	return v
}

// SetLabels sets the volume labels.
func (v Volume) SetLabels(labels map[string]string) Volume {
	v.Labels = labels
	return v
}

// SetScope sets the volume scope.
func (v Volume) SetScope(scope string) Volume {
	v.Scope = scope
	return v
}

// VolumeInspect represents a detailed volume for inspection.
type VolumeInspect struct {
	Name       string            `json:"name"`
	Driver     string            `json:"driver"`
	Mountpoint string            `json:"mountpoint"`
	Labels     map[string]string `json:"labels"`
	Scope      string            `json:"scope"`
}

// NewVolumeInspect creates a new empty VolumeInspect.
func NewVolumeInspect() VolumeInspect {
	return VolumeInspect{}
}

// SetName sets the volume name.
func (v VolumeInspect) SetName(name string) VolumeInspect {
	v.Name = name
	return v
}

// SetDriver sets the volume driver.
func (v VolumeInspect) SetDriver(driver string) VolumeInspect {
	v.Driver = driver
	return v
}

// SetMountpoint sets the volume mountpoint.
func (v VolumeInspect) SetMountpoint(mountpoint string) VolumeInspect {
	v.Mountpoint = mountpoint
	return v
}

// SetLabels sets the volume labels.
func (v VolumeInspect) SetLabels(labels map[string]string) VolumeInspect {
	v.Labels = labels
	return v
}

// SetScope sets the volume scope.
func (v VolumeInspect) SetScope(scope string) VolumeInspect {
	v.Scope = scope
	return v
}

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

// ExecContainerResult holds the result of executing a command in a container.
type ExecContainerResult struct {
	ExecID   string
	ExitCode int
	Stdout   string
	Stderr   string
}
