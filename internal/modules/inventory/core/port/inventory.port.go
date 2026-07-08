package port

import (
	"chem-factory/internal/domain"
	"context"
)

type InventoryRepository interface {
	FindByID(ctx context.Context, id int) (domain.Inventory, error)
	FindByUserID(ctx context.Context, userID int) ([]domain.Inventory, error)
	FindIDByUserIDmatID(ctx context.Context, userID, materialID int) (int, error)
	IncreaseByID(ctx context.Context, id, amount int) error 
	ReduceByID(ctx context.Context, id, amount int) error
	Add(ctx context.Context, inventory domain.Inventory) error
	DeleteByID(ctx context.Context, id int) error
}