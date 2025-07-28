package api

import (
	"axolotl-cloud/internal/app/project"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterProjectRoutes(r *gin.Engine, db *gorm.DB) {
	handler := &project.Handler{
		Repository: &project.ProjectRepository{DB: db},
	}
	projectGroup := r.Group("/projects")
	{
		projectGroup.GET("/", handler.GetAllProjects)
		projectGroup.GET("/:id", handler.GetProjectByID)
		projectGroup.POST("/", handler.CreateProject)
		projectGroup.PUT("/:id", handler.UpdateProject)
		projectGroup.DELETE("/:id", handler.DeleteProject)
	}
}
