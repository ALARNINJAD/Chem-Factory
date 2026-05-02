package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type hashManager struct{}

func NewHashManager() *hashManager { return &hashManager{} }

func (h *hashManager) CheckPassword(password, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("Auth hash, check password: %w", err)
	}
	return nil
}

func (h *hashManager) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 5)
	if err != nil {
		return "", fmt.Errorf("Auth hash, hash password: %w", err)
	}
	return string(hashedPassword), nil
}
