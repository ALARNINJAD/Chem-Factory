package main

import (
	"chem-factory/internal/repository"
	"chem-factory/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	repository.Init()
	server := gin.Default()
	routes.Register(server)
	server.Run(":8090")
}
