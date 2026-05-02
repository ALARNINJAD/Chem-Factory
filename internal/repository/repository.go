package repository

import (
	"database/sql"
	"fmt"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type Manager struct {
	User      *userManager
	Shop      *shopManager
	Mixer     *mixerManager
	Inventory *inventoryManager
	Material  *materialManager
	db        *sql.DB
}

func New() *Manager {
	db, err := sql.Open("sqlite3", filepath.Join(".", "database", "database.db"))
	if err != nil {
		panic(fmt.Errorf("Repository, new, open db: %w", err))
	}
	return &Manager{
		User:      NewUserManager(db),
		Shop:      NewShopManager(db),
		Mixer:     NewMixerManager(db),
		Inventory: NewInventoryManager(db),
		Material:  NewMaterialManager(db),
		db:        db,
	}
}

func (r *Manager) ExecQuery(query string) error {
	if _, err := r.db.Exec(query); err != nil {
		return fmt.Errorf("Repository, exec query: %w", err)
	}
	return nil
}

func (r *Manager) Transaction() (*sql.Tx, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("Repository, transaction: %w", err)
	}
	return tx, nil
}
