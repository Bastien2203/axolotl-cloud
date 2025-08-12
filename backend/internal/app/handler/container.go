package handler

import (
	"axolotl-cloud/infra/docker"
	"axolotl-cloud/infra/git"
	"axolotl-cloud/infra/logger"
	"axolotl-cloud/infra/worker"
	"axolotl-cloud/internal/app/model"
	"axolotl-cloud/internal/app/repository"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"axolotl-cloud/utils"

	"github.com/gin-gonic/gin"
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
	container.Name = utils.FormatContainerName(project.Name, container.Name)
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

	project, err := h.ProjectRepository.FindByID(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Project not found"})
		return
	}

	createdContainers, _, err := utils.ParseComposeFileFromBytes([]byte(request.ComposeFile), project)
	if err != nil {
		logger.Error("Failed to parse compose file", err)
		c.JSON(500, gin.H{"error": "Failed to create containers from compose file"})
		return
	}

	c.JSON(201, createdContainers)
}

type RequestBuildFromSource struct {
	GitURL      string `json:"git_url" binding:"required"`
	AccessToken string `json:"access_token" binding:"omitempty"`
}

func (h *ContainerHandler) BuildFromSource(c *gin.Context) {
	projectID, exists := utils.ParamUInt(c, "id")
	if !exists {
		c.JSON(400, gin.H{"error": "Invalid project ID"})
		return
	}

	var body RequestBuildFromSource
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request data"})
		return
	}

	project, err := h.ProjectRepository.FindByID(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Project not found"})
		return
	}

	h.JobWorker.AddJob(&model.Job{
		Name: fmt.Sprintf("Clone and build image from source for project %d", projectID),
		Run: func(ctx context.Context, log *logger.Logger) error {
			log.Info("Cloning repository from %s", body.GitURL)
			dir, err := git.CloneRepository(body.GitURL, fmt.Sprintf("project-%d", projectID), body.AccessToken)
			if err != nil {
				return fmt.Errorf("failed to clone repository: %w", err)
			}
			log.Info("Successfully cloned repository to %s", dir)

			log.Info("Checking for Compose file in cloned directory")
			var composeFile model.ComposeFile
			var parsedContainers []model.Container
			var composeErr error
			hasComposeFile := false

			// Check for compose files
			composePaths := []string{
				filepath.Join(dir, "compose.yaml"),
				filepath.Join(dir, "docker-compose.yaml"),
			}
			var composeContent []byte
			for _, path := range composePaths {
				if content, err := os.ReadFile(path); err == nil {
					composeContent = content
					hasComposeFile = true
					log.Info("Found Compose file at %s", path)
					break
				}
			}

			if hasComposeFile {
				composeFile, parsedContainers, composeErr = utils.ParseComposeFileFromBytes(composeContent, project)
				if composeErr != nil {
					return fmt.Errorf("failed to parse compose file: %w", composeErr)
				}
				log.Info("Successfully parsed Compose file with %d services", len(composeFile.Services))
			}

			// Case: Compose file exists (single or multiple services)
			if hasComposeFile {
				for serviceName, service := range composeFile.Services {
					if service.Build != nil {
						// Build service image from compose definition
						buildContext := filepath.Join(dir, service.Build.Context)
						dockerfile := service.Build.Dockerfile
						if dockerfile == "" {
							dockerfile = "Dockerfile"
						}

						// Determine image name (use compose image name or generate)
						imageName := service.Image
						if imageName == "" {
							imageName = fmt.Sprintf("project-%d-%s", projectID, serviceName)
						}

						log.Info("Building image for service %s from %s", serviceName, buildContext)
						if err := h.DockerClient.BuildImage(ctx, buildContext, dockerfile, imageName, log); err != nil {
							return fmt.Errorf("failed to build image for service %s: %w", serviceName, err)
						}
						log.Info("Successfully built image %s", imageName)

						// Update container model with actual image name
						for i, c := range parsedContainers {
							if c.Name == utils.FormatContainerName(project.Name, serviceName) {
								parsedContainers[i].DockerImage = imageName
							}
						}
					} else if service.Image == "" {
						return fmt.Errorf("service %s has no image and no build context", serviceName)
					}
				}

				// Create all containers from compose definition
				for _, container := range parsedContainers {
					if err := h.ContainerRepository.Create(ctx, &container); err != nil {
						return fmt.Errorf("failed to create container %s: %w", container.Name, err)
					}
					log.Info("Created container %s", container.Name)
				}
			} else {
				// Case: No compose file - build from root Dockerfile
				imageName := fmt.Sprintf("project-%d-image", projectID)
				log.Info("Building image from source directory %s", dir)
				if err := h.DockerClient.BuildImage(ctx, dir, "Dockerfile", imageName, log); err != nil {
					return fmt.Errorf("failed to build image from source: %w", err)
				}
				log.Info("Successfully built image %s", imageName)

				// Create default container
				container := model.Container{
					ProjectID:   projectID,
					Name:        utils.FormatContainerName(project.Name, "default"),
					DockerImage: imageName,
					Ports:       make(map[string]string),
					Env:         make(map[string]string),
					Volumes:     make(map[string]string),
					Networks:    []string{},
					NetworkMode: "bridge",
				}
				if err := h.ContainerRepository.Create(ctx, &container); err != nil {
					return fmt.Errorf("failed to create container: %w", err)
				}
				log.Info("Created default container %s", container.Name)
			}

			return nil
		},
	}, nil)

	c.JSON(200, gin.H{"message": "Job started"})
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
		logger.Error("Failed to delete container", err)
		c.JSON(500, gin.H{"error": "Failed to delete container"})
		return
	}

	jobId, err := h.JobWorker.AddJob(&model.Job{
		Name: fmt.Sprintf("Remove container %s", container.Name),
		Run: func(ctx context.Context, log *logger.Logger) error {
			if err := h.DockerClient.RemoveContainer(ctx, container.Name, log); err != nil {
				return fmt.Errorf("failed to remove container %s: %w", container.Name, err)
			}
			return nil
		},
	}, &containerID)
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

	containerIsCreated, err := h.DockerClient.ContainerExists(c.Request.Context(), container.Name, logger.GlobalLogger)
	if err != nil {
		logger.Error("Failed to check if container exists", err)
		c.JSON(500, gin.H{"error": "Failed to check if container exists"})
		return
	}

	jobId, err := h.JobWorker.AddJob(&model.Job{
		Name: fmt.Sprintf("Start container %s", container.Name),
		Run: func(ctx context.Context, log *logger.Logger) error {
			if containerIsCreated {
				log.Info("Container already exists, removing it before starting a new one")
				if err := h.DockerClient.RemoveContainer(ctx, container.Name, log); err != nil {
					return fmt.Errorf("failed to remove container %s: %w", container.Name, err)
				}
			}

			if _, err := h.DockerClient.CreateContainer(ctx, container.Name, container.DockerImage, container.Ports, container.Env, container.Volumes, container.NetworkMode, log); err != nil {
				return fmt.Errorf("failed to create container %s: %w", container.Name, err)
			}

			if _, err := h.DockerClient.StartContainer(ctx, container.Name, log); err != nil {
				return fmt.Errorf("failed to start container %s: %w", container.Name, err)
			}
			return nil
		},
	}, &containerID)

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
		Run: func(ctx context.Context, log *logger.Logger) error {
			if err := h.DockerClient.StopContainer(ctx, container.Name, log); err != nil {
				return fmt.Errorf("failed to stop container %s: %w", container.Name, err)
			}
			return nil
		},
	}, &containerID)
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
		logger.Error("Failed to get container logs", err)
		c.JSON(500, gin.H{"error": "Failed to get container logs"})
		return
	}

	c.String(200, logs)
}
