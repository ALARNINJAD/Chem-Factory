package port

import "context"

type MaterialRepository interface {
	FindNameByID(ctx context.Context, id int) (string, error)
}
