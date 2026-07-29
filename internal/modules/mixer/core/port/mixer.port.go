package port

import (
	"chem-factory/internal/domain"
	"context"
)

type MixerRepository interface {
	FindByID(ctx context.Context, id int) (domain.Mix, error)
	Add(ctx context.Context, mixer domain.Mix) error
	FindIDByUserIDIngrID(ctx context.Context, userID, firstID, secID int) (int, error)
	DeleteByID(ctx context.Context, id int) error
}

type Transactor interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}
