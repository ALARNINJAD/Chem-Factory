package bootstrap

import (
	database "chem-factory/internal/database/sqlite"
	"chem-factory/internal/modules/auth/adapter/http/handler"
	"chem-factory/internal/modules/auth/core/usecase"
	"chem-factory/internal/modules/user/adapter/sqlite"
	"chem-factory/utils/jwt"
	"os"
)

type Module struct {
	AuthHandler *handler.AuthHandler
}

func NewModule(db *database.Database) Module {

	userRepo := sqlite.NewUserRepo(db)

	authUsecase := usecase.NewAuthUsecase(userRepo, db)
	jwtService := jwt.NewJWT(os.Getenv("SECRET_KEY"))

	authHandler := handler.NewAuthHandler(authUsecase, jwtService)

	return Module{AuthHandler: authHandler}
}
