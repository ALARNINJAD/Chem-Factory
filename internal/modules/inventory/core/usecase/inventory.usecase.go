package usecase

import (
	"chem-factory/internal/modules/inventory/adapter/http/dto"
	"chem-factory/internal/modules/inventory/core/port"
	"context"
)

type inventoryUsecase struct {
	inventoryRepo port.InventoryRepository
	materialRepo  port.MaterialRepository
	transactor    port.Transactor
}

func NewInventoryUsecase(
	inventoryRepo port.InventoryRepository,
	materialRepo port.MaterialRepository,
	transactor port.Transactor,
) *inventoryUsecase {
	return &inventoryUsecase{
		inventoryRepo: inventoryRepo,
		materialRepo:  materialRepo,
		transactor:    transactor,
	}
}

func (service *inventoryUsecase) Export(ctx context.Context, userID uint) (dto.ExportResponse, error) {

	var response dto.ExportResponse

	InventoryList, err := service.inventoryRepo.FindByUserID(ctx, userID)
	if err != nil {
		return dto.ExportResponse{}, err
	}

	for _, inventory := range InventoryList {

		materialName, err := service.materialRepo.FindNameByID(ctx, inventory.MaterialID)
		if err != nil {
			return dto.ExportResponse{}, err
		}

		response.InventoryList = append(response.InventoryList, dto.InventoryResponse{
			ID:           inventory.ID,
			UserID:       inventory.UserID,
			MaterialID:   inventory.MaterialID,
			MaterialName: materialName,
			Amount:       inventory.Amount,
			DateTime:     inventory.DateTime,
		})
	}

	return response, nil
}
