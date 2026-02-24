package repository

import "chem-factory/internal/model"

type user struct {
	model.User
}

func (u *user) new() error {
	_, err := db.Exec("INSERT INTO user(username, password) VALUES (?, ?)", u.Username, u.Password)
	return err
}
