package main

import (
	"chem-factory/internal/database/sqlite"
	"chem-factory/internal/domain"
	marketsqlite "chem-factory/internal/modules/market/adapter/sqlite"
	materialsqlite "chem-factory/internal/modules/material/adapter/sqlite"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {

	db := sqlite.New()
	ctx := context.Background()

	createTables(ctx, db)
	createMaterials(ctx, db)
	createShopMaterials(ctx, db)
}

func createTables(ctx context.Context, db *sqlite.Database) {

	if _, err := db.Extract(ctx).ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS users (
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
	CREATE TABLE IF NOT EXISTS materials (
		id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		first_ingredient_id INTEGER CHECK("first_ingredient_id" <= "second_ingredient_id"),
		second_ingredient_id INTEGER CHECK("first_ingredient_id" <= "second_ingredient_id"),
		name TEXT NOT NULL UNIQUE,
		price INTEGER NOT NULL CHECK("price" >= 0),
		mix_time INTEGER CHECK("mix_time" >= 0),
		UNIQUE("first_ingredient_id","second_ingredient_id"),
		FOREIGN KEY("user_id") REFERENCES "users"("id")
		FOREIGN KEY("first_ingredient_id") REFERENCES "materials"("id"),
		FOREIGN KEY("second_ingredient_id") REFERENCES "materials"("id")
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
		FOREIGN KEY("material_id") REFERENCES "materials"("id"),
		FOREIGN KEY("user_id") REFERENCES "users"("id")
	)`); err != nil {
		panic(fmt.Errorf("Migration, create inventory table: %w", err))
	}

	if _, err := db.Extract(ctx).ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS market (
		id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
		user_id	INTEGER,
		material_id	INTEGER NOT NULL,
		amount INTEGER NOT NULL CHECK("amount" >= 1),
		date_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UNIQUE("user_id","material_id"),
		FOREIGN KEY("material_id") REFERENCES "materials"("id"),
		FOREIGN KEY("user_id") REFERENCES "users"("id")
	)`); err != nil {
		panic(fmt.Errorf("Migration, create market table: %w", err))
	}

	if _, err := db.Extract(ctx).ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS mixes (
		id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		first_ingredient_id INTEGER CHECK("first_ingredient_id" <= "second_ingredient_id"),
		second_ingredient_id INTEGER CHECK("first_ingredient_id" <= "second_ingredient_id"),
		amount INTEGER NOT NULL CHECK(10 >= "amount" AND "amount" >= 1),
		date_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UNIQUE("user_id","first_ingredient_id","second_ingredient_id"),
		FOREIGN KEY("user_id") REFERENCES "users"("id"),
		FOREIGN KEY("first_ingredient_id") REFERENCES "materials"("id"),
		FOREIGN KEY("second_ingredient_id") REFERENCES "materials"("id")
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
		Name                 string `json:"name"`
		Price                int    `json:"price"`
		MixTime              int    `json:"mix_time"`
		FirstIngredientName  string `json:"first_ingredient_name"`
		SecondIngredientName string `json:"second_ingredient_name"`
	}
	if err := json.Unmarshal(data, &materials); err != nil {
		panic(fmt.Errorf("Migration, create admin materials, parse file: %w ", err))
	}

	repo := materialsqlite.NewMaterialRepo(db)

	for _, m := range materials {
		if _, err := repo.FindIDByName(ctx, m.Name); err == nil {
			continue
		} else if !errors.Is(err, sql.ErrNoRows) {
			panic(fmt.Errorf("Migration, check material existence: %w", err))
		}

		var (
			firstID  uint
			secondID uint
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
			firstID = fID
			secondID = sID
		}

		if err := repo.Add(ctx, domain.Material{
			Name:               m.Name,
			Price:              m.Price,
			MixTime:            m.MixTime,
			FirstIngredientID:  firstID,
			SecondIngredientID: secondID,
		}); err != nil {
			log.Println("Could not add material:", m.Name, "error:", err)
		}
	}
}

func createShopMaterials(ctx context.Context, db *sqlite.Database) {
	rows, err := db.Extract(ctx).QueryContext(ctx, `
		SELECT id, first_ingredient_id, second_ingredient_id, name
		FROM materials
	`)
	if err != nil {
		panic(fmt.Errorf("Migration, read materials for shop seed: %w", err))
	}
	defer rows.Close()

	var rawMaterials []struct {
		ID   uint
		Name string
	}

	for rows.Next() {
		var (
			id       uint
			firstID  sql.NullInt64
			secondID sql.NullInt64
			name     string
		)
		if err := rows.Scan(&id, &firstID, &secondID, &name); err != nil {
			panic(fmt.Errorf("Migration, scan material for shop seed: %w", err))
		}
		if firstID.Valid || secondID.Valid {
			continue
		}
		rawMaterials = append(rawMaterials, struct {
			ID   uint
			Name string
		}{ID: id, Name: name})
	}

	if err := rows.Err(); err != nil {
		panic(fmt.Errorf("Migration, read materials rows for shop seed: %w", err))
	}

	marketRepo := marketsqlite.NewMarketRepo(db)
	for _, mat := range rawMaterials {
		var count int
		if err := db.Extract(ctx).QueryRowContext(ctx,
			`SELECT COUNT(1) FROM market WHERE user_id IS NULL AND material_id = ?`,
			mat.ID,
		).Scan(&count); err != nil {
			panic(fmt.Errorf("Migration, check shop material existence: %w", err))
		}
		if count > 0 {
			continue
		}
		if err := marketRepo.Add(ctx, domain.Market{
			UserID:     0,
			MaterialID: mat.ID,
			Amount:     10,
			DateTime:   time.Now(),
		}); err != nil {
			log.Println("Could not add shop material:", mat.Name, "error:", err)
		}
	}
}
