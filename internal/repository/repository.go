package repository

import (
	"chem-factory/internal/model"
	"database/sql"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type RepositoryManager interface {
	// User
	SaveNewUser(user model.User) error
	ExportUserByUsername(user *model.User) error
	ExportUserByID(user *model.User) error
	ExportPasswordByID(id int) (string, error)
	ExportPasswordByUsername(username string) (string, error)
	ExportIDbyUsername(username string) (int, error)
	ExportUsernameByID(id int) (string, error)
	// Material
	SaveMaterial(material *model.Material) error
	// shop
	// inventory
}

type repositoryManager struct {
	db *sql.DB
}

func Init() *repositoryManager {
	var r repositoryManager
	var err error
	r.db, err = sql.Open("sqlite3", filepath.Join(".", "database", "database.db"))
	if err != nil {
		panic("Could not open database.")
	}
	r.creatTables()
	r.createMaterials()
	return &r
}

func (r *repositoryManager) creatTables() {

	userQuery := `
	CREATE TABLE IF NOT EXISTS user (
		id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		balance INTEGER DEFAULT 100 CHECK("balance" >= 0),
		xp INTEGER DEFAULT 0 CHECK("xp" >= 0),
		level INTEGER DEFAULT 0 CHECK("level" >= 0)
	)`
	if _, err := r.db.Exec(userQuery); err != nil {
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
	if _, err := r.db.Exec(materialQuery); err != nil {
		panic("Could not create material table.")
	}

	inventoryQuery := `
	CREATE TABLE IF NOT EXISTS inventory (
		id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
		user_id	INTEGER NOT NULL,
		material_id	INTEGER NOT NULL,
		number INTEGER NOT NULL DEFAULT 0 CHECK("number" >= 1),
		UNIQUE("user_id","material_id"),
		FOREIGN KEY("material_id") REFERENCES "material"("id"),
		FOREIGN KEY("user_id") REFERENCES "user"("id")
	)`
	if _, err := r.db.Exec(inventoryQuery); err != nil {
		panic("Could not create inventory table.")
	}

	shopQuery := `
	CREATE TABLE IF NOT EXISTS shop (
		id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
		user_id	INTEGER NOT NULL,
		material_id	INTEGER NOT NULL,
		number INTEGER NOT NULL DEFAULT 0 CHECK(10 >= "number" >= 1),
		sell_price INTEGER NOT NULL CHECK("sell_price" >= 0),
		date_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UNIQUE("user_id","material_id"),
		FOREIGN KEY("material_id") REFERENCES "material"("id"),
		FOREIGN KEY("user_id") REFERENCES "user"("id")
	)`
	if _, err := r.db.Exec(shopQuery); err != nil {
		panic("Could not create shop table.")
	}
}
