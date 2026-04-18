package repository

import (
	"chem-factory/internal/model"
	"database/sql"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type RepositoryManager interface {
	// User
	SaveNewUser(username string, password string) error
	ExportUserByUsername(username string) (*user, error)
	ExportUserByID(id int) (*user, error)
	ExportPasswordByID(id int) (string, error)
	ExportPasswordByUsername(username string) (string, error)
	ExportIDbyUsername(username string) (int, error)
	ExportUsernameByID(id int) (string, error)
	IncreaseBalance(username string, amount int) error
	ReduceBalance(username string, amount int) error
	// Material
	SaveMaterial(mtrl model.Material) error
	ExportIDbyMaterialName(name string) (int, error)
	ExportBaseMaterials() ([]baseMaterial, error)
	// shop
	ExportShop() ([]shop, error)
	AddToShop(s model.Shop) error
	RemoveFromShop(s model.Shop) error
	// inventory
	AddToInventory(i model.Inventory) error
	RemoveFromInventory(i model.Inventory) error
	ExportInventory(username string) ([]inventory, error)
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
	return &r
}

func (r *repositoryManager) creatTables() {

	createUserTable(r)
	createMaterialTable(r)
	createShopTable(r)
	createInventoryTable(r)
}
