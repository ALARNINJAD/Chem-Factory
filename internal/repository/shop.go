package repository

import "time"

type shop struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	MaterialID   int       `json:"material_id"`
	Username     string    `json:"username"`
	MaterialName string    `json:"material_name"`
	Number       int       `json:"number"`
	SellPrice    int       `json:"sell_price"`
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
		sell_price INTEGER NOT NULL CHECK("sell_price" >= 0),
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

func NewItemsForSell() error {
	return nil
}

// SellItems()
