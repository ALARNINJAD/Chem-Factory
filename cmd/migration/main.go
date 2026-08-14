package main

import (
	"chem-factory/internal/database/sqlite"
	inventorysqlite "chem-factory/internal/modules/inventory/adapter/sqlite"
	marketsqlite "chem-factory/internal/modules/market/adapter/sqlite"
	materialsqlite "chem-factory/internal/modules/material/adapter/sqlite"
	mixersqlite "chem-factory/internal/modules/mixer/adapter/sqlite"
	usersqlite "chem-factory/internal/modules/user/adapter/sqlite"
	"context"
	"fmt"
)

func main() {
	db := sqlite.New()
	ctx := context.Background()

	createTables(ctx, db)
}

func createTables(ctx context.Context, db *sqlite.Database) {
	gdb := db.Extract(ctx)
	if err := gdb.AutoMigrate(
		&usersqlite.User{},
		&materialsqlite.Material{},
		&inventorysqlite.Inventory{},
		&marketsqlite.Market{},
		&mixersqlite.Mix{},
	); err != nil {
		panic(fmt.Errorf("Migration, auto migrate: %w", err))
	}
}
