package port

import (
	"chem-factory/internal/domain"
	"chem-factory/internal/modules/inventory/adapter/http/dto"
	"context"
)

type InventoryRepository interface {
	FindByID(ctx context.Context, id uint) (domain.Inventory, error)
	FindByUserID(ctx context.Context, userID uint) ([]domain.Inventory, error)
	FindIDByUserIDmatID(ctx context.Context, userID, materialID uint) (uint, error)
	HasByUserIDMaterialID(ctx context.Context, userID, materialID uint) (bool, error)
	IncreaseByID(ctx context.Context, id uint, amount int) error
	ReduceByID(ctx context.Context, id uint, amount int) error
	Add(ctx context.Context, inventory domain.Inventory) error
	DeleteByID(ctx context.Context, id uint) error
}

type MaterialRepository interface {
	FindNameByID(ctx context.Context, id uint) (string, error)
}

type InventoryService interface {
	Export(ctx context.Context, userID uint) (dto.ExportResponse, error)
}

type Transactor interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}
