package port

import (
	"chem-factory/internal/domain"
	"chem-factory/internal/modules/market/adapter/http/dto"
	"context"
)

type UserRepository interface {
	FindUsernameByID(ctx context.Context, id int) (string, error)
	ReduceBalanceByID(ctx context.Context, id, amount int) error
	IncreaseBalanceByID(ctx context.Context, id, amount int) error
}

type Transactor interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type MaterialRepository interface {
	FindNameByID(ctx context.Context, id int) (string, error)
}

type MarketService interface {
	Export(ctx context.Context) (dto.MarketListResponse, error)
	Buy(ctx context.Context, request dto.BuyRequest, userID int) error
	SetForSell(ctx context.Context, request dto.SetForSellRequest, userID int) error
}

type MarketRepository interface {
	Export(ctx context.Context) ([]domain.Market, error)
	FindIDByPriceUserIDMatID(ctx context.Context, price, userID, materialID int) (int, error)
	FindByPriceUserIDMatID(ctx context.Context, price, userID, materialID int) (domain.Market, error)
	FindByID(ctx context.Context, id int) (domain.Market, error)
	FindMatIDByID(ctx context.Context, id int) (int, error)
	ReduceAmountByID(ctx context.Context, id, amount int) error
	IncreaseAmountByID(ctx context.Context, id, amount int) error
	Add(ctx context.Context, market domain.Market) error
	DeleteByID(ctx context.Context, id int) error
	FindUserIDByID(ctx context.Context, id int) (int, error)
}

type InventoryRepository interface {
	Add(ctx context.Context, inventory domain.Inventory) error
	IncreaseByID(ctx context.Context, id, amount int) error
	ReduceByID(ctx context.Context, id, amount int) error
	FindIDByUserIDmatID(ctx context.Context, userID, materialID int) (int, error)
}
