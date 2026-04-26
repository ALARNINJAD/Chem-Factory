package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func (a *authManager) CheckPassword(password, hashedPassword string) error {

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("Auth hash, check password: %w ", err)
	}
	return nil
}

func (a *authManager) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), a.costOfHash)
	if err != nil {
		return "", fmt.Errorf("Auth hash, hash password: %w ", err)
	}
	return string(hashedPassword), nil
}
