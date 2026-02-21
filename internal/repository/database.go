package repository

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init() {
	var err error
	DB, err = sql.Open("sqlite3", "../../database/database.db")
	if err != nil {
		panic("Could not open to database.")
	}
	creatTables()
}

func creatTables() {

	usersQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		balance INTEGER DEFAULT 0 CHECK("balance" >= 0),
		xp INTEGER DEFAULT 0 CHECK("xp" >= 0),
		level INTEGER DEFAULT 0 CHECK("level" >= 0)
	)`
	if _, err := DB.Exec(usersQuery); err != nil {
		panic("Could not create users table.")
	}

	materialsQuery := `
	CREATE TABLE IF NOT EXISTS materials (
		id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		sell_price INTEGER NOT NULL CHECK("sell_price" >= 0),
		buy_price INTEGER NOT NULL CHECK("buy_price" >= 0),
		first_ingredient_id INTEGER,
		second_ingredient_id INTEGER,
		mix_time INTEGER NOT NULL CHECK("mix_time" >= 0),
		UNIQUE("first_ingredient_id","second_ingredient_id"),
		FOREIGN KEY("first_ingredient_id") REFERENCES "materials"("id"),
		FOREIGN KEY("second_ingredient_id") REFERENCES "materials"("id")
	)`
	if _, err := DB.Exec(materialsQuery); err != nil {
		panic("Could not create materials table.")
	}

	inventoryQuery := `
	CREATE TABLE IF NOT EXISTS inventory (
		id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
		user_id	INTEGER NOT NULL ,
		material_id	INTEGER NOT NULL,
		number INTEGER NOT NULL DEFAULT 0 CHECK("number" >= 0),
		UNIQUE("user_id","material_id"),
		FOREIGN KEY("material_id") REFERENCES "materials"("id"),
		FOREIGN KEY("user_id") REFERENCES "users"("id")
	)`
	if _, err := DB.Exec(inventoryQuery); err != nil {
		panic("Could not create inventory table.")
	}
}
