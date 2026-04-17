package repository

import "time"

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
		UNIQUE("user_id","material_id"),
		UNIQUE("username","material_name"),
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
