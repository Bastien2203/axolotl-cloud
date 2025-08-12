package main

import (
	"axolotl-cloud/api"
	"axolotl-cloud/infra/db"
	"axolotl-cloud/infra/docker"
	"axolotl-cloud/infra/logger"
	"axolotl-cloud/infra/settings"
	"axolotl-cloud/infra/shared"
	"axolotl-cloud/infra/websocket"
	"axolotl-cloud/infra/worker"
	"axolotl-cloud/internal/app/repository"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func init() {
	if err := shared.LoadEnv(); err != nil {
		panic(err)
	}
}

func initDB() *gorm.DB {
	db, err := db.InitDB()
	if err != nil {
		panic(err)
	}
	return db
}

func initWorker(db *gorm.DB) (*worker.Worker, context.CancelFunc) {
	jobTimeoutInSecond, err := repository.NewSettingRepository(db).GetByKey(settings.JobTimeout)
	if err != nil {
		panic(err)
	}

	jobWorker := worker.NewWorker(10, &repository.JobRepository{DB: db})
	// context background with job timeout (convert jobTimeoutInSecond.Value to int)
	timeout, err := strconv.Atoi(jobTimeoutInSecond.Value)
	if err != nil {
		panic("Invalid job timeout value: " + jobTimeoutInSecond.Value)
	}
	logger.Info(fmt.Sprintf("Job timeout set to %d seconds", timeout))
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)

	jobWorker.Start(ctx)
	return jobWorker, cancel
}

func initDockerClient() *docker.DockerClient {
	dockerClient, err := docker.NewDockerClient()
	if err != nil {
		panic("Failed to initialize Docker client: " + err.Error())
	}
	return dockerClient
}

func initWSServer() *websocket.WebSocketServer {
	return &websocket.WebSocketServer{
		OnConnect: func(conn websocket.WebSocketConnection) {},
		OnMessage: func(conn websocket.WebSocketConnection, data websocket.WSMessage[any]) {
			websocket.NewWSMessageHandler(conn).HandleMessage(data)
		},
		OnDisconnect: func(conn websocket.WebSocketConnection, err error) {
			if err != nil {
				logger.Error("WebSocket connection error:", err)
			}
			websocket.NewWSMessageHandler(conn).UnsubscribeAllTopics()
		},
	}
}

func main() {
	db := initDB()
	dockerClient := initDockerClient()
	defer dockerClient.Close()
	wss := initWSServer()

	jobWorker, cancel := initWorker(db)
	defer cancel()

	r := gin.Default()
	api.RegisterMiddlewares(r)
	api.RegisterWebSocketRoutes(r, wss)
	api.RegisterRoutes(r, db, dockerClient, jobWorker)
	r.Run(":" + shared.GetEnv("HTTP_PORT"))
}
