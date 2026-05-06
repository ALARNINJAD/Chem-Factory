package main

import (
	"chem-factory/internal/auth"
	"chem-factory/internal/dto/material"
	"chem-factory/internal/repository"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := repository.New()
	a := auth.New()
	createTables(r)
	createAdminUser(r, a)
	createAdminMaterials(r)
}

func createTables(r *repository.Manager) {

	if err := r.ExecQuery(`
	CREATE TABLE IF NOT EXISTS user (
		id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		balance INTEGER DEFAULT 200 CHECK("balance" >= 0),
		xp INTEGER DEFAULT 0 CHECK("xp" >= 0),
		level INTEGER DEFAULT 0 CHECK("level" >= 0)
	)`); err != nil {
		panic(fmt.Errorf("Migration, create user table: %w", err))
	}

	if err := r.ExecQuery(`
	CREATE TABLE IF NOT EXISTS material (
		id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		first_ingredient_id INTEGER,
		second_ingredient_id INTEGER,
		username TEXT NOT NULL,
		name TEXT NOT NULL UNIQUE,
		first_ingredient_name TEXT,
		second_ingredient_name TEXT,
		sell_price INTEGER NOT NULL CHECK("sell_price" >= 0),
		buy_price INTEGER NOT NULL CHECK("buy_price" >= 0),
		mix_time INTEGER CHECK("mix_time" >= 0),
		UNIQUE("first_ingredient_id","second_ingredient_id"),
		UNIQUE("first_ingredient_name","second_ingredient_name"),
		FOREIGN KEY("user_id") REFERENCES "user"("id")
		FOREIGN KEY("first_ingredient_id") REFERENCES "material"("id"),
		FOREIGN KEY("second_ingredient_id") REFERENCES "material"("id")
	)`); err != nil {
		panic(fmt.Errorf("Migration, create material table: %w", err))
	}

	if err := r.ExecQuery(`
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
	)`); err != nil {
		panic(fmt.Errorf("Migration, create inventory table: %w", err))
	}

	if err := r.ExecQuery(`
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
	)`); err != nil {
		panic(fmt.Errorf("Migration, create shop table: %w", err))
	}

	if err := r.ExecQuery(`
	CREATE TABLE IF NOT EXISTS mixer (
		id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		first_ingredient_id INTEGER CHECK("first_ingredient_id" <= "second_ingredient_id"),
		second_ingredient_id INTEGER CHECK("first_ingredient_id" <= "second_ingredient_id"),
		username TEXT NOT NULL,
		first_ingredient_name TEXT NOT NULL,
		second_ingredient_name TEXT NOT NULL,
		number INTEGER NOT NULL CHECK("number" >= 1),
		date_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UNIQUE("user_id","first_ingredient_id","second_ingredient_id"),
		UNIQUE("username","first_ingredient_name","second_ingredient_name"),
		FOREIGN KEY("user_id") REFERENCES "user"("id"),
		FOREIGN KEY("first_ingredient_id") REFERENCES "material"("id"),
		FOREIGN KEY("second_ingredient_id") REFERENCES "material"("id")
	)`); err != nil {
		panic(fmt.Errorf("Migration, create mixer table: %w", err))
	}
}

func createAdminUser(r *repository.Manager, a *auth.Manager) {

	_, err := r.User.FindIDbyUsername(os.Getenv("ADMIN_NAME"))
	if err == nil {
		return
	}

	hashedPassword, err := a.Hash.HashPassword(os.Getenv("ADMIN_PASSWORD"))
	if err != nil {
		panic(fmt.Errorf("Migration, create admin: %w ", err))
	}

	tx, err := r.Transaction()
	if err != nil {
		panic(fmt.Errorf("Migration, create admin, transaction: %w ", err))
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err := r.User.Add(tx, os.Getenv("ADMIN_NAME"), hashedPassword); err != nil {
		panic(fmt.Errorf("Migration, create admin: %w ", err))
	}

	if err = tx.Commit(); err != nil {
		panic(fmt.Errorf("Migration, create admin, commit: %w ", err))
	}
}

func createAdminMaterials(r *repository.Manager) {

	adminID, err := r.User.FindIDbyUsername(os.Getenv("ADMIN_NAME"))
	if err != nil {
		panic(fmt.Errorf("Migration, create admin materials: %w ", err))
	}

	data, err := os.ReadFile(filepath.Join(".", "cmd", "migration", "materials.json"))
	if err != nil {
		panic(fmt.Errorf("Migration, create admin materials, read file: %w ", err))
	}

	materials := r.Material.EmptySlice()

	if err = json.Unmarshal(data, &materials); err != nil {
		panic(fmt.Errorf("Migration, create admin materials, unmarshal: %w ", err))
	}

	for _, m := range materials {

		m.Username = os.Getenv("ADMIN_NAME")
		m.UserID = adminID

		tx, err := r.Transaction()
		if err != nil {
			panic(fmt.Errorf("Migration, create admin, transaction: %w ", err))
		}
		defer func() {
			if err != nil {
				tx.Rollback()
			}
		}()

		if m.Name == m.FirstIngredientName && m.Name == m.SecondIngredientName {
			err = r.Material.AddBase(tx, material.Base{
				Name:      m.Name,
				UserID:    m.UserID,
				Username:  m.Username,
				SellPrice: m.SellPrice,
				BuyPrice:  m.BuyPrice,
			})
			if err != nil {
				panic(fmt.Errorf("Migration, create admin materials, %v: %w ", m, err))
			}
			if err = tx.Commit(); err != nil {
				panic(fmt.Errorf("Migration, create admin materials, commit: %w ", err))
			}
			continue
		}

		m.FirstIngredientID, err = r.Material.FindIDByName(m.FirstIngredientName)
		if err != nil {
			panic(fmt.Errorf("Migration, create admin materials, %v: %w ", m, err))
		}

		m.SecondIngredientID, err = r.Material.FindIDByName(m.SecondIngredientName)
		if err != nil {
			panic(fmt.Errorf("Migration, create admin materials, %v: %w ", m, err))
		}

		if m.FirstIngredientID > m.SecondIngredientID {
			m.FirstIngredientID, m.SecondIngredientID = m.SecondIngredientID, m.FirstIngredientID
			m.FirstIngredientName, m.SecondIngredientName = m.SecondIngredientName, m.FirstIngredientName
		}

		if err = r.Material.Add(m); err != nil {
			panic(fmt.Errorf("Migration, create admin materials, %v: %w ", m, err))
		}

		if err = tx.Commit(); err != nil {
			panic(fmt.Errorf("Migration, create admin materials, commit: %w ", err))
		}
	}
}
