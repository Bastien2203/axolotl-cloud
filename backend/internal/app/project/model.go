package project

import "time"

type Project struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `json:"name" binding:"required"`
	RepoURL    string    `json:"repo_url" binding:"required"`
	DeployMode string    `json:"deploy_mode" binding:"required"`
	Port       int       `json:"port" binding:"required"`
	Env        string    `json:"env" binding:"required"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
