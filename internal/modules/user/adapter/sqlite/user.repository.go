package sqlite

import (
	"chem-factory/internal/domain"
	"chem-factory/internal/database/sqlite"
	"context"
	"fmt"
)

type UserRepo struct{ db *sqlite.Database }

func NewUserRepo(db *sqlite.Database) *UserRepo { return &UserRepo{db: db} }

func (r *UserRepo) FindByUsername(ctx context.Context, username string) (domain.User, error) {
	var user domain.User

	err := r.db.Extract(ctx).QueryRowContext(ctx,
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

	err := r.db.Extract(ctx).QueryRowContext(ctx,
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

	err := r.db.Extract(ctx).QueryRowContext(ctx,
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

	err := r.db.Extract(ctx).QueryRowContext(ctx,
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

	err := r.db.Extract(ctx).QueryRowContext(ctx,
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

	err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT username FROM user WHERE id = ?`,
		id,
	).Scan(&username)

	if err != nil {
		return "", fmt.Errorf("find username by id: %w", err)
	}

	return username, nil
}

func (r *UserRepo) Add(ctx context.Context, user domain.User) error {
	_, err := r.db.Extract(ctx).ExecContext(ctx,
		`INSERT INTO user(username, password) VALUES (?, ?)`,
		user.Username,
		user.Password,
	)

	if err != nil {
		return fmt.Errorf("add user: %w", err)
	}

	return nil
}

func (r *UserRepo) IncreaseBalanceByID(ctx context.Context, id, amount int) error {
	_, err := r.db.Extract(ctx).ExecContext(ctx,
		`UPDATE user SET balance = balance + ? WHERE id = ?`,
		amount,
		id,
	)
	if err != nil {
		return fmt.Errorf("increase balance by id: %w", err)
	}
	return nil
}

func (r *UserRepo) ReduceBalanceByID(ctx context.Context, id, amount int) error {
	_, err := r.db.Extract(ctx).ExecContext(ctx,
		`UPDATE user SET balance = balance - ? WHERE id = ?`,
		amount,
		id,
	)
	if err != nil {
		return fmt.Errorf("reduce balance by id: %w", err)
	}
	return nil
}

func (r *UserRepo) FindBalanceByID(ctx context.Context, id int) (int, error) {
	var balance int
	err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT balance FROM user WHERE id = ?`,
		id,
	).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("find balance by id: %w", err)
	}
	return balance, nil
}

func (r *UserRepo) UpdateBalanceByID(ctx context.Context, id, balance int) error {
	_, err := r.db.Extract(ctx).ExecContext(ctx,
		`UPDATE user SET balance = ? WHERE id = ?`,
		balance,
		id,
	)
	if err != nil {
		return fmt.Errorf("update balance by id: %w", err)
	}
	return nil
}

func (r *UserRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var exists int

	err := r.db.Extract(ctx).QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM user WHERE username = ?)`,
		username,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check username existence: %w", err)
	}

	return exists == 1, nil
}

func (r *UserRepo) UpdateXPAndLevelByID(ctx context.Context, id, xp, level int) error {
	_, err := r.db.Extract(ctx).ExecContext(ctx,
		`UPDATE user SET xp = ?, level = ? WHERE id = ?`,
		xp,
		level,
		id,
	)
	if err != nil {
		return fmt.Errorf("update xp and level by id: %w", err)
	}
	return nil
}

func (r *UserRepo) ListTopByXP(ctx context.Context, limit int) ([]domain.User, error) {
	if limit <= 0 {
		limit = 10
	}

	rows, err := r.db.Extract(ctx).QueryContext(ctx,
		`SELECT id, username, password, balance, xp, level FROM user ORDER BY xp DESC, level DESC, username ASC LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("list top users by xp: %w", err)
	}
	defer rows.Close()

	users := make([]domain.User, 0, limit)
	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Password,
			&user.Balance,
			&user.XP,
			&user.Level,
		)
		if err != nil {
			return nil, fmt.Errorf("scan top users by xp: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate top users by xp: %w", err)
	}

	return users, nil
}
