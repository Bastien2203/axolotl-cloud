package main

import (
	"axolotl-cloud/api"
	"axolotl-cloud/infra/db"
	"axolotl-cloud/infra/docker"
	"axolotl-cloud/infra/shared"

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

	r := gin.Default()
	api.RegisterMiddlewares(r)
	api.RegisterRoutes(r, db, dockerClient)
	r.Run(":" + shared.GetEnv("HTTP_PORT"))
}
