package port

import (
	"chem-factory/internal/domain"
	"chem-factory/internal/modules/market/adapter/http/dto"
	"context"
)

type UserRepository interface {
	FindUsernameByID(ctx context.Context, id uint) (string, error)
	ReduceBalanceByID(ctx context.Context, id uint, amount int) error
	IncreaseBalanceByID(ctx context.Context, id uint, amount int) error
}

type Transactor interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type MaterialRepository interface {
	FindNameByID(ctx context.Context, id uint) (string, error)
}

type MarketService interface {
	Export(ctx context.Context) (dto.MarketListResponse, error)
	Buy(ctx context.Context, request dto.BuyRequest, userID uint) error
	SetForSell(ctx context.Context, request dto.SetForSellRequest, userID uint) error
}

type MarketRepository interface {
	Export(ctx context.Context) ([]domain.Market, error)
	FindIDByPriceUserIDMatID(ctx context.Context, price int, userID, materialID uint) (uint, error)
	FindByPriceUserIDMatID(ctx context.Context, price int, userID, materialID uint) (domain.Market, error)
	FindByID(ctx context.Context, id uint) (domain.Market, error)
	FindMatIDByID(ctx context.Context, id uint) (uint, error)
	ReduceAmountByID(ctx context.Context, id uint, amount int) error
	IncreaseAmountByID(ctx context.Context, id uint, amount int) error
	Add(ctx context.Context, market domain.Market) error
	DeleteByID(ctx context.Context, id uint) error
	FindUserIDByID(ctx context.Context, id uint) (uint, error)
}

type InventoryRepository interface {
	Add(ctx context.Context, inventory domain.Inventory) error
	IncreaseByID(ctx context.Context, id uint, amount int) error
	ReduceByID(ctx context.Context, id uint, amount int) error
	FindIDByUserIDmatID(ctx context.Context, userID, materialID uint) (uint, error)
}
