package handler

import (
	"axolotl-cloud/infra/docker"
	"axolotl-cloud/infra/logger"
	"axolotl-cloud/internal/app/model"
	"axolotl-cloud/internal/app/repository"

	"github.com/gin-gonic/gin"
)

type VolumeHandler struct {
	ContainerRepository *repository.ContainerRepository
	DockerClient        *docker.DockerClient
}

func (h *VolumeHandler) GetVolumes(c *gin.Context) {
	containers, err := h.ContainerRepository.GetAllContainers(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve containers"})
		return
	}
	var allVolumes []*model.Volume
	for _, container := range containers {
		volumes, err := h.DockerClient.ContainerVolumes(c.Request.Context(), container.Name)
		if err != nil {
			logger.Error("Failed to get volumes for container", err)
			c.JSON(500, gin.H{"error": "Failed to retrieve container volumes"})
			return
		}
		for _, volume := range volumes {
			volume.ContainerID = container.ID
			volume.ProjectID = container.ProjectID
		}
		allVolumes = append(allVolumes, volumes...)
	}
	c.JSON(200, allVolumes)

}
