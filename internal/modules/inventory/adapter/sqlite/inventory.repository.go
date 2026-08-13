package sqlite

import (
	database "chem-factory/internal/database/sqlite"
	"chem-factory/internal/domain"
	"chem-factory/pkg/reedam"
	"context"
	"database/sql"
	"errors"
)

type InventoryRepo struct{ db *database.Database }

func NewInventoryRepo(db *database.Database) *InventoryRepo { return &InventoryRepo{db: db} }

func (r *InventoryRepo) FindByID(ctx context.Context, id uint) (domain.Inventory, error) {
	var inventory domain.Inventory

	err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT id, user_id, material_id, amount, date_time FROM inventory WHERE id = ?`,
		id,
	).Scan(
		&inventory.ID,
		&inventory.UserID,
		&inventory.MaterialID,
		&inventory.Amount,
		&inventory.DateTime,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Inventory{}, reedam.InventoryNotFound(err)
		}
		return domain.Inventory{}, reedam.Unexpected(err).WithLog(id, database.IsTx(ctx))
	}

	return inventory, nil
}

func (r *InventoryRepo) FindByUserID(ctx context.Context, userID uint) ([]domain.Inventory, error) {
	rows, err := r.db.Extract(ctx).QueryContext(ctx,
		`SELECT id, user_id, material_id, amount, date_time FROM inventory WHERE user_id = ?`,
		userID,
	)
	if err != nil {
		return nil, reedam.Unexpected(err).WithLog(userID, database.IsTx(ctx))
	}
	defer rows.Close()

	var list []domain.Inventory

	for rows.Next() {
		var inventory domain.Inventory
		if err := rows.Scan(
			&inventory.ID,
			&inventory.UserID,
			&inventory.MaterialID,
			&inventory.Amount,
			&inventory.DateTime,
		); err != nil {
			return nil, reedam.Unexpected(err).WithLog(userID, database.IsTx(ctx))
		}
		list = append(list, inventory)
	}

	if err := rows.Err(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, reedam.InventoryNotFound(err)
		}
		return nil, reedam.Unexpected(err).WithLog(userID, database.IsTx(ctx))
	}

	return list, nil
}

func (r *InventoryRepo) FindByUserIDmatID(ctx context.Context, userID, materialID uint) (domain.Inventory, error) {
	var inventory domain.Inventory

	err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT id, user_id, material_id, amount, date_time FROM inventory WHERE user_id = ? AND material_id = ?`,
		userID, materialID,
	).Scan(
		&inventory.ID,
		&inventory.UserID,
		&inventory.MaterialID,
		&inventory.Amount,
		&inventory.DateTime,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Inventory{}, reedam.InventoryNotFound(err)
		}
		return domain.Inventory{}, reedam.Unexpected(err).WithLog(userID, materialID, database.IsTx(ctx))
	}

	return inventory, nil
}

func (r *InventoryRepo) FindIDByUserIDmatID(ctx context.Context, userID, materialID uint) (uint, error) {
	var id uint
	if err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT id FROM inventory WHERE user_id = ? AND material_id = ?`,
		userID, materialID,
	).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, reedam.InventoryNotFound(err)
		}
		return 0, reedam.Unexpected(err).WithLog(userID, materialID, database.IsTx(ctx))
	}
	return id, nil
}

func (r *InventoryRepo) IncreaseByID(ctx context.Context, id uint, amount int) error {
	_, err := r.db.Extract(ctx).ExecContext(ctx,
		`UPDATE inventory SET amount = amount + ? WHERE id = ?`,
		amount,
		id,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return reedam.InventoryNotFound(err)
		}
		return reedam.Unexpected(err).WithLog(id, amount, database.IsTx(ctx))
	}
	return nil
}

func (r *InventoryRepo) ReduceByID(ctx context.Context, id uint, amount int) error {
	_, err := r.db.Extract(ctx).ExecContext(ctx,
		`UPDATE inventory SET amount = amount - ? WHERE id = ?`,
		amount,
		id,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return reedam.InventoryNotFound(err)
		}
		return reedam.Unexpected(err).WithLog(id, amount, database.IsTx(ctx))
	}
	return nil
}

func (r *InventoryRepo) Add(ctx context.Context, inventory domain.Inventory) error {
	_, err := r.db.Extract(ctx).ExecContext(ctx,
		`INSERT INTO inventory(user_id, material_id, amount, date_time)
		VALUES (?, ?, ?, ?)`,
		inventory.UserID,
		inventory.MaterialID,
		inventory.Amount,
		inventory.DateTime,
	)
	if err != nil {
		return reedam.Unexpected(err).WithLog(inventory, database.IsTx(ctx))
	}
	return nil
}

func (r *InventoryRepo) DeleteByID(ctx context.Context, id uint) error {
	_, err := r.db.Extract(ctx).ExecContext(ctx,
		`DELETE FROM inventory WHERE id = ?`,
		id,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return reedam.InventoryNotFound(err)
		}
		return reedam.Unexpected(err).WithLog(id, database.IsTx(ctx))
	}
	return nil
}
