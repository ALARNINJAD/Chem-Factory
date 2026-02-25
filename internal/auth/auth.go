package auth

import (
	"chem-factory/internal/model"
	"chem-factory/internal/repository"
)

func CheckUserPassword(usr model.User) bool {
	savedPassword, err := repository.ExportUserPassword(usr)
	if err != nil {
		return false
	}
	return savedPassword == usr.Password
}
