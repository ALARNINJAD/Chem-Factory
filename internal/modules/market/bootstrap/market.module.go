package bootstrap

import (
	database "chem-factory/internal/database/sqlite"
	"chem-factory/internal/modules/market/adapter/http/handler"
	"chem-factory/internal/modules/market/core/usecase"
	materialRepository "chem-factory/internal/modules/material/adapter/sqlite"
	userRepository "chem-factory/internal/modules/user/adapter/sqlite"
	marketRepository "chem-factory/internal/modules/market/adapter/sqlite"
	inventoryRepository "chem-factory/internal/modules/inventory/adapter/sqlite"
	"log"
)

type Module struct {
	MarketHandler *handler.MarketHandler
}

func NewModule(db *database.Database) Module {
	log.Println("Creating market module")

	materialRepo := materialRepository.NewMaterialRepo(db)
	userRepo := userRepository.NewUserRepo(db)
	marketRepo := marketRepository.NewMarketRepo(db)
	inventoryRepo := inventoryRepository.NewInventoryRepo(db)

	marketUsecase := usecase.NewMarketUsecase(materialRepo, userRepo, marketRepo, inventoryRepo, db)

	marketHandler := handler.NewMarketHandler(marketUsecase)

	return Module{MarketHandler: marketHandler}
}
