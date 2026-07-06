package main

import (
	"chem-factory/internal/database/sqlite"
	"context"
	"fmt"
)

func main() {

	db := sqlite.New()
	ctx := context.Background()
	createTables(ctx, db)
}

func createTables(ctx context.Context, db *sqlite.Database) {

	if _, err := db.Extract(ctx).ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS user (
		id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		balance INTEGER DEFAULT 0 CHECK("balance" >= 0),
		xp INTEGER DEFAULT 0 CHECK("xp" >= 0),
		level INTEGER DEFAULT 0 CHECK("level" >= 0)
	)`); err != nil {
		panic(fmt.Errorf("Migration, create user table: %w", err))
	}

	if _, err := db.Extract(ctx).ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS material (
		id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		first_ingredient_id INTEGER CHECK("first_ingredient_id" <= "second_ingredient_id"),
		second_ingredient_id INTEGER CHECK("first_ingredient_id" <= "second_ingredient_id"),
		name TEXT NOT NULL UNIQUE,
		price INTEGER NOT NULL CHECK("price" >= 0),
		mix_time INTEGER CHECK("mix_time" >= 0),
		UNIQUE("first_ingredient_id","second_ingredient_id"),
		FOREIGN KEY("user_id") REFERENCES "user"("id")
		FOREIGN KEY("first_ingredient_id") REFERENCES "material"("id"),
		FOREIGN KEY("second_ingredient_id") REFERENCES "material"("id")
	)`); err != nil {
		panic(fmt.Errorf("Migration, create material table: %w", err))
	}

	if _, err := db.Extract(ctx).ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS inventory (
		id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
		user_id	INTEGER NOT NULL,
		material_id	INTEGER NOT NULL,
		amount INTEGER NOT NULL CHECK("amount" >= 1),
		date_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UNIQUE("user_id","material_id"),
		FOREIGN KEY("material_id") REFERENCES "material"("id"),
		FOREIGN KEY("user_id") REFERENCES "user"("id")
	)`); err != nil {
		panic(fmt.Errorf("Migration, create inventory table: %w", err))
	}

	if _, err := db.Extract(ctx).ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS market (
		id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
		user_id	INTEGER NOT NULL,
		material_id	INTEGER NOT NULL,
		amount INTEGER NOT NULL CHECK("amount" >= 1),
		price INTEGER NOT NULL CHECK("price" >= 0),
		date_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UNIQUE("user_id","material_id","price"),
		FOREIGN KEY("material_id") REFERENCES "material"("id"),
		FOREIGN KEY("user_id") REFERENCES "user"("id")
	)`); err != nil {
		panic(fmt.Errorf("Migration, create market table: %w", err))
	}

	if _, err := db.Extract(ctx).ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS mixer (
		id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		first_ingredient_id INTEGER CHECK("first_ingredient_id" <= "second_ingredient_id"),
		second_ingredient_id INTEGER CHECK("first_ingredient_id" <= "second_ingredient_id"),
		amount INTEGER NOT NULL CHECK(10 >= "amount" AND "amount" >= 1),
		date_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UNIQUE("user_id","first_ingredient_id","second_ingredient_id"),
		FOREIGN KEY("user_id") REFERENCES "user"("id"),
		FOREIGN KEY("first_ingredient_id") REFERENCES "material"("id"),
		FOREIGN KEY("second_ingredient_id") REFERENCES "material"("id")
	)`); err != nil {
		panic(fmt.Errorf("Migration, create mixer table: %w", err))
	}
}
