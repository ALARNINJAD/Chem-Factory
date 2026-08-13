package hash

import (
	"chem-factory/pkg/lang"
	"chem-factory/pkg/reedam"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type Hash struct{ costOfHash int }

func New(costOfHash int) Hash { return Hash{costOfHash: costOfHash} }

func (hash Hash) CheckPassword(password, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return reedam.New().WithError(err).WithErrName(lang.ErrorWrongPassword).WithMessage(lang.MessageWrongPassword).WithStatus(http.StatusUnauthorized)
	}
	return nil
}

func (hash Hash) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), hash.costOfHash)
	if err != nil {
		return "", reedam.Unexpected(err).WithLog()
	}
	return string(hashedPassword), nil
}
