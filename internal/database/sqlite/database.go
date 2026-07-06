package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
)

type Database struct { database *sql.DB }

func New() *Database {
	database, err := sql.Open("sqlite3", filepath.Join(".", "internal", "database", "sqlite", "database.db"))
	if err != nil {
		panic(fmt.Errorf("could not open sqlite database: %w", err))
	}
	return &Database{database: database}
}

type txKey struct{}

func (db *Database) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := db.database.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin sqlite transaction: %w", err)
	}
	if err := fn(context.WithValue(ctx, txKey{}, tx)); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (db *Database) Extract(ctx context.Context) interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
} {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}
	return db.database
}
