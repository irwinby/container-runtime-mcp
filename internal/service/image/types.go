package image

// Image represents a summarized image for listing.
type Image struct {
	ID         string   `json:"id"`
	RepoTags   []string `json:"repo_tags"`
	Size       int64    `json:"size"`
	Created    int64    `json:"created"`
	Containers int64    `json:"containers"`
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
