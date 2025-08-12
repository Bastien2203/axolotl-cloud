package main

import (
	"axolotl-cloud/api"
	"axolotl-cloud/infra/db"
	"axolotl-cloud/infra/docker"
	"axolotl-cloud/infra/logger"
	"axolotl-cloud/infra/settings"
	"axolotl-cloud/infra/shared"
	"axolotl-cloud/infra/worker"
	"axolotl-cloud/internal/app/repository"
	"context"
	"fmt"
	"strconv"
	"time"

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
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	jobWorker.Start(ctx)

	r := gin.Default()
	api.RegisterMiddlewares(r)
	api.RegisterRoutes(r, db, dockerClient, jobWorker)
	r.Run(":" + shared.GetEnv("HTTP_PORT"))
}
