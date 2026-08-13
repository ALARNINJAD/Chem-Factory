package domain

import (
	"chem-factory/pkg/lang"
	"chem-factory/pkg/reedam"
	"errors"
	"net/http"
)

type User struct {
	ID       uint   `json:"id,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
	Balance  int    `json:"balance"`
	XP       int    `json:"xp"`
	Level    int    `json:"level"`
}

func (user *User) New(username, password string) error {
	if username == "" {
		return reedam.New().WithError(errors.New("username is empty")).WithErrName(lang.ErrorInvalidUsername).WithMessage(lang.MessageInvalidUsername).WithStatus(http.StatusBadRequest)
	}
	if password == "" {
		return reedam.New().WithError(errors.New("password is empty")).WithErrName(lang.ErrorInvalidPassword).WithMessage(lang.MessageInvalidPassword).WithStatus(http.StatusBadRequest)
	}
	user.Username = username
	user.Password = password
	user.Balance = 500 // for test . must be inside config
	user.XP = 0
	user.Level = 1
	return nil
}

func (user *User) GetXP(xp int) error {
	if xp < 0 {
		return reedam.Unexpected(errors.New("xp is negative")).WithLog(xp)
	}
	user.XP += xp
	for user.XP >= user.Level*1000 {
		user.XP -= user.Level * 1000
		user.Level++
	}
	return nil
}
