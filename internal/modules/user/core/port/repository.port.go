package port

import (
	"chem-factory/internal/domain"
	"context"
)

type UserRepository interface {
	FindByID(ctx context.Context, id int) (domain.User, error)
}
