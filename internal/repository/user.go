package repository

import "fmt"

type user struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Balance  int    `json:"balance"`
	XP       int    `json:"xp"`
	Level    int    `json:"level"`
}

func createUserTable(r *repositoryManager) {
	query := `
	CREATE TABLE IF NOT EXISTS user (
		id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		balance INTEGER DEFAULT 200 CHECK("balance" >= 0),
		xp INTEGER DEFAULT 0 CHECK("xp" >= 0),
		level INTEGER DEFAULT 0 CHECK("level" >= 0)
	)`
	if _, err := r.db.Exec(query); err != nil {
		panic(fmt.Errorf("Repository user, create table: %w ", err))
	}
}

func (r *repositoryManager) FindUserByUsername(username string) (*user, error) {

	var u user

	err := r.db.QueryRow(`
		SELECT id,username,password,balance,xp,level
		FROM user WHERE username = ?`, username).Scan(
		&u.ID, &u.Username, &u.Password, &u.Balance, &u.XP, &u.Level)
	if err != nil {
		return nil, fmt.Errorf("Repository user, find user by username: %w ", err)
	}

	return &u, nil
}

func (r *repositoryManager) FindUserByID(id int) (*user, error) {

	var u user

	err := r.db.QueryRow(`
		SELECT id,username,password,balance,xp,level
		FROM user WHERE id = ?`, id).Scan(
		&u.ID, &u.Username, &u.Password, &u.Balance, &u.XP, &u.Level)
	if err != nil {
		return nil, fmt.Errorf("Repository user, find user by id: %w ", err)
	}

	return &u, nil
}

func (r *repositoryManager) FindPasswordByID(id int) (string, error) {

	var password string
	err := r.db.QueryRow("SELECT password FROM user WHERE id = ?", id).Scan(&password)
	if err != nil {
		return "", fmt.Errorf("Repository uesr, find password by id: %w ", err)
	}
	return password, nil
}

func (r *repositoryManager) FindPasswordByUsername(username string) (string, error) {

	var password string
	err := r.db.QueryRow("SELECT password FROM user WHERE username = ?", username).Scan(&password)
	if err != nil {
		return "", fmt.Errorf("Repository user, find password by username: %w ", err)
	}
	return password, nil
}

func (r *repositoryManager) FindIDbyUsername(username string) (int, error) {

	var id int
	err := r.db.QueryRow("SELECT id FROM user WHERE username = ?", username).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("Repository user, find id by username: %w ", err)
	}
	return id, nil
}

func (r *repositoryManager) FindUsernameByID(id int) (string, error) {

	var username string
	err := r.db.QueryRow("SELECT username FROM user WHERE id = ?", id).Scan(&username)
	if err != nil {
		return "", fmt.Errorf("Repository user, find username by id: %w ", err)
	}
	return username, nil
}

func (r *repositoryManager) SaveUser(username string, password string) error {

	_, err := r.db.Exec("INSERT INTO user(username, password) VALUES (?, ?)", username, password)
	if err != nil {
		return fmt.Errorf("Repository user, save user: %w ", err)
	}
	return nil
}

func (r *repositoryManager) IncreaseBalance(username string, amount int) error {

	_, err := r.db.Exec("UPDATE user SET balance = balance + ? WHERE username = ?", amount, username)
	if err != nil {
		return fmt.Errorf("Repository user, increase balance: %w ", err)
	}

	return nil
}

func (r *repositoryManager) ReduceBalance(username string, amount int) error {

	_, err := r.db.Exec("UPDATE user SET balance = balance - ? WHERE username = ?", amount, username)
	if err != nil {
		return fmt.Errorf("Repository user, reduce balance: %w ", err)
	}

	return nil
}
