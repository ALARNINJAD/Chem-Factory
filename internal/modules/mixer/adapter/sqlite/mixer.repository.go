package sqlite

import (
	database "chem-factory/internal/database/sqlite"
	"chem-factory/internal/domain"
	"context"
	"fmt"
)

type MixRepo struct{ db *database.Database }

func NewMixRepo(db *database.Database) *MixRepo { return &MixRepo{db: db} }

func (r *MixRepo) FindByID(ctx context.Context, id uint) (domain.Mix, error) {
	var mix domain.Mix

	err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT id, user_id, first_ingredient_id, second_ingredient_id, amount, date_time FROM mixes WHERE id = ?`,
		id,
	).Scan(
		&mix.ID,
		&mix.UserID,
		&mix.FirstIngredientID,
		&mix.SecondIngredientID,
		&mix.Amount,
		&mix.DateTime,
	)
	if err != nil {
		return domain.Mix{}, fmt.Errorf("find mix by id: %w", err)
	}

	return mix, nil
}

func (r *MixRepo) Add(ctx context.Context, mix domain.Mix) error {
	_, err := r.db.Extract(ctx).ExecContext(ctx,
		`INSERT INTO mixes(user_id, first_ingredient_id, second_ingredient_id, amount, date_time)
		VALUES (?, ?, ?, ?, ?)`,
		mix.UserID,
		mix.FirstIngredientID,
		mix.SecondIngredientID,
		mix.Amount,
		mix.DateTime,
	)
	if err != nil {
		return fmt.Errorf("add mix: %w", err)
	}

	return nil
}

func (r *MixRepo) FindIDByUserIDIngrID(ctx context.Context, userID, firstID, secID uint) (uint, error) {
	var id uint

	err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT id FROM mixes WHERE user_id = ? AND first_ingredient_id = ? AND second_ingredient_id = ?`,
		userID,
		firstID,
		secID,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("find mix id by user id and ingredients id: %w", err)
	}

	return id, nil
}

func (r *MixRepo) DeleteByID(ctx context.Context, id uint) error {
	_, err := r.db.Extract(ctx).ExecContext(ctx,
		`DELETE FROM mixes WHERE id = ?`,
		id,
	)
	if err != nil {
		return fmt.Errorf("delete mix by id: %w", err)
	}

	return nil
}
