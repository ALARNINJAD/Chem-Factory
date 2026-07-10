package port

import (
	"chem-factory/internal/domain"
	"chem-factory/internal/modules/user/adapter/http/dto"
	"context"
)

type UserRepository interface {
	FindByID(ctx context.Context, id int) (domain.User, error)
}

type Transactor interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type UserService interface {
	GetProfile(ctx context.Context, userID int) (dto.ProfileResponse, error)
}
