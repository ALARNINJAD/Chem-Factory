package port

import (
	"chem-factory/internal/domain"
	"chem-factory/internal/modules/auth/adapter/http/dto"
	"context"
)

type UserRepository interface {
	FindIDByUsername(ctx context.Context, username string) (uint, error)
	FindPasswordByID(ctx context.Context, userID uint) (string, error)
	Add(ctx context.Context, user domain.User) error
}

type AuthService interface {
	Register(ctx context.Context, request dto.RegisterRequest) error
	Login(ctx context.Context, request dto.LoginRequest) (uint, error)
}

type Transactor interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}
