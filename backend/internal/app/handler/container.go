package handler

import (
	"axolotl-cloud/infra/docker"
	"axolotl-cloud/infra/logger"
	"axolotl-cloud/infra/worker"
	"axolotl-cloud/internal/app/model"
	"axolotl-cloud/internal/app/repository"
	"context"
	"fmt"
	"strings"

	"axolotl-cloud/utils"
	"regexp"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

type ContainerHandler struct {
	ContainerRepository *repository.ContainerRepository
	ProjectRepository   *repository.ProjectRepository
	JobWorker           *worker.Worker
	DockerClient        *docker.DockerClient
}

func (h *ContainerHandler) CreateContainer(c *gin.Context) {
	projectID, exists := utils.ParamUInt(c, "id")
	if !exists {
		c.JSON(400, gin.H{"error": "Invalid project ID"})
		return
	}

	project, err := h.ProjectRepository.FindByID(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Project not found"})
		return
	}

	var container model.Container
	if err := c.ShouldBindJSON(&container); err != nil {
		c.JSON(400, gin.H{"error": "Invalid container data"})
		return
	}
	container.ProjectID = projectID
	container.Name = formatContainerName(project.Name, container.Name)
	if err := h.ContainerRepository.Create(c.Request.Context(), &container); err != nil {
		c.JSON(500, gin.H{"error": "Failed to create container"})
		return
	}

	c.JSON(201, container)
}

func (h *ContainerHandler) ImportComposeFile(c *gin.Context) {
	projectID, exists := utils.ParamUInt(c, "id")
	if !exists {
		c.JSON(400, gin.H{"error": "Invalid project ID"})
		return
	}

	var request struct {
		ComposeFile string `json:"compose_file"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request data"})
		return
	}

	var compose model.ComposeFile
	if err := yaml.Unmarshal([]byte(request.ComposeFile), &compose); err != nil {
		logger.Error("Failed to parse compose file", err)
		c.JSON(400, gin.H{"error": "Invalid YAML"})
		return
	}

	project, err := h.ProjectRepository.FindByID(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Project not found"})
		return
	}

	var createdContainers []model.Container
	for name, service := range compose.Services {
		container := model.Container{
			ProjectID:   projectID,
			Name:        formatContainerName(project.Name, name),
			DockerImage: service.Image,
			Ports:       parsePorts(service.Ports),
			Env:         service.Env,
			Volumes:     parseVolumes(service.Volumes),
			Networks:    parseNetworks(service.Networks),
			NetworkMode: service.NetworkMode,
		}
		if err := h.ContainerRepository.Create(c.Request.Context(), &container); err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to create container %s", name)})
			return
		}
		createdContainers = append(createdContainers, container)
	}

	c.JSON(201, createdContainers)
}

func (h *ContainerHandler) GetAllContainers(c *gin.Context) {
	projectID, exists := utils.ParamUInt(c, "id")
	if !exists {
		c.JSON(400, gin.H{"error": "Invalid project ID"})
		return
	}

	containers, err := h.ContainerRepository.FindAllByProjectID(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve containers"})
		return
	}

	c.JSON(200, containers)
}

func (h *ContainerHandler) GetContainerByID(c *gin.Context) {
	containerID, exists := utils.ParamUInt(c, "containerId")
	if !exists {
		c.JSON(400, gin.H{"error": "Invalid container ID"})
		return
	}

	container, err := h.ContainerRepository.FindByID(c.Request.Context(), containerID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Container not found"})
		return
	}
	c.JSON(200, container)
}

func (h *ContainerHandler) UpdateContainer(c *gin.Context) {
	projectID, exists := utils.ParamUInt(c, "id")
	if !exists {
		c.JSON(400, gin.H{"error": "Invalid project ID"})
		return
	}

	containerID, exists := utils.ParamUInt(c, "containerId")
	if !exists {
		c.JSON(400, gin.H{"error": "Invalid container ID"})
		return
	}

	var container model.Container
	if err := c.ShouldBindJSON(&container); err != nil {
		c.JSON(400, gin.H{"error": "Invalid container data"})
		return
	}
	container.ID = containerID
	container.ProjectID = projectID

	if err := h.ContainerRepository.Save(c.Request.Context(), &container); err != nil {
		c.JSON(500, gin.H{"error": "Failed to update container"})
		return
	}
}

func (h *ContainerHandler) DeleteContainer(c *gin.Context) {
	containerID, exists := utils.ParamUInt(c, "containerId")
	if !exists {
		c.JSON(400, gin.H{"error": "Invalid container ID"})
		return
	}

	container, err := h.ContainerRepository.FindByID(c.Request.Context(), containerID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Container not found"})
		return
	}

	if err := h.ContainerRepository.Delete(c.Request.Context(), containerID); err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete container"})
		return
	}

	jobId, err := h.JobWorker.AddJob(&model.Job{
		Name: fmt.Sprintf("Remove container %s", container.Name),
		Run: func(ctx context.Context, log func(string)) error {
			if err := h.DockerClient.RemoveContainer(ctx, container.Name, log); err != nil {
				return fmt.Errorf("failed to remove container %s: %w", container.Name, err)
			}
			return nil
		},
	}, containerID)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to add job to remove container %s", container.Name)})
		return
	}

	c.JSON(201, gin.H{
		"job_id": jobId,
	})
}

/// Docker operations

func (h *ContainerHandler) GetContainerStatus(c *gin.Context) {
	containerID, exists := utils.ParamUInt(c, "containerId")
	if !exists {
		c.JSON(400, gin.H{"error": "Invalid container ID"})
		return
	}

	container, err := h.ContainerRepository.FindByID(c.Request.Context(), containerID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Container not found"})
		return
	}

	status, err := h.DockerClient.ContainerStatus(c.Request.Context(), container.Name)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get container status"})
		return
	}

	c.JSON(200, gin.H{"status": status})
}

func (h *ContainerHandler) StartContainer(c *gin.Context) {
	containerID, exists := utils.ParamUInt(c, "containerId")
	if !exists {
		c.JSON(400, gin.H{"error": "Invalid container ID"})
		return
	}

	container, err := h.ContainerRepository.FindByID(c.Request.Context(), containerID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Container not found"})
		return
	}

	containerIsCreated, err := h.DockerClient.ContainerExists(c.Request.Context(), container.Name, logger.Info)
	if err != nil {
		logger.Error("Failed to check if container exists", err)
		c.JSON(500, gin.H{"error": "Failed to check if container exists"})
		return
	}

	jobId, err := h.JobWorker.AddJob(&model.Job{
		Name: fmt.Sprintf("Start container %s", container.Name),
		Run: func(ctx context.Context, log func(string)) error {
			if !containerIsCreated {
				if _, err := h.DockerClient.CreateContainer(ctx, container.Name, container.DockerImage, container.Ports, container.Env, container.Volumes, container.NetworkMode, log); err != nil {
					return fmt.Errorf("failed to create container %s: %w", container.Name, err)
				}
			}

			if _, err := h.DockerClient.StartContainer(ctx, container.Name, log); err != nil {
				return fmt.Errorf("failed to start container %s: %w", container.Name, err)
			}
			return nil
		},
	}, containerID)

	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to add job to start container %s", container.Name)})
		return
	}

	c.JSON(201, gin.H{
		"job_id": jobId,
	})
}

func (h *ContainerHandler) StopContainer(c *gin.Context) {
	containerID, exists := utils.ParamUInt(c, "containerId")
	if !exists {
		c.JSON(400, gin.H{"error": "Invalid container ID"})
		return
	}

	container, err := h.ContainerRepository.FindByID(c.Request.Context(), containerID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Container not found"})
		return
	}

	jobId, err := h.JobWorker.AddJob(&model.Job{
		Name: fmt.Sprintf("Stop container %s", container.Name),
		Run: func(ctx context.Context, log func(string)) error {
			if err := h.DockerClient.StopContainer(ctx, container.Name, log); err != nil {
				return fmt.Errorf("failed to stop container %s: %w", container.Name, err)
			}
			return nil
		},
	}, containerID)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to add job to stop container %s", container.Name)})
		return
	}

	c.JSON(201, gin.H{
		"job_id": jobId,
	})
}

func (h *ContainerHandler) GetContainerLogs(c *gin.Context) {
	containerID, exists := utils.ParamUInt(c, "containerId")
	if !exists {
		c.JSON(400, gin.H{"error": "Invalid container ID"})
		return
	}

	container, err := h.ContainerRepository.FindByID(c.Request.Context(), containerID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Container not found"})
		return
	}

	tail := c.Query("tail")
	logs, err := h.DockerClient.GetContainerLogs(c.Request.Context(), container.Name, tail)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get container logs"})
		return
	}

	c.String(200, logs)
}

// utils -------
func formatContainerName(projectName string, containerName string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	cleanProjectName := reg.ReplaceAllString(projectName, "")
	cleanContainerName := reg.ReplaceAllString(containerName, "")

	return cleanProjectName + "_" + cleanContainerName
}

func parsePorts(portDefs []string) map[string]string {
	ports := make(map[string]string)
	for _, def := range portDefs {
		parts := strings.Split(def, ":")
		if len(parts) == 2 {
			ports[parts[0]] = parts[1]
		}
	}
	return ports
}

func parseVolumes(volumeDefs []string) map[string]string {
	volumes := make(map[string]string)
	for _, def := range volumeDefs {
		parts := strings.Split(def, ":")
		if len(parts) == 2 {
			volumes[parts[0]] = parts[1]
		}
	}
	return volumes
}

func parseNetworks(networkDefs []string) []string {
	networks := make([]string, 0, len(networkDefs))
	for _, def := range networkDefs {
		if def != "" {
			networks = append(networks, def)
		}
	}
	return networks
}
