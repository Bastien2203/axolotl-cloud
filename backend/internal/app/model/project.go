package model

import (
	"time"
)

type Project struct {
	ID         uint        `gorm:"primaryKey" json:"id"`
	Name       string      `json:"name" binding:"required"`
	IconURL    string      `json:"icon_url"`
	CreatedAt  time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time   `json:"updated_at" gorm:"autoUpdateTime"`
	Containers []Container `json:"containers" gorm:"foreignKey:ProjectID"`
}
