package sqlite

import (
	database "chem-factory/internal/database/sqlite"
	"chem-factory/internal/domain"
	"chem-factory/pkg/reedam"
	"chem-factory/utils/convert"
	"errors"

	"context"
	"database/sql"
)

type MarketRepo struct{ db *database.Database }

func NewMarketRepo(db *database.Database) *MarketRepo { return &MarketRepo{db: db} }

func (r *MarketRepo) Export(ctx context.Context) ([]domain.Market, error) {
	rows, err := r.db.Extract(ctx).QueryContext(ctx,
		`SELECT id, user_id, material_id, amount, date_time FROM market`,
	)
	if err != nil {
		return nil, reedam.Unexpected(err).WithLog(database.IsTx(ctx))
	}
	defer rows.Close()

	var list []domain.Market
	for rows.Next() {
		var (
			market domain.Market
			userID sql.NullInt64
		)
		if err := rows.Scan(
			&market.ID,
			&userID,
			&market.MaterialID,
			&market.Amount,
			&market.DateTime,
		); err != nil {
			return nil, reedam.Unexpected(err).WithLog(database.IsTx(ctx))
		}
		market.UserID = convert.SQLiteNullInt64ToUint(userID)
		list = append(list, market)
	}

	if err := rows.Err(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, reedam.MarketNotFound(err)
		}
		return nil, reedam.Unexpected(err).WithLog(database.IsTx(ctx))
	}
	return list, nil
}

func (r *MarketRepo) FindIDByUserIDMatID(ctx context.Context, userID, materialID uint) (uint, error) {
	var id uint
	if err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT id FROM market WHERE user_id = ? AND material_id = ?`,
		userID, materialID,
	).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, reedam.MarketNotFound(err)
		}
		return 0, reedam.Unexpected(err).WithLog(userID, materialID, database.IsTx(ctx))
	}
	return id, nil
}

func (r *MarketRepo) FindByID(ctx context.Context, id uint) (domain.Market, error) {
	var (
		market domain.Market
		userID sql.NullInt64
	)
	if err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT id, user_id, material_id, amount, date_time
		FROM market WHERE id = ?`,
		id,
	).Scan(
		&market.ID,
		&userID,
		&market.MaterialID,
		&market.Amount,
		&market.DateTime,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Market{}, reedam.MarketNotFound(err)
		}
		return domain.Market{}, reedam.Unexpected(err).WithLog(id, database.IsTx(ctx))
	}
	market.UserID = convert.SQLiteNullInt64ToUint(userID)
	return market, nil
}

func (r *MarketRepo) ReduceAmountByID(ctx context.Context, id uint, amount int) error {
	_, err := r.db.Extract(ctx).ExecContext(ctx,
		`UPDATE market SET amount = amount - ? WHERE id = ?`,
		amount,
		id,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return reedam.MarketNotFound(err)
		}
		return reedam.Unexpected(err).WithLog(id, amount, database.IsTx(ctx))
	}
	return nil
}

func (r *MarketRepo) IncreaseAmountByID(ctx context.Context, id uint, amount int) error {
	_, err := r.db.Extract(ctx).ExecContext(ctx,
		`UPDATE market SET amount = amount + ? WHERE id = ?`,
		amount,
		id,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return reedam.MarketNotFound(err)
		}
		return reedam.Unexpected(err).WithLog(id, amount, database.IsTx(ctx))
	}
	return nil
}

func (r *MarketRepo) Add(ctx context.Context, market domain.Market) error {
	_, err := r.db.Extract(ctx).ExecContext(ctx,
		`INSERT INTO market(user_id, material_id, amount, date_time)
		VALUES (?, ?, ?, ?)`,
		convert.UintToNullInt64(market.UserID),
		market.MaterialID,
		market.Amount,
		market.DateTime,
	)
	if err != nil {
		return reedam.Unexpected(err).WithLog(market, database.IsTx(ctx))
	}
	return nil
}

func (r *MarketRepo) DeleteByID(ctx context.Context, id uint) error {
	_, err := r.db.Extract(ctx).ExecContext(ctx,
		`DELETE FROM market WHERE id = ?`,
		id,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return reedam.MarketNotFound(err)
		}
		return reedam.Unexpected(err).WithLog(id, database.IsTx(ctx))
	}
	return nil
}

func (r *MarketRepo) FindUserIDByID(ctx context.Context, id uint) (uint, error) {
	var userID sql.NullInt64
	if err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT user_id FROM market WHERE id = ?`,
		id,
	).Scan(&userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, reedam.MarketNotFound(err)
		}
		return 0, reedam.Unexpected(err).WithLog(id, database.IsTx(ctx))
	}
	return convert.SQLiteNullInt64ToUint(userID), nil
}

func (r *MarketRepo) FindMatIDByID(ctx context.Context, id uint) (uint, error) {
	var materialID uint
	if err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT material_id FROM market WHERE id = ?`,
		id,
	).Scan(&materialID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, reedam.MarketNotFound(err)
		}
		return 0, reedam.Unexpected(err).WithLog(id, database.IsTx(ctx))
	}
	return materialID, nil
}
