package main

import (
	"chem-factory/internal/database/sqlite"
	authBootstrap "chem-factory/internal/modules/auth/bootstrap"
	inventoryBootstrap "chem-factory/internal/modules/inventory/bootstrap"
	userBootstrap "chem-factory/internal/modules/user/bootstrap"
	marketBootstrap "chem-factory/internal/modules/market/bootstrap"
	routesHTTP "chem-factory/internal/routes/http"
	jwtManager "chem-factory/utils/jwt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	log.Println("Loading .env")
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	buildServer().Start("127.0.0.1:8090")
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
