package main

import (
	"chem-factory/internal/auth"
	"chem-factory/internal/repository"
	"log"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := repository.Init()
	createTables(r)
	
	a := auth.Init()
	createAdminUser(r, a)
	
	createAdminMaterials(r)
}
