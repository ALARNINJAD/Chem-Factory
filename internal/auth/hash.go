package auth

import "golang.org/x/crypto/bcrypt"

func (auth *authManager) CheckPassword(password, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

func (auth *authManager) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), auth.costOfHash)
	return string(hashedPassword), err
}
