package system

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

// SystemVersion represents Docker version information.
type SystemVersion struct {
	Version       string `json:"version"`
	APIVersion    string `json:"api_version"`
	MinAPIVersion string `json:"min_api_version"`
	Os            string `json:"os"`
	Arch          string `json:"arch"`
	PlatformName  string `json:"platform_name"`
}

// PingResult represents a Docker ping response.
type PingResult struct {
	APIVersion     string `json:"api_version"`
	OSType         string `json:"os_type"`
	Experimental   bool   `json:"experimental"`
	BuilderVersion string `json:"builder_version"`
}
