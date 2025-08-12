package websocket

import "axolotl-cloud/internal/app/model"

type JobLogPayload struct {
	JobID uint         `json:"job_id"`
	Log   model.JobLog `json:"log"`
}

type ContainerStatusPayload struct {
	ContainerID string `json:"container_id"`
	Status      string `json:"status"`
}
