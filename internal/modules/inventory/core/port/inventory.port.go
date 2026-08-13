package port

import (
	"chem-factory/internal/domain"
	"chem-factory/internal/modules/inventory/adapter/http/dto"
	"context"
)

type InventoryRepository interface {
	FindByUserID(ctx context.Context, userID uint) ([]domain.Inventory, error)
}

type MaterialRepository interface {
	FindNameByID(ctx context.Context, id uint) (string, error)
}

type InventoryService interface {
	Export(ctx context.Context, userID uint) (dto.ExportResponse, error)
}
