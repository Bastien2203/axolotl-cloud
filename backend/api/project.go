package api

import (
	"axolotl-cloud/internal/app/handler"
	"axolotl-cloud/internal/app/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterProjectRoutes(r *gin.RouterGroup, db *gorm.DB) {
	projectHandler := &handler.ProjectHandler{
		ProjectRepository: &repository.ProjectRepository{DB: db},
	}
	projectGroup := r.Group("/projects")
	{
		projectGroup.GET("", projectHandler.GetAllProjects)
		projectGroup.GET("/:id", projectHandler.GetProjectByID)
		projectGroup.POST("", projectHandler.CreateProject)
		projectGroup.PUT("/:id", projectHandler.UpdateProject)
		projectGroup.DELETE("/:id", projectHandler.DeleteProject)
	}
}
