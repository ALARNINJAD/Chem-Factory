package port

import (
	"chem-factory/internal/modules/user/adapter/http/dto"
	"context"
)

type UserService interface {
	GetProfile(ctx context.Context, userID int) (dto.ProfileResponse, error)
}
