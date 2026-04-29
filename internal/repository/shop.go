package repository

import (
	"database/sql"
	"fmt"
	"time"
)

type shop struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	MaterialID   int       `json:"material_id"`
	Username     string    `json:"username"`
	MaterialName string    `json:"material_name"`
	Number       int       `json:"number"`
	Price        int       `json:"price"`
	DateTime     time.Time `json:"date_time"`
}

func (r *repositoryManager) ExportShop() ([]shop, error) {

	rows, err := r.db.Query("SELECT * FROM shop")
	if err != nil {
		return nil, fmt.Errorf("Repository shop, export shop select all query: %w ", err)
	}
	defer rows.Close()

	var list []shop

	for rows.Next() {
		var s shop

		err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.MaterialID,
			&s.Username,
			&s.MaterialName,
			&s.Number,
			&s.Price,
			&s.DateTime,
		)

		if err != nil {
			return nil, fmt.Errorf("Repository shop, export shop rows scan: %w ", err)
		}

		list = append(list, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Repository shop, export shop rows error: %w ", err)
	}

	return list, nil
}

func (r *repositoryManager) EmptyShopStruct() *shop {
	return &shop{}
}

func (r *repositoryManager) FindShopIDByInfo(userID, materialID, price int) (int, error) {

	var id int

	err := r.db.QueryRow("SELECT id FROM shop WHERE user_id = ? AND material_id = ? AND price = ?",
		userID, materialID, price).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("Repository shop, find shop id by info: %w ", err)
	}

	return id, nil
}

func (r *repositoryManager) FindShopByInfo(userID, materialID, price int) (*shop, error) {

	var shp shop

	err := r.db.QueryRow(`
		SELECT id, user_id, material_id, username, material_name, number, price, date_time
		FROM shop WHERE user_id = ? AND material_id = ? AND price = ?`,
		userID, materialID, price).Scan(
		&shp.ID, &shp.UserID, &shp.MaterialID,
		&shp.Username, &shp.MaterialName, &shp.Number, &shp.Price, &shp.DateTime)
	if err != nil {
		return nil, fmt.Errorf("Repository shop, find shop by info: %w ", err)
	}

	return &shp, nil
}

func (r *repositoryManager) ReduceShopNumberByID(tx *sql.Tx, id, number int) error {

	_, err := tx.Exec("UPDATE shop SET number = number - ? WHERE id = ?", number, id)
	if err != nil {
		return fmt.Errorf("Repository shop, reduce shop number by id: %w ", err)
	}
	return nil
}

func (r *repositoryManager) IncreaseShopNumberByID(tx *sql.Tx, id, number int) error {

	_, err := tx.Exec("UPDATE shop SET number = number + ? WHERE id = ?", number, id)
	if err != nil {
		return fmt.Errorf("Repository shop, increase shop number by id: %w ", err)
	}
	return nil
}

func (r *repositoryManager) AddToShop(tx *sql.Tx, s shop) error {

	_, err := tx.Exec(`
		INSERT INTO shop(user_id, material_id, username, material_name, number, price)
		VALUES (?, ?, ?, ?, ?, ?)`,
		s.UserID, s.MaterialID, s.Username, s.MaterialName, s.Number, s.Price)
	if err != nil {
		return fmt.Errorf("Repository shop, add to shop: %w ", err)
	}

	return nil
}

func (r *repositoryManager) DeleteFromShopByID(tx *sql.Tx, id int) error {

	_, err := tx.Exec("DELETE FROM shop WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("Repository shop, delete from shop by id: %w ", err)
	}
	return nil
}
