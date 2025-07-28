package project

import (
	"axolotl-cloud/utils"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Repository *ProjectRepository
}

func (h *Handler) CreateProject(c *gin.Context) {
	var project Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	if err := h.Repository.Create(c.Request.Context(), &project); err != nil {
		c.JSON(500, gin.H{"error": "Failed to create project"})
		return
	}
	c.JSON(201, project)
}

func (h *Handler) GetAllProjects(c *gin.Context) {
	projects, err := h.Repository.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve projects"})
		return
	}
	c.JSON(200, projects)
}

func (h *Handler) GetProjectByID(c *gin.Context) {
	id, exists := utils.ParamUInt(c, "id")
	if !exists {
		return
	}

	project, err := h.Repository.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Project not found"})
		return
	}
	c.JSON(200, project)
}

func (h *Handler) UpdateProject(c *gin.Context) {
	id, exists := utils.ParamUInt(c, "id")
	if !exists {
		return
	}

	var project Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	project.ID = id

	if err := h.Repository.Save(c.Request.Context(), &project); err != nil {
		c.JSON(500, gin.H{"error": "Failed to update project"})
		return
	}
	c.JSON(200, project)
}

func (h *Handler) DeleteProject(c *gin.Context) {
	id, exists := utils.ParamUInt(c, "id")
	if !exists {
		return
	}

	if err := h.Repository.Delete(c.Request.Context(), id); err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete project"})
		return
	}
	c.Status(204) // No Content
}
