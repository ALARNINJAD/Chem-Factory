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

	routes.New(
		gin.Default(), service.New(
			auth.New(), repository.New(), notification.New())).Start()
}

// repository interface
// transaction structure