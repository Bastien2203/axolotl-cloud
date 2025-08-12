package api

import (
	"axolotl-cloud/internal/app/handler"
	"axolotl-cloud/internal/app/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterSettingRoutes(router *gin.RouterGroup, db *gorm.DB) {
	settingHandler := &handler.SettingHandler{
		SettingRepository: repository.NewSettingRepository(db),
	}
	settingGroup := router.Group("/settings")
	{
		settingGroup.GET("", settingHandler.GetAllSettings)
		settingGroup.GET("/:key", settingHandler.GetSettingByKey)
		settingGroup.POST("", settingHandler.SaveSetting)
		settingGroup.DELETE("/:key", settingHandler.DeleteSetting)
	}
}
