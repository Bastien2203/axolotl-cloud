package api

import (
	"axolotl-cloud/infra/docker"
	"axolotl-cloud/internal/app/handler"
	"axolotl-cloud/internal/app/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterVolumeRoutes(router *gin.RouterGroup, db *gorm.DB, dockerClient *docker.DockerClient) {
	volumeHandler := &handler.VolumeHandler{
		ContainerRepository: &repository.ContainerRepository{DB: db},
		DockerClient:        dockerClient,
	}
	router.GET("/volumes", volumeHandler.GetVolumes)
}
