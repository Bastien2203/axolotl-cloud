package handler

import (
	"axolotl-cloud/internal/app/model"
	"axolotl-cloud/internal/app/repository"
	"axolotl-cloud/utils"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	ProjectRepository *repository.ProjectRepository
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var project model.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	if err := h.ProjectRepository.Create(c.Request.Context(), &project); err != nil {
		c.JSON(500, gin.H{"error": "Failed to create project"})
		return
	}

	c.JSON(201, project)
}

func (h *ProjectHandler) GetAllProjects(c *gin.Context) {
	projects, err := h.ProjectRepository.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve projects"})
		return
	}
	c.JSON(200, projects)
}

func (h *ProjectHandler) GetProjectByID(c *gin.Context) {
	id, exists := utils.ParamUInt(c, "id")
	if !exists {
		return
	}

	project, err := h.ProjectRepository.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Project not found"})
		return
	}
	c.JSON(200, project)
}

func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	id, exists := utils.ParamUInt(c, "id")
	if !exists {
		return
	}

	var project model.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	project.ID = id

	if err := h.ProjectRepository.Save(c.Request.Context(), &project); err != nil {
		c.JSON(500, gin.H{"error": "Failed to update project"})
		return
	}
	c.JSON(200, project)
}

func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	id, exists := utils.ParamUInt(c, "id")
	if !exists {
		return
	}

	if err := h.ProjectRepository.Delete(c.Request.Context(), id); err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete project"})
		return
	}

	c.Status(204) // No Content
}
