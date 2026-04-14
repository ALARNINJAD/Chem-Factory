package main

import (
	"chem-factory/internal/auth"
	"chem-factory/internal/repository"
	"chem-factory/internal/routes"
	"chem-factory/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	routes.Init(
		gin.Default(), service.Init(
			auth.Init(), repository.Init())).Start()
}
