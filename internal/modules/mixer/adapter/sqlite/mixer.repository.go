package sqlite

import (
	database "chem-factory/internal/database/sqlite"
	"chem-factory/internal/domain"
	"context"
	"fmt"
)

type MixerRepo struct{ db *database.Database }

func NewMixerRepo(db *database.Database) *MixerRepo { return &MixerRepo{db: db} }

func (r *MixerRepo) FindByID(ctx context.Context, id uint) (domain.Mixer, error) {
	var mixer domain.Mixer

	err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT id, user_id, first_ingredient_id, second_ingredient_id, amount, date_time FROM mixer WHERE id = ?`,
		id,
	).Scan(
		&mixer.ID,
		&mixer.UserID,
		&mixer.FirstIngredientID,
		&mixer.SecondIngredientID,
		&mixer.Amount,
		&mixer.DateTime,
	)
	if err != nil {
		return domain.Mixer{}, fmt.Errorf("find mixer by id: %w", err)
	}

	return mixer, nil
}

func (r *MixerRepo) Add(ctx context.Context, mixer domain.Mixer) error {
	_, err := r.db.Extract(ctx).ExecContext(ctx,
		`INSERT INTO mixer(user_id, first_ingredient_id, second_ingredient_id, amount, date_time)
		VALUES (?, ?, ?, ?, ?)`,
		mixer.UserID,
		mixer.FirstIngredientID,
		mixer.SecondIngredientID,
		mixer.Amount,
		mixer.DateTime,
	)
	if err != nil {
		return fmt.Errorf("add mixer: %w", err)
	}

	return nil
}

func (r *MixerRepo) FindIDByUserIDIngrID(ctx context.Context, userID, firstID, secID uint) (uint, error) {
	var id uint

	err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT id FROM mixer WHERE user_id = ? AND first_ingredient_id = ? AND second_ingredient_id = ?`,
		userID,
		firstID,
		secID,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("find mixer id by user id and ingredients id: %w", err)
	}

	return id, nil
}

func (r *MixerRepo) DeleteByID(ctx context.Context, id uint) error {
	_, err := r.db.Extract(ctx).ExecContext(ctx,
		`DELETE FROM mixer WHERE id = ?`,
		id,
	)
	if err != nil {
		return fmt.Errorf("delete mixer by id: %w", err)
	}

	return nil
}
