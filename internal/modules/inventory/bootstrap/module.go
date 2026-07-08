package bootstrap

import (
	database "chem-factory/internal/database/sqlite"
	"chem-factory/internal/modules/inventory/adapter/http/handler"
	"chem-factory/internal/modules/inventory/adapter/sqlite"
	materialRepository "chem-factory/internal/modules/material/adapter/sqlite"
	"chem-factory/internal/modules/inventory/core/usecase"
	"log"
)

type Module struct {
	InventoryHandler *handler.InventoryHandler
}

func NewModule(db *database.Database) Module {
	log.Println("Creating inventory module")

	inventoryRepo := sqlite.NewInventoryRepo(db)
	materialRepo := materialRepository.NewMaterialRepo(db)

	inventoryUsecase := usecase.NewInventoryUsecase(inventoryRepo, materialRepo, db)

	inventoryHandler := handler.NewInventoryHandler(inventoryUsecase)

	return Module{InventoryHandler: inventoryHandler}
}
