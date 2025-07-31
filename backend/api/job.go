package api

import (
	"axolotl-cloud/infra/worker"
	"axolotl-cloud/internal/app/handler"
	"axolotl-cloud/internal/app/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterJobsRoutes(r *gin.RouterGroup, db *gorm.DB, w *worker.Worker) {
	jobHandler := &handler.JobHandler{
		JobRepository: &repository.JobRepository{DB: db},
		Worker:        w,
	}
	jobGroup := r.Group("/jobs")
	{
		jobGroup.GET("", jobHandler.GetAllJobs)
		jobGroup.GET("/:id", jobHandler.GetJobByID)
		jobGroup.DELETE("/:id", jobHandler.DeleteJob)
	}
}
