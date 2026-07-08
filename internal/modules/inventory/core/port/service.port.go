package port

import (
	"chem-factory/internal/modules/inventory/adapter/http/dto"
	"context"
)

type InventoryService interface {
	Export(ctx context.Context, userID int) (dto.ExportResponse, error)
}
