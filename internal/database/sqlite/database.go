package sqlite

import (
	"context"
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	db *gorm.DB
}

func New() *Database {
	log.Println("Initializing sqlite with GORM")

	db, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("could not open sqlite database with gorm: %w", err))
	}
	return &Database{db: db}
}

func (database *Database) DB() *gorm.DB {
	return database.db
}

type txKey struct{}

func (db *Database) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return db.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, txKey{}, tx))
	})
}

func (db *Database) Extract(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return db.db.WithContext(ctx)
}

func IsTx(ctx context.Context) bool {
	_, ok := ctx.Value(txKey{}).(*gorm.DB)
	return ok
}
