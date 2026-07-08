package main

import (
	"chem-factory/internal/database/sqlite"
	"chem-factory/internal/domain"
	materialsqlite "chem-factory/internal/modules/material/adapter/sqlite"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func main() {

	db := sqlite.New()
	ctx := context.Background()

	createTables(ctx, db)
	createMaterials(ctx, db)
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

func createMaterials(ctx context.Context, db *sqlite.Database) {

	data, err := os.ReadFile(filepath.Join(".", "pkg", "material", "default_materials.json"))
	if err != nil {
		panic(fmt.Errorf("Migration, create admin materials, read file: %w ", err))
	}

	var materials []struct {
		Name                  string `json:"name"`
		Price                 int    `json:"price"`
		MixTime               int    `json:"mix_time"`
		FirstIngredientName   string `json:"first_ingredient_name"`
		SecondIngredientName  string `json:"second_ingredient_name"`
	}
	if err := json.Unmarshal(data, &materials); err != nil {
		panic(fmt.Errorf("Migration, create admin materials, parse file: %w ", err))
	}

	repo := materialsqlite.NewMaterialRepo(db)

	for _, m := range materials {
		var (
			firstID  *int
			secondID *int
		)

		if m.FirstIngredientName != "" && m.SecondIngredientName != "" {
			fID, err := repo.FindIDByName(ctx, m.FirstIngredientName)
			if err != nil {
				panic(fmt.Errorf("Migration, create admin materials, find first ingredient: %w", err))
			}
			sID, err := repo.FindIDByName(ctx, m.SecondIngredientName)
			if err != nil {
				panic(fmt.Errorf("Migration, create admin materials, find second ingredient: %w", err))
			}

			if fID > sID {
				fID, sID = sID, fID
			}
			firstID = &fID
			secondID = &sID
		}

		if err := repo.Add(ctx, domain.Material{
			Name:              m.Name,
			Price:             m.Price,
			MixTime:           m.MixTime,
			FirstIngredientID: firstID,
			SecondIngredientID: secondID,
		}); err != nil {
			panic(fmt.Errorf("Migration, create admin materials, add material: %w", err))
		}
	}
}
