package main

import (
	"chem-factory/internal/database/sqlite"
	userBootstrap "chem-factory/internal/modules/user/bootstrap"
	authBootstrap "chem-factory/internal/modules/auth/bootstrap"
	routesHTTP "chem-factory/internal/routes/http"
	"chem-factory/utils/jwt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	BuildServer().Start(":8080")
}

func BuildServer() *routesHTTP.Server {

	db := sqlite.New()

	j := jwt.NewJWT(os.Getenv("SECRET_KEY"))
	user := userBootstrap.NewModule(db)
	auth := authBootstrap.NewModule(db)

	return routesHTTP.NewServer(j, user, auth)
}