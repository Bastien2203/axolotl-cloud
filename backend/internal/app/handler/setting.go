package handler

import (
	"axolotl-cloud/internal/app/model"
	"axolotl-cloud/internal/app/repository"

	"github.com/gin-gonic/gin"
)

type SettingHandler struct {
	SettingRepository *repository.SettingRepository
}

func (h *SettingHandler) GetAllSettings(c *gin.Context) {
	settings, err := h.SettingRepository.GetAll()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve settings"})
		return
	}
	c.JSON(200, settings)
}

func (h *SettingHandler) GetSettingByKey(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(400, gin.H{"error": "Invalid setting key"})
		return
	}

	setting, err := h.SettingRepository.GetByKey(model.SettingKey(key))
	if err != nil {
		c.JSON(404, gin.H{"error": "Setting not found"})
		return
	}
	c.JSON(200, setting)
}

func (h *SettingHandler) SaveSetting(c *gin.Context) {
	var setting model.Setting
	if err := c.ShouldBindJSON(&setting); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	if err := h.SettingRepository.Save(&setting); err != nil {
		c.JSON(500, gin.H{"error": "Failed to save setting"})
		return
	}

	c.JSON(200, setting)
}

func (h *SettingHandler) DeleteSetting(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(400, gin.H{"error": "Invalid setting key"})
		return
	}

	if err := h.SettingRepository.RemoveByKey(key); err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete setting"})
		return
	}

	c.Status(204)
}
