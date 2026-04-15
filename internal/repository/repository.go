package repository

import (
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
	// Material
	SaveMaterial(m material) error
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

	createUserTable(r)
	createMaterialTable(r)
	createShopTable(r)
	createShopTable(r)
}
