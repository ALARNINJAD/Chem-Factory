package sqlite

import (
	"chem-factory/internal/core/domain"
	"context"
	"database/sql"
	"fmt"
)

type UserRepo struct{ db *DB }

func NewUserRepo(db *DB) *UserRepo { return &UserRepo{db: db} }

func (r *UserRepo) FindByUsername(ctx context.Context, username string) (domain.User, error) {
	var user domain.User

	err := r.db.extract(ctx).QueryRowContext(ctx,
		`SELECT id, username, password, balance, xp, level FROM user WHERE username = ?`,
		username,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Balance,
		&user.XP,
		&user.Level,
	)

	if err != nil {
		return domain.User{}, fmt.Errorf("find user by username: %w", err)
	}

	return user, nil
}

func (r *UserRepo) FindByID(ctx context.Context, id int) (domain.User, error) {
	var user domain.User

	err := r.db.extract(ctx).QueryRowContext(ctx,
		`SELECT id, username, password, balance, xp, level FROM user WHERE id = ?`,
		id,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Balance,
		&user.XP,
		&user.Level,
	)

	if err != nil {
		return domain.User{}, fmt.Errorf("find user by id: %w", err)
	}

	return user, nil
}

func (r *UserRepo) FindPasswordByID(ctx context.Context, id int) (string, error) {
	var password string

	err := r.db.extract(ctx).QueryRowContext(ctx,
		`SELECT password FROM user WHERE id = ?`,
		id,
	).Scan(&password)

	if err != nil {
		return "", fmt.Errorf("find password by id: %w", err)
	}

	return password, nil
}

func (r *UserRepo) FindPasswordByUsername(ctx context.Context, username string) (string, error) {
	var password string

	err := r.db.extract(ctx).QueryRowContext(ctx,
		`SELECT password FROM user WHERE username = ?`,
		username,
	).Scan(&password)

	if err != nil {
		return "", fmt.Errorf("find password by username: %w", err)
	}

	return password, nil
}

func (r *UserRepo) FindIDByUsername(ctx context.Context, username string) (int, error) {
	var id int

	err := r.db.extract(ctx).QueryRowContext(ctx,
		`SELECT id FROM user WHERE username = ?`,
		username,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("find id by username: %w", err)
	}

	return id, nil
}

func (r *UserRepo) FindUsernameByID(ctx context.Context, id int) (string, error) {
	var username string

	err := r.db.extract(ctx).QueryRowContext(ctx,
		`SELECT username FROM user WHERE id = ?`,
		id,
	).Scan(&username)

	if err != nil {
		return "", fmt.Errorf("find username by id: %w", err)
	}

	return username, nil
}

func (r *UserRepo) Add(ctx context.Context, tx *sql.Tx, username, password string) error {
	_, err := tx.ExecContext(ctx,
		`INSERT INTO user(username, password) VALUES (?, ?)`,
		username,
		password,
	)

	if err != nil {
		return fmt.Errorf("add user: %w", err)
	}

	return nil
}

func (r *UserRepo) IncreaseBalance(ctx context.Context, tx *sql.Tx, username string, amount int) error {
	_, err := tx.ExecContext(ctx,
		`UPDATE user SET balance = balance + ? WHERE username = ?`,
		amount,
		username,
	)

	if err != nil {
		return fmt.Errorf("increase balance: %w", err)
	}

	return nil
}

func (r *UserRepo) ReduceBalance(ctx context.Context, tx *sql.Tx, username string, amount int) error {
	_, err := tx.ExecContext(ctx,
		`UPDATE user SET balance = balance - ? WHERE username = ?`,
		amount,
		username,
	)

	if err != nil {
		return fmt.Errorf("reduce balance: %w", err)
	}

	return nil
}