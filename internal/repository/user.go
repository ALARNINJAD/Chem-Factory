package repository

import (
	"chem-factory/internal/model"
)

type user struct {
	model.User
}

func (u *user) save() error {
	_, err := db.Exec("INSERT INTO user(username, password) VALUES (?, ?)", u.Username, u.Password)
	return err
}

func (u *user) password() (string, error) {
	var savedPassword string
	err := db.QueryRow("SELECT password FROM user WHERE username = ?", u.Username).Scan(&savedPassword)
	return savedPassword, err
}
