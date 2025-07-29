package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterFrontRoutes(r *gin.Engine) {
	r.LoadHTMLFiles("./dist/index.html")
	r.GET("/icon.png", func(c *gin.Context) {
		c.File("./dist/icon.png")
	})
	r.Static("/assets", "./dist/assets")
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
}
