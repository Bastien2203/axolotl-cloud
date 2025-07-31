package model

import (
	"context"
)

type Job struct {
	ID          uint                                              `gorm:"primaryKey" json:"id"`
	Name        string                                            `json:"name"`
	Run         func(ctx context.Context, log func(string)) error `gorm:"-" json:"-"`
	Logs        []JobLog                                          `gorm:"foreignKey:JobID" json:"logs"`
	Status      JobStatus                                         `json:"status"`
	CreatedAt   int64                                             `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   int64                                             `json:"updated_at" gorm:"autoUpdateTime"`
	ContainerID uint                                              `json:"container_id"`
}

type JobLog struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	JobID     uint
	Line      string `json:"line"`
	CreatedAt int64  `json:"created_at" gorm:"autoCreateTime"`
}

type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
)
