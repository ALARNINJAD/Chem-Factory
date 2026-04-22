package repository

import (
	"database/sql"
	"fmt"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type RepositoryManager interface {
	// User
	FindUserByUsername(username string) (*user, error)
	FindUserByID(id int) (*user, error)
	FindPasswordByID(id int) (string, error)
	FindPasswordByUsername(username string) (string, error)
	FindIDbyUsername(username string) (int, error)
	FindUsernameByID(id int) (string, error)
	SaveUser(username string, password string) error
	IncreaseBalance(username string, amount int) error
	ReduceBalance(username string, amount int) error
	// Material
	SaveMaterial(m material) error
	FindIDbyMaterialName(name string) (int, error)
	FindBaseMaterials() ([]material, error)
	FindMaterialByID(id int) (*material, error)
	FindMaterialByIngrID(firstID int, secondID int) (*material, error)
	FindMaterialIDByIngrID(firstID int, secondID int) (int, error)
	FindMatMixTimeByID(id int) (int, error)
	// shop
	ExportShop() ([]shop, error)
	EmptyShopStruct() *shop
	FindShopIDByInfo(userID, materialID, price int) (int, error)
	FindShopByInfo(userID, materialID, price int) (*shop, error)
	ReduceShopNumberByID(id, number int) error
	IncreaseShopNumberByID(id, number int) error
	DeleteFromShopByID(id int) error
	// inventory
	FindInvenIDByUserIDmatID(userID, materialID int) (int, error)
	IncreaseInventoryByID(id, number int) error
	ReduceInventoryByID(id, number int) error
	AddToInventory(i inventory) error
	DeleteInventoryByID(id int) error
	GetInventoryByUsername(username string) ([]inventory, error)
	// mixer
	FindMixIDByUserIDIngrID(userID, firstID, secID int) (int, error)
	AddToMixer(m mixer) error
	FindMixRowByID(id int) (*mixer, error)
}

type repositoryManager struct {
	db *sql.DB
}

func Init() *repositoryManager {

	var r repositoryManager
	var err error

	r.db, err = sql.Open("sqlite3", filepath.Join(".", "database", "database.db"))
	if err != nil {
		panic(fmt.Errorf("Repository, init open db: %w", err))
	}

	createUserTable(&r)
	createMaterialTable(&r)
	createShopTable(&r)
	createInventoryTable(&r)
	createMixerTable(&r)

	return &r
}
