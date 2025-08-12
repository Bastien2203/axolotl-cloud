package model

import "axolotl-cloud/types"

type Container struct {
	ID          uint             `gorm:"primaryKey" json:"id"`
	DockerImage string           `json:"docker_image" binding:"required"`
	Ports       types.StringMap  `gorm:"type:text" json:"ports"`
	Env         types.StringMap  `gorm:"type:text" json:"env"`
	Volumes     types.StringMap  `gorm:"type:text" json:"volumes"`
	Name        string           `json:"name" binding:"required"`
	ProjectID   uint             `json:"project_id"`
	Networks    types.StringList `gorm:"type:text" json:"networks"`
	NetworkMode string           `json:"network_mode" binding:"required,oneof=bridge host none" gorm:"default:bridge"`
	Jobs        []Job            `gorm:"foreignKey:ContainerID;constraint:OnDelete:CASCADE" json:"jobs"`
	LastJob     Job              `gorm:"foreignKey:ContainerID;constraint:OnDelete:SET NULL" json:"last_job,omitempty"`
}
