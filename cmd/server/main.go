package main

import (
	"chem-factory/internal/database/sqlite"
	authBootstrap "chem-factory/internal/modules/auth/bootstrap"
	inventoryBootstrap "chem-factory/internal/modules/inventory/bootstrap"
	marketBootstrap "chem-factory/internal/modules/market/bootstrap"
	userBootstrap "chem-factory/internal/modules/user/bootstrap"
	routesHTTP "chem-factory/internal/routes/http"
	"chem-factory/pkg/constants"
	jwtManager "chem-factory/utils/jwt"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	log.Println("Loading .env")
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT environment variable is not set")
	}

	serverAddress := fmt.Sprintf("%s:%s", constants.LocalHost, port)

	buildServer().Start(serverAddress)
}

func buildServer() *routesHTTP.Server {
	log.Println("Building server")

	db := sqlite.New()

	jwt := jwtManager.NewJWT(os.Getenv("SECRET_KEY"))
	user := userBootstrap.NewModule(db)
	auth := authBootstrap.NewModule(db)
	inventory := inventoryBootstrap.NewModule(db)
	market := marketBootstrap.NewModule(db)

	return routesHTTP.NewServer(jwt, user, auth, inventory, market)
}
