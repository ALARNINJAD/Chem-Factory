package sqlite

import (
	"chem-factory/internal/database/sqlite"
	"chem-factory/internal/domain"
	"context"
	"fmt"
)

type MarketRepo struct{ db *sqlite.Database }

func NewMarketRepo(db *sqlite.Database) *MarketRepo { return &MarketRepo{db: db} }

func (r *MarketRepo) Export(ctx context.Context) ([]domain.Market, error) {
	rows, err := r.db.Extract(ctx).QueryContext(ctx,
		`SELECT id, user_id, material_id, amount, price, date_time FROM market`,
	)
	if err != nil {
		return nil, fmt.Errorf("export market select all query: %w", err)
	}
	defer rows.Close()

	var list []domain.Market

	for rows.Next() {
		var market domain.Market
		if err := rows.Scan(
			&market.ID,
			&market.UserID,
			&market.MaterialID,
			&market.Amount,
			&market.Price,
			&market.DateTime,
		); err != nil {
			return nil, fmt.Errorf("export market rows scan: %w", err)
		}
		list = append(list, market)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("export market rows error: %w", err)
	}
	return list, nil
}

func (r *MarketRepo) FindIDByPriceUserIDMatID(ctx context.Context, price, userID, materialID int) (int, error) {
	var id int
	if err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT id FROM market WHERE price = ? AND user_id = ? AND material_id = ?`,
		price, userID, materialID,
	).Scan(&id); err != nil {
		return 0, fmt.Errorf("find market id by price, user id and material id: %w", err)
	}
	return id, nil
}

func (r *MarketRepo) FindByPriceUserIDMatID(ctx context.Context, price, userID, materialID int) (domain.Market, error) {
	var market domain.Market
	if err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT id, user_id, material_id, amount, price, date_time
		FROM market WHERE price = ? AND user_id = ? AND material_id = ?`,
		price, userID, materialID,
	).Scan(
		&market.ID,
		&market.UserID,
		&market.MaterialID,
		&market.Amount,
		&market.Price,
		&market.DateTime,
	); err != nil {
		return domain.Market{}, fmt.Errorf("find market by price, user id and material id: %w", err)
	}
	return market, nil
}

func (r *MarketRepo) FindByID(ctx context.Context, id int) (domain.Market, error) {
	var market domain.Market
	if err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT id, user_id, material_id, amount, price, date_time
		FROM market WHERE id = ?`,
		id,
	).Scan(
		&market.ID,
		&market.UserID,
		&market.MaterialID,
		&market.Amount,
		&market.Price,
		&market.DateTime,
	); err != nil {
		return domain.Market{}, fmt.Errorf("find market by id: %w", err)
	}
	return market, nil
}

func (r *MarketRepo) ReduceAmountByID(ctx context.Context, id, amount int) error {
	_, err := r.db.Extract(ctx).ExecContext(ctx,
		`UPDATE market SET amount = amount - ? WHERE id = ?`,
		amount,
		id,
	)
	if err != nil {
		return fmt.Errorf("reduce market amount by id: %w", err)
	}
	return nil
}

func (r *MarketRepo) IncreaseAmountByID(ctx context.Context, id, amount int) error {
	_, err := r.db.Extract(ctx).ExecContext(ctx,
		`UPDATE market SET amount = amount + ? WHERE id = ?`,
		amount,
		id,
	)
	if err != nil {
		return fmt.Errorf("increase market amount by id: %w", err)
	}
	return nil
}

func (r *MarketRepo) Add(ctx context.Context, market domain.Market) error {
	_, err := r.db.Extract(ctx).ExecContext(ctx,
		`INSERT INTO market(user_id, material_id, amount, price, date_time)
		VALUES (?, ?, ?, ?, ?)`,
		market.UserID,
		market.MaterialID,
		market.Amount,
		market.Price,
		market.DateTime,
	)
	if err != nil {
		return fmt.Errorf("add market: %w", err)
	}
	return nil
}

func (r *MarketRepo) DeleteByID(ctx context.Context, id int) error {
	_, err := r.db.Extract(ctx).ExecContext(ctx,
		`DELETE FROM market WHERE id = ?`,
		id,
	)
	if err != nil {
		return fmt.Errorf("delete market by id: %w", err)
	}
	return nil
}

func (r *MarketRepo) FindUserIDByID(ctx context.Context, id int) (int, error) {
	var userID int
	if err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT user_id FROM market WHERE id = ?`,
		id,
	).Scan(&userID); err != nil {
		return 0, fmt.Errorf("find user id by market id: %w", err)
	}
	return userID, nil
}

func (r *MarketRepo) FindMatIDByID(ctx context.Context, id int) (int, error) {
	var materialID int
	if err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT material_id FROM market WHERE id = ?`,
		id,
	).Scan(&materialID); err != nil {
		return 0, fmt.Errorf("find material id by market id: %w", err)
	}
	return materialID, nil
}
