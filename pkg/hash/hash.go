package hash

import (
	"chem-factory/pkg/lang"
	"chem-factory/pkg/reedam"

	"golang.org/x/crypto/bcrypt"
)

type Hash struct{ costOfHash int }

func New(costOfHash int) Hash { return Hash{costOfHash: costOfHash} }

func (hash Hash) CheckPassword(password, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return reedam.New().WithError(err).WithMessage(lang.ErrorUnexpected).WithStatus(reedam.StatusInternalServerError).WithLog()
	}
	return nil
}

func (hash Hash) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), hash.costOfHash)
	if err != nil {
		return "", reedam.New().WithError(err).WithMessage(lang.ErrorUnexpected).WithStatus(reedam.StatusInternalServerError).WithLog()
	}
	return string(hashedPassword), nil
}
