package main

import (
	"axolotl-cloud/api"
	"axolotl-cloud/infra/db"
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

	r := gin.Default()
	api.RegisterRoutes(r, db)
	r.Run(":" + shared.GetEnv("HTTP_PORT"))
}
