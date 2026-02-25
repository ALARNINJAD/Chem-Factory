package repository

import (
	"chem-factory/internal/model"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func Init() {
	rootPath, err := os.Getwd()
	if err != nil {
		panic("Could not find database.")
	}
	dbPath := filepath.Join(rootPath, "database", "database.db")
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		panic("Could not open database.")
	}
	creatTables()
	createMaterials()
}

func creatTables() {

	userQuery := `
	CREATE TABLE IF NOT EXISTS user (
		id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		balance INTEGER DEFAULT 100 CHECK("balance" >= 0),
		xp INTEGER DEFAULT 0 CHECK("xp" >= 0),
		level INTEGER DEFAULT 0 CHECK("level" >= 0)
	)`
	if _, err := db.Exec(userQuery); err != nil {
		panic("Could not create users table.")
	}

	materialQuery := `
	CREATE TABLE IF NOT EXISTS material (
		id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		sell_price INTEGER NOT NULL CHECK("sell_price" >= 0),
		buy_price INTEGER NOT NULL CHECK("buy_price" >= 0),
		first_ingredient_id INTEGER,
		second_ingredient_id INTEGER,
		mix_time INTEGER NOT NULL CHECK("mix_time" >= 0),
		UNIQUE("first_ingredient_id","second_ingredient_id"),
		FOREIGN KEY("first_ingredient_id") REFERENCES "material"("id"),
		FOREIGN KEY("second_ingredient_id") REFERENCES "material"("id")
	)`
	if _, err := db.Exec(materialQuery); err != nil {
		panic("Could not create material table.")
	}

	inventoryQuery := `
	CREATE TABLE IF NOT EXISTS inventory (
		id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
		user_id	INTEGER NOT NULL ,
		material_id	INTEGER NOT NULL,
		number INTEGER NOT NULL DEFAULT 0 CHECK("number" >= 0),
		UNIQUE("user_id","material_id"),
		FOREIGN KEY("material_id") REFERENCES "material"("id"),
		FOREIGN KEY("user_id") REFERENCES "user"("id")
	)`
	if _, err := db.Exec(inventoryQuery); err != nil {
		panic("Could not create inventory table.")
	}
}

func createMaterials() {

	var materials []material

	rootPath, err := os.Getwd()
	if err != nil {
		panic("Could not find materials file.")
	}
	materialsPath := filepath.Join(rootPath, "configs", "materials.json")

	data, err := os.ReadFile(materialsPath)
	if err != nil {
		panic("Could not open materials file.")
	}

	json.Unmarshal(data, &materials)

	for _, m := range materials {
		if err = m.new(); err != nil {
			panic(fmt.Sprintf("Could not add the %s material to database.", m.Name))
		}
	}
}

func ExportUserPassword(usr model.User) (string, error) {
	u := user{User: usr}
	return u.password()
}

func SaveUser(usr model.User) error {
	u := user{User: usr}
	return u.save()
}
