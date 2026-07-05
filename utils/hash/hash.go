package hash

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Hash struct{ costOfHash int }

func New(costOfHash int) Hash { return Hash{costOfHash: costOfHash} }

func (hash Hash) CheckPassword(password, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("check password failed: %w", err)
	}
	return nil
}

func (hash Hash) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), hash.costOfHash)
	if err != nil {
		return "", fmt.Errorf("could not hash password: %w", err)
	}
	return string(hashedPassword), nil
}
