package main

import (
	"axolotl-cloud/api"
	"axolotl-cloud/infra/db"
	"axolotl-cloud/infra/docker"
	"axolotl-cloud/infra/shared"
	"axolotl-cloud/infra/worker"
	"axolotl-cloud/internal/app/repository"
	"context"

	"github.com/gin-gonic/gin"
)

func init() {
	if err := shared.LoadEnv(); err != nil {
		panic(err)
	}
}

func main() {
	db, err := db.InitDB()
	if err != nil {
		panic(err)
	}

	dockerClient, err := docker.NewDockerClient()
	if err != nil {
		panic(err)
	}
	defer dockerClient.Close()

	jobWorker := worker.NewWorker(10, &repository.JobRepository{DB: db})
	ctx := context.Background()
	jobWorker.Start(ctx)

	r := gin.Default()
	api.RegisterMiddlewares(r)
	api.RegisterRoutes(r, db, dockerClient, jobWorker)
	r.Run(":" + shared.GetEnv("HTTP_PORT"))
}
