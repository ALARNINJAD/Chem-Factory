package repository

import (
	"chem-factory/internal/model"
	"errors"
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

func createShopTable(r *repositoryManager) {
	query := `
	CREATE TABLE IF NOT EXISTS shop (
		id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
		user_id	INTEGER NOT NULL,
		material_id	INTEGER NOT NULL,
		username TEXT NOT NULL,
		material_name TEXT NOT NULL,
		number INTEGER NOT NULL DEFAULT 0 CHECK(10 >= "number" >= 1),
		price INTEGER NOT NULL CHECK("price" >= 0),
		date_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UNIQUE("user_id","material_id","price"),
		UNIQUE("username","material_name","price"),
		FOREIGN KEY("material_id") REFERENCES "material"("id"),
		FOREIGN KEY("user_id") REFERENCES "user"("id")
	)`
	if _, err := r.db.Exec(query); err != nil {
		panic("Could not create shop table.")
	}
}

func (r *repositoryManager) ExportShop() ([]shop, error) {

	rows, err := r.db.Query("SELECT * FROM shop")
	if err != nil {
		return nil, err
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
			return nil, err
		}

		list = append(list, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *repositoryManager) AddToShop(s model.Shop) error {

	userID, err := r.ExportIDbyUsername(s.Username)
	if err != nil {
		return err
	}

	materialID, err := r.ExportIDbyMaterialName(s.MaterialName)
	if err != nil {
		return err
	}

	shp := shop{
		UserID:       userID,
		MaterialID:   materialID,
		Username:     s.Username,
		MaterialName: s.MaterialName,
		Number:       s.Number,
		Price:        s.Price,
	}

	r.db.QueryRow("SELECT id FROM shop WHERE user_id = ? AND material_id = ? AND price = ?",
		shp.UserID, shp.MaterialID, shp.Price).Scan(&shp.ID)

	if shp.ID != 0 {
		_, err := r.db.Exec("UPDATE shop SET number = number + ? WHERE id = ?", shp.Number, shp.ID)
		if err != nil {
			return err
		}
		return nil
	}

	_, err = r.db.Exec("INSERT INTO shop(user_id,material_id,username,material_name,number,price) VALUES(?,?,?,?,?,?)",
		shp.UserID, shp.MaterialID, shp.Username, shp.MaterialName, shp.Number, shp.Price)
	if err != nil {
		return err
	}

	return nil
}

func (r *repositoryManager) RemoveFromShop(s model.Shop) error {

	shp := shop{
		Username:     s.Username,
		MaterialName: s.MaterialName,
		Number:       s.Number,
		Price:        s.Price,
	}

	err := r.db.QueryRow("SELECT id FROM shop WHERE username = ? AND material_name = ? AND price = ?",
		shp.Username, shp.MaterialName, shp.Price).Scan(&shp.ID)
	if err != nil {
		return err
	}

	var number int
	err = r.db.QueryRow("SELECT number FROM shop WHERE id = ?", shp.ID).Scan(&number)
	if err != nil {
		return err
	}

	if number < shp.Number {

		return errors.New("Could not remove items from shop.")

	} else if number == shp.Number {

		_, err = r.db.Exec("DELETE FROM shop WHERE id = ?", shp.ID)
		if err != nil {
			return err
		}

	} else {

		_, err = r.db.Exec("UPDATE shop SET number = number - ? WHERE id = ?", shp.Number, shp.ID)
		if err != nil {
			return err
		}

	}

	return nil
}
