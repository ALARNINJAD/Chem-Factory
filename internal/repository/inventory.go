package repository

import (
	"database/sql"
	"fmt"
	"time"
)

type inventory struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	MaterialID   int       `json:"material_id"`
	Username     string    `json:"username"`
	MaterialName string    `json:"material_name"`
	Number       int       `json:"number"`
	DateTime     time.Time `json:"date_time"`
}

type inventoryManager struct{ db *sql.DB }

func NewInventoryManager(db *sql.DB) *inventoryManager { return &inventoryManager{db: db} }

func (i *inventoryManager) EmptyStruct() *inventory {	return &inventory{} }

func (i *inventoryManager) EmptySlice() []inventory { return []inventory{} }

func (i *inventoryManager) FindByID(id int) (*inventory, error) {
	var inv inventory
	err := i.db.QueryRow(`
		SELECT id, user_id, material_id, username, material_name, number, date_time
		FROM inventory WHERE id = ?`, id).Scan(
		&inv.ID, &inv.UserID, &inv.MaterialID, &inv.Username, &inv.MaterialName, &inv.Number, &inv.DateTime)
	if err != nil {
		return nil, fmt.Errorf("Repository inventory, find user inventory by id: %w", err)
	}
	return &inv, nil
}

func (i *inventoryManager) FindIDByUserIDmatID(userID, materialID int) (int, error) {
	var id int
	err := i.db.QueryRow("SELECT id FROM inventory WHERE user_id = ? AND material_id = ?",
		userID, materialID).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("Repository inventory, find inventory id by user and material id: %w", err)
	}
	return id, nil
}

func (i *inventoryManager) IncreaseByID(tx *sql.Tx, id, number int) error {
	_, err := tx.Exec("UPDATE inventory SET number = number + ? WHERE id = ?", number, id)
	if err != nil {
		return fmt.Errorf("Repository inventory, increase inventoy number: %w", err)
	}
	return nil
}

func (i *inventoryManager) ReduceByID(tx *sql.Tx, id, number int) error {
	_, err := tx.Exec("UPDATE inventory SET number = number - ? WHERE id = ?", number, id)
	if err != nil {
		return fmt.Errorf("Repository inventory, reduce inventoy number: %w", err)
	}
	return nil
}

func (i *inventoryManager) Add(tx *sql.Tx, inv inventory) error {
	_, err := tx.Exec(`
		INSERT INTO 
		inventory(user_id,material_id,username,material_name,number)
		VALUES(?,?,?,?,?)`,
		inv.UserID, inv.MaterialID, inv.Username, inv.MaterialName, inv.Number)
	if err != nil {
		return fmt.Errorf("Repository inventory, add to inventory: %w ", err)
	}
	return nil
}

func (i *inventoryManager) DeleteByID(tx *sql.Tx, id int) error {
	_, err := tx.Exec("DELETE FROM inventory WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("Repository inventory, delelte inventory by id: %w", err)
	}
	return nil
}

func (i *inventoryManager) FindByUsername(username string) ([]inventory, error) {
	rows, err := i.db.Query(`
        SELECT id, user_id, material_id, username, material_name, number, date_time
        FROM inventory 
        WHERE username = ?`, username)
	if err != nil {
		return nil, fmt.Errorf("Repository inventory, get inventory by username query: %w", err)
	}
	defer rows.Close()

	var list []inventory

	for rows.Next() {
		var inv inventory
		err := rows.Scan(
			&inv.ID,
			&inv.UserID,
			&inv.MaterialID,
			&inv.Username,
			&inv.MaterialName,
			&inv.Number,
			&inv.DateTime,
		)
		if err != nil {
			return nil, fmt.Errorf("Repository inventory, get inventory by username scan: %w", err)
		}
		list = append(list, inv)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Repository inventory, get inventory by username rows: %w", err)
	}
	return list, nil
}
