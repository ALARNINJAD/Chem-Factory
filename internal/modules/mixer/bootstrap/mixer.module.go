package bootstrap

import (
	database "chem-factory/internal/database/sqlite"
	inventorySQLite "chem-factory/internal/modules/inventory/adapter/sqlite"
	materialSQLite "chem-factory/internal/modules/material/adapter/sqlite"
	"chem-factory/internal/modules/mixer/adapter/http/handler"
	mixerSQLite "chem-factory/internal/modules/mixer/adapter/sqlite"
	"chem-factory/internal/modules/mixer/core/usecase"
	userSQLite "chem-factory/internal/modules/user/adapter/sqlite"
	"log"
)

type Module struct {
	MixerHandler *handler.MixerHandler
}

func NewModule(db *database.Database) Module {
	log.Println("Creating mixer module")

	mixerRepo := mixerSQLite.NewMixerRepo(db)
	userRepo := userSQLite.NewUserRepo(db)
	materialRepo := materialSQLite.NewMaterialRepo(db)
	inventoryRepo := inventorySQLite.NewInventoryRepo(db)

	mixerUsecase := usecase.NewMixerUsecase(
		mixerRepo,
		userRepo,
		materialRepo,
		inventoryRepo,
		db,
	)

	mixerHandler := handler.NewMixerHandler(mixerUsecase)

	return Module{MixerHandler: mixerHandler}
}
