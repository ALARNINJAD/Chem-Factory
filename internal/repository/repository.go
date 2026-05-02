package repository

import (
	mat "chem-factory/internal/dto/material"
	"database/sql"
	"fmt"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type RepositoryManager interface {
	// repository
	ExecQuery(query string) error
	Transaction() (*sql.Tx, error)
	// User
	FindUserByUsername(username string) (*user, error)
	FindUserByID(id int) (*user, error)
	FindPasswordByID(id int) (string, error)
	FindPasswordByUsername(username string) (string, error)
	FindIDbyUsername(username string) (int, error)
	FindUsernameByID(id int) (string, error)
	SaveUser(tx *sql.Tx, username, password string) error
	IncreaseBalance(tx *sql.Tx, username string, amount int) error
	ReduceBalance(tx *sql.Tx, username string, amount int) error
	// Material
	EmptyMaterialStruct() *material
	EmptyMaterialSlice() []material
	FindBaseMaterialByID(id int) (*baseMaterial, error)
	GetMaterialsByUsername(username string) ([]material, error)
	SaveMaterial(m material) error
	SaveBaseMaterial(tx *sql.Tx, m mat.BaseMaterial) error
	FindIDbyMaterialName(name string) (int, error)
	FindMaterialByID(id int) (*material, error)
	FindMaterialByIngrID(firstID int, secondID int) (*material, error)
	FindMaterialIDByIngrID(firstID int, secondID int) (int, error)
	FindMatMixTimeByID(id int) (int, error)
	// shop
	ExportShop() ([]shop, error)
	EmptyShopStruct() *shop
	AddToShop(tx *sql.Tx, s shop) error
	FindShopIDByInfo(userID, materialID, price int) (int, error)
	FindShopByInfo(userID, materialID, price int) (*shop, error)
	ReduceShopNumberByID(tx *sql.Tx, id, number int) error
	IncreaseShopNumberByID(tx *sql.Tx, id, number int) error
	DeleteFromShopByID(tx *sql.Tx, id int) error
	// inventory
	EmptyInventoryStruct() *inventory
	EmptyInventorySlice() []inventory
	FindUserInvenByID(id int) (*inventory, error)
	FindInvenIDByUserIDmatID(userID, materialID int) (int, error)
	IncreaseInventoryByID(tx *sql.Tx, id, number int) error
	ReduceInventoryByID(tx *sql.Tx, id, number int) error
	AddToInventory(tx *sql.Tx, i inventory) error
	DeleteInventoryByID(tx *sql.Tx, id int) error
	GetInventoryItemsByUsername(username string) ([]inventory, error)
	// mixer
	EmptyMixerStruct() *mixer
	FindMixIDByUserIDIngrID(userID, firstID, secID int) (int, error)
	AddToMixer(tx *sql.Tx, m mixer) error
	FindMixRowByID(id int) (*mixer, error)
	DeleteMixByID(tx *sql.Tx, id int) error
}

type repositoryManager struct {
	db *sql.DB
}

func Init() *repositoryManager {
	var r repositoryManager
	var err error
	r.db, err = sql.Open("sqlite3", filepath.Join(".", "database", "database.db"))
	if err != nil {
		panic(fmt.Errorf("Repository, init open db: %w ", err))
	}
	return &r
}

func (r *repositoryManager) ExecQuery(query string) error {
	if _, err := r.db.Exec(query); err != nil {
		return fmt.Errorf("Repository, exec query: %w", err)
	}
	return nil
}

func (r *repositoryManager) Transaction() (*sql.Tx, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("Repository, transaction: %w", err)
	}
	return tx, nil
}
