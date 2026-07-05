package port

import (
	"chem-factory/internal/modules/auth/adapter/http/dto"
	"context"
)

type AuthService interface {
	Register(ctx context.Context, request dto.RegisterRequest) error
	Login(ctx context.Context, request dto.LoginRequest) (int, error)
}
