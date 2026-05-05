package container

// Container represents a summarized container for listing.
type Container struct {
	ID      string   `json:"id"`
	Names   []string `json:"names"`
	Image   string   `json:"image"`
	State   string   `json:"state"`
	Status  string   `json:"status"`
	Created int64    `json:"created"`
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

// ContainerLogsResult holds the container log output.
type ContainerLogsResult struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

// ExecContainerResult holds the result of executing a command in a container.
type ExecContainerResult struct {
	ExecID   string `json:"exec_id"`
	ExitCode int    `json:"exit_code"`
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
}
