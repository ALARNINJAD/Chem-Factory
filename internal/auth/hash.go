package auth

import (
	"chem-factory/internal/model"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(user *model.User) error {
	HashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 15)
	user.Password = string(HashedPassword)
	return err
}
