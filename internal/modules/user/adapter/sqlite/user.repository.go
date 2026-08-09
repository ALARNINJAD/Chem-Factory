package sqlite

import (
	database "chem-factory/internal/database/sqlite"
	"chem-factory/internal/domain"
	"chem-factory/pkg/reedam"
	"context"
)

type UserRepo struct{ db *database.Database }

func NewUserRepo(db *database.Database) *UserRepo { return &UserRepo{db: db} }

func (r *UserRepo) FindByUsername(ctx context.Context, username string) (domain.User, error) {
	var user domain.User

	err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT id, username, password, balance, xp, level FROM users WHERE username = ?`,
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
		return domain.User{}, reedam.InternalError(err)
	}

	return user, nil
}

func (r *UserRepo) FindByID(ctx context.Context, id uint) (domain.User, error) {
	var user domain.User

	err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT id, username, password, balance, xp, level FROM users WHERE id = ?`,
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
		return domain.User{}, reedam.InternalError(err)
	}

	return user, nil
}

func (r *UserRepo) FindPasswordByID(ctx context.Context, id uint) (string, error) {
	var password string

	err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT password FROM users WHERE id = ?`,
		id,
	).Scan(&password)

	if err != nil {
		return "", reedam.InternalError(err)
	}

	return password, nil
}

func (r *UserRepo) FindPasswordByUsername(ctx context.Context, username string) (string, error) {
	var password string

	err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT password FROM users WHERE username = ?`,
		username,
	).Scan(&password)

	if err != nil {
		return "", reedam.InternalError(err)
	}

	return password, nil
}

func (r *UserRepo) FindIDByUsername(ctx context.Context, username string) (uint, error) {
	var id uint

	err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT id FROM users WHERE username = ?`,
		username,
	).Scan(&id)
	if err != nil {
		return 0, reedam.InternalError(err)
	}

	return id, nil
}

func (r *UserRepo) FindUsernameByID(ctx context.Context, id uint) (string, error) {
	var username string

	err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT username FROM users WHERE id = ?`,
		id,
	).Scan(&username)

	if err != nil {
		return "", reedam.InternalError(err)
	}

	return username, nil
}

func (r *UserRepo) Add(ctx context.Context, user domain.User) error {
	_, err := r.db.Extract(ctx).ExecContext(ctx,
		`INSERT INTO users(username, password, balance) VALUES (?, ?, ?)`,
		user.Username,
		user.Password,
		user.Balance,
	)

	if err != nil {
		return reedam.InternalError(err)
	}

	return nil
}

func (r *UserRepo) IncreaseBalanceByID(ctx context.Context, id uint, amount int) error {
	_, err := r.db.Extract(ctx).ExecContext(ctx,
		`UPDATE users SET balance = balance + ? WHERE id = ?`,
		amount,
		id,
	)
	if err != nil {
		return reedam.InternalError(err)
	}
	return nil
}

func (r *UserRepo) ReduceBalanceByID(ctx context.Context, id uint, amount int) error {
	_, err := r.db.Extract(ctx).ExecContext(ctx,
		`UPDATE users SET balance = balance - ? WHERE id = ?`,
		amount,
		id,
	)
	if err != nil {
		return reedam.InternalError(err)
	}
	return nil
}

func (r *UserRepo) FindBalanceByID(ctx context.Context, id uint) (int, error) {
	var balance int
	err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT balance FROM users WHERE id = ?`,
		id,
	).Scan(&balance)
	if err != nil {
		return 0, reedam.InternalError(err)
	}
	return balance, nil
}

func (r *UserRepo) UpdateBalanceByID(ctx context.Context, id uint, balance int) error {
	_, err := r.db.Extract(ctx).ExecContext(ctx,
		`UPDATE users SET balance = ? WHERE id = ?`,
		balance,
		id,
	)
	if err != nil {
		return reedam.InternalError(err)
	}
	return nil
}

func (r *UserRepo) UpdateLevelXPByID(ctx context.Context, user domain.User) error {
	_, err := r.db.Extract(ctx).ExecContext(ctx,
		`UPDATE users SET level = ?, xp = ? WHERE id = ?`,
		user.Level,
		user.XP,
		user.ID,
	)
	if err != nil {
		return reedam.InternalError(err)
	}
	return nil
}
