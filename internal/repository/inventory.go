package repository

import (
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

func createInventoryTable(r *repositoryManager) {
	query := `
	CREATE TABLE IF NOT EXISTS inventory (
		id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
		user_id	INTEGER NOT NULL,
		material_id	INTEGER NOT NULL,
		username TEXT NOT NULL,
		material_name TEXT NOT NULL,
		number INTEGER NOT NULL CHECK("number" >= 1),
		date_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UNIQUE("user_id","material_id"),
		UNIQUE("username","material_name"),
		FOREIGN KEY("material_id") REFERENCES "material"("id"),
		FOREIGN KEY("user_id") REFERENCES "user"("id")
	)`
	if _, err := r.db.Exec(query); err != nil {
		panic(fmt.Errorf("Repository inventory, create inventory table: %w", err))
	}
}

func (r *repositoryManager) FindInvenIDByUserIDmatID(userID, materialID int) (int, error) {

	var id int
	err := r.db.QueryRow("SELECT id FROM inventory WHERE user_id = ? AND material_id = ?",
		userID, materialID).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("Repository inventory, find inventory id by user&material id: %w", err)
	}
	return id, nil
}

func (r *repositoryManager) IncreaseInventoryByID(id, number int) error {

	_, err := r.db.Exec("UPDATE inventory SET number = number + ? WHERE id = ?", number, id)
	if err != nil {
		return fmt.Errorf("Repository inventory, increase inventoy number: %w", err)
	}
	return nil
}

func (r *repositoryManager) ReduceInventoryByID(id, number int) error {

	_, err := r.db.Exec("UPDATE inventory SET number = number - ? WHERE id = ?", number, id)
	if err != nil {
		return fmt.Errorf("Repository inventory, reduce inventoy number: %w", err)
	}
	return nil
}

func (r *repositoryManager) AddToInventory(i inventory) error {

	_, err := r.db.Exec(`
		INSERT INTO 
		inventory(user_id,material_id,username,material_name,number)
		VALUES(?,?,?,?,?)`,
		i.UserID, i.MaterialID, i.Username, i.MaterialName, i.Number)
	if err != nil {
		return fmt.Errorf("Repository inventory, add to inventory: %w", err)
	}

	return nil
}

func (r *repositoryManager) DeleteInventoryByID(id int) error {

	_, err := r.db.Exec("DELETE FROM inventory WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("Repository inventory, delelte inventory by id: %w", err)
	}
	return nil
}

func (r *repositoryManager) GetInventoryByUsername(username string) ([]inventory, error) {

	rows, err := r.db.Query(`
        SELECT id, user_id, material_id, username, material_name, number, date_time
        FROM inventory 
        WHERE username = ?`, username)
	if err != nil {
		return nil, fmt.Errorf("Repository inventory, get inventory by username query: %w", err)
	}
	defer rows.Close()

	var list []inventory

	for rows.Next() {
		var i inventory

		err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.MaterialID,
			&i.Username,
			&i.MaterialName,
			&i.Number,
			&i.DateTime,
		)

		if err != nil {
			return nil, fmt.Errorf("Repository inventory, get inventory by username scan: %w", err)
		}

		list = append(list, i)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Repository inventory, get inventory by username rows: %w", err)
	}

	return list, nil
}
