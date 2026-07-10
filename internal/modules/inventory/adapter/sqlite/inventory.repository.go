package sqlite

import (
	"chem-factory/internal/database/sqlite"
	"chem-factory/internal/domain"
	"context"
	"database/sql"
	"fmt"
)

type InventoryRepo struct{ db *sqlite.Database }

func NewInventoryRepo(db *sqlite.Database) *InventoryRepo { return &InventoryRepo{db: db} }

func (r *InventoryRepo) FindByID(ctx context.Context, id int) (domain.Inventory, error) {
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
		return domain.Inventory{}, fmt.Errorf("find inventory by id: %w", err)
	}

	return inventory, nil
}

func (r *InventoryRepo) FindByUserID(ctx context.Context, userID int) ([]domain.Inventory, error) {
	rows, err := r.db.Extract(ctx).QueryContext(ctx,
		`SELECT id, user_id, material_id, amount, date_time FROM inventory WHERE user_id = ?`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("find inventory by user id select query: %w", err)
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
			return nil, fmt.Errorf("find inventory by user id rows scan: %w", err)
		}
		list = append(list, inventory)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("find inventory by user id rows error: %w", err)
	}

	return list, nil
}

func (r *InventoryRepo) FindIDByUserIDmatID(ctx context.Context, userID, materialID int) (int, error) {
	var id int
	if err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT id FROM inventory WHERE user_id = ? AND material_id = ?`,
		userID, materialID,
	).Scan(&id); err != nil {
		return 0, fmt.Errorf("find inventory id by user and material id: %w", err)
	}
	return id, nil
}

func (r *InventoryRepo) HasByUserIDMaterialID(ctx context.Context, userID, materialID int) (bool, error) {
	var exists int
	if err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT 1 FROM inventory WHERE user_id = ? AND material_id = ? LIMIT 1`,
		userID, materialID,
	).Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("check inventory exists by user and material id: %w", err)
	}
	return true, nil
}

func (r *InventoryRepo) IncreaseByID(ctx context.Context, id, amount int) error {
	_, err := r.db.Extract(ctx).ExecContext(ctx,
		`UPDATE inventory SET amount = amount + ? WHERE id = ?`,
		amount,
		id,
	)
	if err != nil {
		return fmt.Errorf("increase inventory amount by id: %w", err)
	}
	return nil
}

func (r *InventoryRepo) ReduceByID(ctx context.Context, id, amount int) error {
	_, err := r.db.Extract(ctx).ExecContext(ctx,
		`UPDATE inventory SET amount = amount - ? WHERE id = ?`,
		amount,
		id,
	)
	if err != nil {
		return fmt.Errorf("reduce inventory amount by id: %w", err)
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
		return fmt.Errorf("add inventory: %w", err)
	}
	return nil
}

func (r *InventoryRepo) DeleteByID(ctx context.Context, id int) error {
	_, err := r.db.Extract(ctx).ExecContext(ctx,
		`DELETE FROM inventory WHERE id = ?`,
		id,
	)
	if err != nil {
		return fmt.Errorf("delete inventory by id: %w", err)
	}
	return nil
}
