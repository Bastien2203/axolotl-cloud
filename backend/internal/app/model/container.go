package model

import "axolotl-cloud/utils"

type Container struct {
	ID          uint             `gorm:"primaryKey" json:"id"`
	DockerImage string           `json:"docker_image" binding:"required"`
	Ports       utils.StringMap  `gorm:"type:text" json:"ports"`
	Env         utils.StringMap  `gorm:"type:text" json:"env"`
	Volumes     utils.StringMap  `gorm:"type:text" json:"volumes"`
	Name        string           `json:"name" binding:"required"`
	ProjectID   uint             `json:"project_id"`
	Networks    utils.StringList `gorm:"type:text" json:"networks"`
	NetworkMode string           `json:"network_mode" binding:"required,oneof=bridge host none" gorm:"default:bridge"`
	LastJob     Job              `gorm:"foreignKey:ContainerID" json:"last_job"`
}
