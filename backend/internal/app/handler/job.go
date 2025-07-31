package handler

import (
	"axolotl-cloud/infra/worker"
	"axolotl-cloud/internal/app/repository"
	"axolotl-cloud/utils"

	"github.com/gin-gonic/gin"
)

type JobHandler struct {
	JobRepository *repository.JobRepository
	Worker        *worker.Worker
}

func (h *JobHandler) GetAllJobs(c *gin.Context) {
	jobs, err := h.JobRepository.GetAll()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve jobs"})
		return
	}
	c.JSON(200, jobs)
}

func (h *JobHandler) GetJobByID(c *gin.Context) {
	id, exists := utils.ParamUInt(c, "id")
	if !exists {
		c.JSON(400, gin.H{"error": "Invalid job ID"})
		return
	}

	job, err := h.JobRepository.GetByID(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Job not found"})
		return
	}
	c.JSON(200, job)
}

func (h *JobHandler) DeleteJob(c *gin.Context) {
	id, exists := utils.ParamUInt(c, "id")
	if !exists {
		c.JSON(400, gin.H{"error": "Invalid job ID"})
		return
	}

	if err := h.JobRepository.RemoveByID(id); err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete job"})
		return
	}

	c.Status(204)
}
