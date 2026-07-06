package bootstrap

import (
	database "chem-factory/internal/database/sqlite"
	"chem-factory/internal/modules/user/adapter/http/handler"
	"chem-factory/internal/modules/user/adapter/sqlite"
	"chem-factory/internal/modules/user/core/usecase"
	"log"
)

type Module struct {
	UserHandler *handler.UserHandler
}

func NewModule(db *database.Database) Module {
	log.Println("Creating user module")

	userRepo := sqlite.NewUserRepo(db)

	userUsecase := usecase.NewUserUsecase(userRepo, db)

	userHandler := handler.NewUserHandler(userUsecase)

	return Module{UserHandler: userHandler}
}
