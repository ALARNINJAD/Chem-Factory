package sqlite

import (
	"chem-factory/internal/domain"
	"chem-factory/internal/database/sqlite"
	"context"
	"database/sql"
	"fmt"
)

type MaterialRepo struct{ db *sqlite.Database }

func NewMaterialRepo(db *sqlite.Database) *MaterialRepo { return &MaterialRepo{db: db} }

func nullInt(v *int) any {
	if v == nil {
		return nil
	}
	return *v
}

func toIntPtr(v sql.NullInt64) *int {
	if !v.Valid {
		return nil
	}
	val := int(v.Int64)
	return &val
}

func (r *MaterialRepo) FindMixTimeByID(ctx context.Context, id int) (int, error) {
	var mixTime int
	if err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT mix_time FROM material WHERE id = ?`, id,
	).Scan(&mixTime); err != nil {
		return 0, fmt.Errorf("find material mix time by id: %w", err)
	}
	return mixTime, nil
}

func (r *MaterialRepo) FindIDByIngrID(ctx context.Context, firstID int, secondID int) (int, error) {
	var id int
	if err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT id FROM material WHERE first_ingredient_id = ? AND second_ingredient_id = ?`,
		firstID, secondID,
	).Scan(&id); err != nil {
		return 0, fmt.Errorf("find material id by ingredient id: %w", err)
	}
	return id, nil
}

func (r *MaterialRepo) FindByIngrID(ctx context.Context, firstID int, secondID int) (domain.Material, error) {
	var mat domain.Material
	var (
		userID, firstIngrID, secondIngrID sql.NullInt64
	)
	if err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT id, user_id, first_ingredient_id, second_ingredient_id,
		name, price, mix_time
		FROM material WHERE first_ingredient_id = ? AND second_ingredient_id = ?`,
		firstID, secondID,
	).Scan(
		&mat.ID,
		&userID, &firstIngrID, &secondIngrID,
		&mat.Name, &mat.Price, &mat.MixTime,
	); err != nil {
		return domain.Material{}, fmt.Errorf("find material by ingredient id: %w", err)
	}
	mat.UserID = toIntPtr(userID)
	mat.FirstIngredientID = toIntPtr(firstIngrID)
	mat.SecondIngredientID = toIntPtr(secondIngrID)
	return mat, nil
}

func (r *MaterialRepo) FindByID(ctx context.Context, id int) (domain.Material, error) {
	var mat domain.Material
	var (
		userID, firstIngrID, secondIngrID sql.NullInt64
	)
	if err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT id, user_id, first_ingredient_id, second_ingredient_id,
		name, price, mix_time
		FROM material WHERE id = ?`,
		id,
	).Scan(
		&mat.ID,
		&userID, &firstIngrID, &secondIngrID,
		&mat.Name, &mat.Price, &mat.MixTime,
	); err != nil {
		return domain.Material{}, fmt.Errorf("find material by id: %w", err)
	}
	mat.UserID = toIntPtr(userID)
	mat.FirstIngredientID = toIntPtr(firstIngrID)
	mat.SecondIngredientID = toIntPtr(secondIngrID)
	return mat, nil
}

func (r *MaterialRepo) FindIDByName(ctx context.Context, name string) (int, error) {
	var id int
	if err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT id FROM material WHERE name = ?`, name,
	).Scan(&id); err != nil {
		return 0, fmt.Errorf("find id by material name: %w", err)
	}
	return id, nil
}

func (r *MaterialRepo) Add(ctx context.Context, material domain.Material) error {
	_, err := r.db.Extract(ctx).ExecContext(ctx,
		`INSERT INTO material(user_id, first_ingredient_id, second_ingredient_id,
		name, price, mix_time)
		VALUES (?, ?, ?, ?, ?, ?)`,
		nullInt(material.UserID), nullInt(material.FirstIngredientID), nullInt(material.SecondIngredientID),
		material.Name, material.Price, nullInt(&material.MixTime),
	)
	if err != nil {
		return fmt.Errorf("save material: %w", err)
	}
	return nil
}

func (r *MaterialRepo) FindByUserID(ctx context.Context, id int) ([]domain.Material, error) {
	rows, err := r.db.Extract(ctx).QueryContext(ctx,
		`SELECT id, user_id, first_ingredient_id, second_ingredient_id,
		name, price, mix_time
		FROM material WHERE user_id = ?`, id)
	if err != nil {
		return nil, fmt.Errorf("find materials by user id select query: %w", err)
	}
	defer rows.Close()

	var list []domain.Material

	for rows.Next() {
		var mat domain.Material
		var (
			userID, firstIngrID, secondIngrID sql.NullInt64
		)
		if err := rows.Scan(
			&mat.ID,
			&userID, &firstIngrID, &secondIngrID,
			&mat.Name, &mat.Price, &mat.MixTime,
		); err != nil {
			return nil, fmt.Errorf("find materials by user id rows scan: %w", err)
		}
		mat.UserID = toIntPtr(userID)
		mat.FirstIngredientID = toIntPtr(firstIngrID)
		mat.SecondIngredientID = toIntPtr(secondIngrID)
		list = append(list, mat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("find materials by user id rows error: %w", err)
	}
	return list, nil
}
