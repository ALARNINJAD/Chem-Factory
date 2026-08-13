package port

import (
	"chem-factory/internal/domain"
	"chem-factory/internal/modules/mixer/adapter/http/dto"
	"context"
)

type MixerRepository interface {
	FindByID(ctx context.Context, id uint) (domain.Mix, error)
	Add(ctx context.Context, mixer domain.Mix) error
	DeleteByID(ctx context.Context, id uint) error
	GetByUserID(ctx context.Context, userID uint) ([]domain.Mix, error)
	FindByUserIDIngrID(ctx context.Context, userID, firstID, secID uint) (domain.Mix, error)
}

type UserRepository interface {
	FindUsernameByID(ctx context.Context, id uint) (string, error)
	UpdateLevelXPByID(ctx context.Context, user domain.User) error
	FindByID(ctx context.Context, id uint) (domain.User, error)
}

type MaterialRepository interface {
	FindNameByID(ctx context.Context, id uint) (string, error)
	FindByIngrID(ctx context.Context, firstID uint, secondID uint) (domain.Material, error)
	Add(ctx context.Context, material domain.Material) error
}

type InventoryRepository interface {
	IncreaseByID(ctx context.Context, id uint, amount int) error
	ReduceByID(ctx context.Context, id uint, amount int) error
	Add(ctx context.Context, inventory domain.Inventory) error
	DeleteByID(ctx context.Context, id uint) error
	FindByUserIDmatID(ctx context.Context, userID, materialID uint) (domain.Inventory, error)
}

type Transactor interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type MixerService interface {
	Mixes(ctx context.Context, userID uint) (dto.MixesResponse, error)
	Mix(ctx context.Context, request dto.MixRequest, userID uint) (dto.MixResponse, error)
	Check(ctx context.Context, request dto.CheckRequest, userID uint) (dto.MixResponse, error)
	Pick(ctx context.Context, request dto.PickRequest, userID uint) (dto.PickResponse, error)
	NewMaterial(ctx context.Context, request dto.NewMaterialRequest, userID uint) (dto.MixResponse, error)
}
