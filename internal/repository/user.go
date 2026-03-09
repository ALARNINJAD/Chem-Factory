package repository

import "chem-factory/internal/model"

func (r *repositoryManager) ExportUserByUsername(u *model.User) error {
	return r.db.QueryRow("SELECT id,password,balance,xp,level FROM user WHERE username = ?", u.Username).Scan(
		&u.ID, &u.Password, &u.Balance, &u.XP, &u.Level)
}

func (r *repositoryManager) ExportUserByID(u *model.User) error {
	return r.db.QueryRow("SELECT username,password,balance,xp,level FROM user WHERE id = ?", u.ID).Scan(
		&u.Username, &u.Password, &u.Balance, &u.XP, &u.Level)
}

func (r *repositoryManager) ExportPasswordByID(id int) (string, error) {
	var savedPassword string
	err := r.db.QueryRow("SELECT password FROM user WHERE id = ?", id).Scan(&savedPassword)
	return savedPassword, err
}

func (r *repositoryManager) ExportPasswordByUsername(username string) (string, error) {
	var savedPassword string
	err := r.db.QueryRow("SELECT password FROM user WHERE username = ?", username).Scan(&savedPassword)
	return savedPassword, err
}

func (r *repositoryManager) ExportIDbyUsername(username string) (int, error) {
	var id int
	err := r.db.QueryRow("SELECT id FROM user WHERE username = ?", username).Scan(&id)
	return id, err
}

func (r *repositoryManager) ExportUsernameByID(id int) (string, error) {
	var username string
	err := r.db.QueryRow("SELECT username FROM user WHERE id = ?", id).Scan(&username)
	return username, err
}

func (r *repositoryManager) SaveNewUser(u model.User) error {
	_, err := r.db.Exec("INSERT INTO user(username, password) VALUES (?, ?)", u.Username, u.Password)
	return err
}
