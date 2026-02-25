package auth

import (
	"chem-factory/internal/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const holyMolySuperSecretKey = "JAVID_SHAH"

func GenerateJWT(user model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour).Unix(),
	})
	return token.SignedString([]byte(holyMolySuperSecretKey))
}
