package main

import (
	"chem-factory/internal/auth"
	"chem-factory/internal/notification"
	"chem-factory/internal/repository"
	"chem-factory/internal/routes"
	"chem-factory/internal/service"
	"log"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	routes.New(
		service.New(
			auth.New(), repository.New(), notification.New())).Start()
}

// transaction context
// middleware
// handler
// remove extra columns
// acid repo
// cun and chan
// try https
