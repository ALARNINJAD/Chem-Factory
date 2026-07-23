package port

import (
	"chem-factory/internal/domain"
	"chem-factory/internal/modules/user/adapter/http/dto"
	"context"
)

type UserRepository interface {
	FindByID(ctx context.Context, id uint) (domain.User, error)
}

type UserService interface {
	GetProfile(ctx context.Context, userID uint) (dto.ProfileResponse, error)
}
