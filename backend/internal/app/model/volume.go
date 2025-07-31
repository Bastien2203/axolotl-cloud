package model

type Volume struct {
	Size        int64  `json:"size"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Type        string `json:"type"` // e.g., "local", "nfs", etc.
	ContainerID uint   `json:"container_id"`
	ProjectID   uint   `json:"project_id"`
}
