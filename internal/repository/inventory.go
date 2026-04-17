package repository

import (
	"chem-factory/internal/model"
	"errors"
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
		panic("Could not create inventory table.")
	}
}

func (r *repositoryManager) AddToInventory(i model.Inventory) error {

	userID, err := r.ExportIDbyUsername(i.Username)
	if err != nil {
		return err
	}

	materialID, err := r.ExportIDbyMaterialName(i.MaterialName)
	if err != nil {
		return err
	}

	inv := inventory{
		UserID:       userID,
		MaterialID:   materialID,
		Username:     i.Username,
		MaterialName: i.MaterialName,
		Number:       i.Number,
	}

	r.db.QueryRow("SELECT id FROM inventory WHERE username = ? AND material_name = ?", inv.Username, inv.MaterialName).Scan(&inv.ID)

	if inv.ID != 0 {
		_, err := r.db.Exec("UPDATE inventory SET number = number + ? WHERE id = ?", inv.Number, inv.ID)
		if err != nil {
			return err
		}
		return nil
	}

	_, err = r.db.Exec("INSERT INTO inventory(user_id,material_id,username,material_name,number) VALUES(?,?,?,?,?)",
		inv.UserID, inv.MaterialID, inv.Username, inv.MaterialName, inv.Number)
	if err != nil {
		return err
	}

	return nil
}

func (r *repositoryManager) RemoveFromInventory(i model.Inventory) error {

	inv := inventory{
		Username:     i.Username,
		MaterialName: i.MaterialName,
		Number:       i.Number,
	}

	err := r.db.QueryRow("SELECT id FROM inventory WHERE username = ? AND material_name = ?", inv.Username, inv.MaterialName).Scan(&inv.ID)
	if err != nil {
		return err
	}

	var number int
	err = r.db.QueryRow("SELECT number FROM inventory WHERE id = ?", inv.ID).Scan(&number)
	if err != nil {
		return err
	}

	if number < inv.Number {

		return errors.New("Could not remove items from inventory.")

	} else if number == inv.Number {

		_, err = r.db.Exec("DELETE FROM inventory WHERE id = ?", inv.ID)
		if err != nil {
			return err
		}

	} else {

		_, err = r.db.Exec("UPDATE inventory SET number = number - ? WHERE id = ?", inv.Number, inv.ID)

	}

	return nil
}

func (r *repositoryManager) ExportInventory(username string) ([]inventory, error) {

	rows, err := r.db.Query(`
        SELECT id, user_id, material_id, username, material_name, number, date_time
        FROM inventory 
        WHERE username = ?`, username)
	if err != nil {
		return nil, err
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
			return nil, err
		}

		list = append(list, i)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}
