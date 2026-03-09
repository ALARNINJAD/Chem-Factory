package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (auth *authManager) GenerateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour).Unix(),
	})
	return token.SignedString([]byte(auth.secretKey))
}

func (auth *authManager) VerifyJWT(token string) (string, error) {

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid JWT.")
		}
		return []byte(auth.secretKey), nil
	})
	if err != nil {
		return "", err
	}

	if !parsedToken.Valid {
		return "", errors.New("Invalid JWT.")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("Invalid JWT.")
	}

	return claims["username"].(string), nil
}
