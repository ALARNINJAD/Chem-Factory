package repository

import (
	"database/sql"
	"fmt"
)

type user struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Balance  int    `json:"balance"`
	XP       int    `json:"xp"`
	Level    int    `json:"level"`
}

type userManager struct{ db *sql.DB }

func NewUserManager(db *sql.DB) *userManager { return &userManager{db: db} }

func (u *userManager) FindByUsername(username string) (*user, error) {
	var usr user
	err := u.db.QueryRow(`
		SELECT id, username, password, balance, xp, level
		FROM user WHERE username = ?`, username).Scan(
		&usr.ID, &usr.Username, &usr.Password, &usr.Balance, &usr.XP, &usr.Level)
	if err != nil {
		return nil, fmt.Errorf("Repository user, find user by username: %w", err)
	}
	return &usr, nil
}

func (u *userManager) FindByID(id int) (*user, error) {
	var usr user
	err := u.db.QueryRow(`
		SELECT id,username,password,balance,xp,level
		FROM user WHERE id = ?`, id).Scan(
		&usr.ID, &usr.Username, &usr.Password, &usr.Balance, &usr.XP, &usr.Level)
	if err != nil {
		return nil, fmt.Errorf("Repository user, find user by id: %w", err)
	}
	return &usr, nil
}

func (u *userManager) FindPasswordByID(id int) (string, error) {
	var password string
	err := u.db.QueryRow("SELECT password FROM user WHERE id = ?", id).Scan(&password)
	if err != nil {
		return "", fmt.Errorf("Repository uesr, find password by id: %w", err)
	}
	return password, nil
}

func (u *userManager) FindPasswordByUsername(username string) (string, error) {
	var password string
	err := u.db.QueryRow("SELECT password FROM user WHERE username = ?", username).Scan(&password)
	if err != nil {
		return "", fmt.Errorf("Repository user, find password by username: %w", err)
	}
	return password, nil
}

func (u *userManager) FindIDbyUsername(username string) (int, error) {
	var id int
	err := u.db.QueryRow("SELECT id FROM user WHERE username = ?", username).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("Repository user, find id by username: %w", err)
	}
	return id, nil
}

func (u *userManager) FindUsernameByID(id int) (string, error) {
	var username string
	err := u.db.QueryRow("SELECT username FROM user WHERE id = ?", id).Scan(&username)
	if err != nil {
		return "", fmt.Errorf("Repository user, find username by id: %w", err)
	}
	return username, nil
}

func (u *userManager) Add(tx *sql.Tx, username, password string) error {
	_, err := tx.Exec("INSERT INTO user(username, password) VALUES (?, ?)", username, password)
	if err != nil {
		return fmt.Errorf("Repository user, save user: %w", err)
	}
	return nil
}

func (u *userManager) IncreaseBalance(tx *sql.Tx, username string, amount int) error {
	_, err := tx.Exec("UPDATE user SET balance = balance + ? WHERE username = ?", amount, username)
	if err != nil {
		return fmt.Errorf("Repository user, increase balance: %w", err)
	}
	return nil
}

func (u *userManager) ReduceBalance(tx *sql.Tx, username string, amount int) error {
	_, err := tx.Exec("UPDATE user SET balance = balance - ? WHERE username = ?", amount, username)
	if err != nil {
		return fmt.Errorf("Repository user, reduce balance: %w", err)
	}
	return nil
}
