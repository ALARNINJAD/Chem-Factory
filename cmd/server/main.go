package main

import (
	"chem-factory/internal/repository"
	"chem-factory/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	repository.Init()
	// repository.NewUser("alarninjad", "12345678")
	server := gin.Default()
	routes.Register(server)
	server.Run(":8090")
}
