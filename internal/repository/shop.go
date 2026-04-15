package repository

import "time"

type Shop struct {
	ID           int
	UserID       int
	MaterialID   int
	Username     string
	MaterialName string
	Number       int
	SellPrice    int
	DateTime     time.Time
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
