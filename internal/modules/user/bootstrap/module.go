package bootstrap

import (
	database "chem-factory/internal/database/sqlite"
	"chem-factory/internal/modules/user/adapter/http/handler"
	"chem-factory/internal/modules/user/adapter/sqlite"
	"chem-factory/internal/modules/user/core/usecase"
)

type Module struct {
	UserHandler *handler.UserHandler
}

func NewModule(db *database.Database) Module {

	userRepo := sqlite.NewUserRepo(db)

	userUsecase := usecase.NewUserUsecase(userRepo, db)

	userHandler := handler.NewUserHandler(userUsecase)

	return Module{UserHandler: userHandler}
}
