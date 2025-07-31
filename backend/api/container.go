package api

import (
	"axolotl-cloud/infra/docker"
	"axolotl-cloud/infra/worker"
	"axolotl-cloud/internal/app/handler"
	"axolotl-cloud/internal/app/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterContainerRoutes(r *gin.RouterGroup, db *gorm.DB, dockerClient *docker.DockerClient, w *worker.Worker) {
	containerHandler := &handler.ContainerHandler{
		ContainerRepository: &repository.ContainerRepository{DB: db},
		ProjectRepository:   &repository.ProjectRepository{DB: db},
		DockerClient:        dockerClient,
		JobWorker:           w,
	}
	containerGroup := r.Group("/projects/:id/containers")
	{
		containerGroup.POST("", containerHandler.CreateContainer)
		containerGroup.GET("", containerHandler.GetAllContainers)
		containerGroup.GET("/:containerId", containerHandler.GetContainerByID)
		containerGroup.PUT("/:containerId", containerHandler.UpdateContainer)
		containerGroup.DELETE("/:containerId", containerHandler.DeleteContainer)

		containerGroup.GET("/:containerId/status", containerHandler.GetContainerStatus)
		containerGroup.POST("/:containerId/start", containerHandler.StartContainer)
		containerGroup.POST("/:containerId/stop", containerHandler.StopContainer)
		containerGroup.POST("/:containerId/logs", containerHandler.GetContainerLogs)
		containerGroup.POST("/import", containerHandler.ImportComposeFile)
	}
}
