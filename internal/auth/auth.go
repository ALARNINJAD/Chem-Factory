package auth

import (
	"chem-factory/internal/model"
	"chem-factory/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

func CheckPassword(user model.User) bool {
	savedHashedPassword, err := repository.ExportPassword(user)
	if err != nil {
		return false
	}
	err = bcrypt.CompareHashAndPassword([]byte(savedHashedPassword), []byte(user.Password))
	return err == nil
}
