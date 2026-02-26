package auth

import (
	"chem-factory/internal/model"
	"chem-factory/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

func CheckPassword(user *model.User) bool {
	password := user.Password
	if err := repository.Export(user); err != nil {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func HashPassword(user *model.User) error {
	HashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 15)
	user.Password = string(HashedPassword)
	return err
}
