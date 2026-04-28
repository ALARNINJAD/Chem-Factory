package main

import (
	"chem-factory/internal/auth"
	"chem-factory/internal/notification"
	"chem-factory/internal/repository"
	"chem-factory/internal/routes"
	"chem-factory/internal/service"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	routes.Init(
		gin.Default(), service.Init(
			auth.Init(), repository.Init(), notification.Init())).Start()
}
