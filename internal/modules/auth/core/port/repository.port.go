package port

import (
	"chem-factory/internal/domain"
	"context"
)

type UserRepository interface {
	FindIDByUsername(ctx context.Context, username string) (int, error)
	FindPasswordByID(ctx context.Context, userID int) (string, error)
	Add(ctx context.Context, user domain.User) error
}
