package api

import (
	"axolotl-cloud/infra/websocket"

	"github.com/gin-gonic/gin"
)

func RegisterWebSocketRoutes(r *gin.Engine, wss *websocket.WebSocketServer) {
	r.GET("/ws", func(c *gin.Context) {
		wss.HandleHTTP(c.Writer, c.Request)
	})
}
