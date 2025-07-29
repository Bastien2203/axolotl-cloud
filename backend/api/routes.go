package api

import (
	"axolotl-cloud/infra/docker"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB, dockerClient *docker.DockerClient) {
	apiGroup := r.Group("/api")
	{
		RegisterProjectRoutes(apiGroup, db)
		RegisterContainerRoutes(apiGroup, db, dockerClient)
	}

	RegisterFrontRoutes(r)
}
