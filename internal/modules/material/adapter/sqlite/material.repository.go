package sqlite

import (
	database "chem-factory/internal/database/sqlite"
	"chem-factory/internal/domain"
	"chem-factory/pkg/lang"
	"chem-factory/pkg/reedam"
	"chem-factory/utils/convert"
	"context"
	"database/sql"
)

type MaterialRepo struct{ db *database.Database }

func NewMaterialRepo(db *database.Database) *MaterialRepo { return &MaterialRepo{db: db} }

func (r *MaterialRepo) FindMixTimeByID(ctx context.Context, id uint) (int, error) {
	var mixTime sql.NullInt64
	if err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT mix_time FROM materials WHERE id = ?`, id,
	).Scan(&mixTime); err != nil {
		return 0, reedam.New().WithError(err).WithMessage(lang.ErrorUnexpected).WithStatus(reedam.StatusInternalServerError).WithLog(id, database.IsTx(ctx))
	}
	return convert.SQLiteNullInt64ToInt(mixTime), nil
}

func (r *MaterialRepo) FindIDByIngrID(ctx context.Context, firstID uint, secondID uint) (uint, error) {
	var id uint
	if err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT id FROM materials WHERE first_ingredient_id = ? AND second_ingredient_id = ?`,
		firstID, secondID,
	).Scan(&id); err != nil {
		return 0, reedam.New().WithError(err).WithMessage(lang.ErrorUnexpected).WithStatus(reedam.StatusInternalServerError).WithLog(firstID, secondID, database.IsTx(ctx))
	}
	return id, nil
}

func (r *MaterialRepo) FindByIngrID(ctx context.Context, firstID uint, secondID uint) (domain.Material, error) {
	var (
		mat                                        domain.Material
		userID, firstIngrID, secondIngrID, mixTime sql.NullInt64
	)
	if err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT id, user_id, first_ingredient_id, second_ingredient_id,
		name, price, mix_time
		FROM materials WHERE first_ingredient_id = ? AND second_ingredient_id = ?`,
		firstID, secondID,
	).Scan(
		&mat.ID,
		&userID, &firstIngrID, &secondIngrID,
		&mat.Name, &mat.Price, &mixTime,
	); err != nil {
		return domain.Material{}, reedam.New().WithError(err).WithMessage(lang.ErrorUnexpected).WithStatus(reedam.StatusInternalServerError).WithLog(firstID, secondID, database.IsTx(ctx))
	}
	mat.UserID = convert.SQLiteNullInt64ToUint(userID)
	mat.FirstIngredientID = convert.SQLiteNullInt64ToUint(firstIngrID)
	mat.SecondIngredientID = convert.SQLiteNullInt64ToUint(secondIngrID)
	mat.MixTime = convert.SQLiteNullInt64ToInt(mixTime)
	return mat, nil
}

func (r *MaterialRepo) FindNameByIngrID(ctx context.Context, firstID uint, secondID uint) (string, error) {
	var name string
	if err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT name FROM materials WHERE first_ingredient_id = ? AND second_ingredient_id = ?`,
		firstID, secondID,
	).Scan(&name); err != nil {
		return "", reedam.New().WithError(err).WithMessage(lang.ErrorUnexpected).WithStatus(reedam.StatusInternalServerError).WithLog(firstID, secondID, database.IsTx(ctx))
	}
	return name, nil
}

func (r *MaterialRepo) FindByID(ctx context.Context, id uint) (domain.Material, error) {
	var mat domain.Material
	var (
		userID, firstIngrID, secondIngrID, mixTime sql.NullInt64
	)
	if err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT id, user_id, first_ingredient_id, second_ingredient_id,
		name, price, mix_time
		FROM materials WHERE id = ?`,
		id,
	).Scan(
		&mat.ID,
		&userID, &firstIngrID, &secondIngrID,
		&mat.Name, &mat.Price, &mixTime,
	); err != nil {
		return domain.Material{}, reedam.New().WithError(err).WithMessage(lang.ErrorUnexpected).WithStatus(reedam.StatusInternalServerError).WithLog(id, database.IsTx(ctx))
	}
	mat.UserID = convert.SQLiteNullInt64ToUint(userID)
	mat.FirstIngredientID = convert.SQLiteNullInt64ToUint(firstIngrID)
	mat.SecondIngredientID = convert.SQLiteNullInt64ToUint(secondIngrID)
	mat.MixTime = convert.SQLiteNullInt64ToInt(mixTime)
	return mat, nil
}

func (r *MaterialRepo) FindIDByName(ctx context.Context, name string) (uint, error) {
	var id uint
	if err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT id FROM materials WHERE name = ?`, name,
	).Scan(&id); err != nil {
		return 0, reedam.New().WithError(err).WithMessage(lang.ErrorUnexpected).WithStatus(reedam.StatusInternalServerError).WithLog(name, database.IsTx(ctx))
	}
	return id, nil
}

func (r *MaterialRepo) FindNameByID(ctx context.Context, id uint) (string, error) {
	var name string
	if err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT name FROM materials WHERE id = ?`, id,
	).Scan(&name); err != nil {
		return "", reedam.New().WithError(err).WithMessage(lang.ErrorUnexpected).WithStatus(reedam.StatusInternalServerError).WithLog(id, database.IsTx(ctx))
	}
	return name, nil
}

func (r *MaterialRepo) Add(ctx context.Context, material domain.Material) error {
	_, err := r.db.Extract(ctx).ExecContext(ctx,
		`INSERT INTO materials(user_id, first_ingredient_id, second_ingredient_id,
		name, price, mix_time)
		VALUES (?, ?, ?, ?, ?, ?)`,
		convert.UintToNullInt64(material.UserID), convert.UintToNullInt64(material.FirstIngredientID), convert.UintToNullInt64(material.SecondIngredientID),
		material.Name, material.Price, convert.IntToNullInt64(material.MixTime),
	)
	if err != nil {
		return reedam.New().WithError(err).WithMessage(lang.ErrorUnexpected).WithStatus(reedam.StatusInternalServerError).WithLog(material, database.IsTx(ctx))
	}
	return nil
}

func (r *MaterialRepo) FindByUserID(ctx context.Context, id uint) ([]domain.Material, error) {
	rows, err := r.db.Extract(ctx).QueryContext(ctx,
		`SELECT id, user_id, first_ingredient_id, second_ingredient_id,
		name, price, mix_time
		FROM materials WHERE user_id = ?`, id)
	if err != nil {
		return nil, reedam.New().WithError(err).WithMessage(lang.ErrorUnexpected).WithStatus(reedam.StatusInternalServerError).WithLog(id, database.IsTx(ctx))
	}
	defer rows.Close()

	var list []domain.Material

	for rows.Next() {
		var (
			mat                                        domain.Material
			userID, firstIngrID, secondIngrID, mixTime sql.NullInt64
		)
		if err := rows.Scan(
			&mat.ID,
			&userID, &firstIngrID, &secondIngrID,
			&mat.Name, &mat.Price, &mixTime,
		); err != nil {
			return nil, reedam.New().WithError(err).WithMessage(lang.ErrorUnexpected).WithStatus(reedam.StatusInternalServerError).WithLog(id, database.IsTx(ctx))
		}
		mat.UserID = convert.SQLiteNullInt64ToUint(userID)
		mat.FirstIngredientID = convert.SQLiteNullInt64ToUint(firstIngrID)
		mat.SecondIngredientID = convert.SQLiteNullInt64ToUint(secondIngrID)
		mat.MixTime = convert.SQLiteNullInt64ToInt(mixTime)
		list = append(list, mat)
	}

	if err := rows.Err(); err != nil {
		return nil, reedam.New().WithError(err).WithMessage(lang.ErrorUnexpected).WithStatus(reedam.StatusInternalServerError).WithLog(id, database.IsTx(ctx))
	}
	return list, nil
}

func (r *MaterialRepo) Get(ctx context.Context) ([]domain.Material, error) {
	rows, err := r.db.Extract(ctx).QueryContext(ctx,
		`SELECT id, user_id, first_ingredient_id, second_ingredient_id,
		name, price, mix_time
		FROM materials`)
	if err != nil {
		return nil, reedam.New().WithError(err).WithMessage(lang.ErrorUnexpected).WithStatus(reedam.StatusInternalServerError).WithLog(database.IsTx(ctx))
	}
	defer rows.Close()

	var list []domain.Material

	for rows.Next() {
		var (
			mat                                        domain.Material
			userID, firstIngrID, secondIngrID, mixTime sql.NullInt64
		)
		if err := rows.Scan(
			&mat.ID,
			&userID, &firstIngrID, &secondIngrID,
			&mat.Name, &mat.Price, &mixTime,
		); err != nil {
			return nil, reedam.New().WithError(err).WithMessage(lang.ErrorUnexpected).WithStatus(reedam.StatusInternalServerError).WithLog(database.IsTx(ctx))
		}
		mat.UserID = convert.SQLiteNullInt64ToUint(userID)
		mat.FirstIngredientID = convert.SQLiteNullInt64ToUint(firstIngrID)
		mat.SecondIngredientID = convert.SQLiteNullInt64ToUint(secondIngrID)
		mat.MixTime = convert.SQLiteNullInt64ToInt(mixTime)
		list = append(list, mat)
	}

	if err := rows.Err(); err != nil {
		return nil, reedam.New().WithError(err).WithMessage(lang.ErrorUnexpected).WithStatus(reedam.StatusInternalServerError).WithLog(database.IsTx(ctx))
	}
	return list, nil
}

func (r *MaterialRepo) FindPriceByID(ctx context.Context, id uint) (int, error) {
	var price int
	if err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT price FROM materials WHERE id = ?`, id,
	).Scan(&price); err != nil {
		return 0, reedam.New().WithError(err).WithMessage(lang.ErrorUnexpected).WithStatus(reedam.StatusInternalServerError).WithLog(id, database.IsTx(ctx))
	}
	return price, nil
}
