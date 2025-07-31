package api

import (
	"axolotl-cloud/infra/docker"
	"axolotl-cloud/infra/worker"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB, dockerClient *docker.DockerClient, jobWorker *worker.Worker) {
	apiGroup := r.Group("/api")
	{
		RegisterProjectRoutes(apiGroup, db)
		RegisterContainerRoutes(apiGroup, db, dockerClient, jobWorker)
		RegisterJobsRoutes(apiGroup, db, jobWorker)
		RegisterVolumeRoutes(apiGroup, db, dockerClient)
	}

	RegisterFrontRoutes(r)
}
