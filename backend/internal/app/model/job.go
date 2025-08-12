package model

import (
	"axolotl-cloud/infra/logger"
	"context"
)

type Job struct {
	ID          uint                                                      `gorm:"primaryKey" json:"id"`
	Name        string                                                    `json:"name"`
	Run         func(ctx context.Context, jobLogger *logger.Logger) error `gorm:"-" json:"-"`
	Logs        []JobLog                                                  `gorm:"foreignKey:JobID;references:ID;constraint:OnDelete:CASCADE" json:"logs"`
	Status      JobStatus                                                 `json:"status"`
	CreatedAt   int64                                                     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   int64                                                     `json:"updated_at" gorm:"autoUpdateTime"`
	ContainerID *uint                                                     `json:"container_id" gorm:"default:null"`
}

type JobLog struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	JobID     uint   `json:"-"`
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
