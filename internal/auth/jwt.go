package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (a *authManager) GenerateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour).Unix(),
	})
	t, err := token.SignedString([]byte(a.secretKey))
	if err != nil {
		return "", fmt.Errorf("Auth jwt, generate jwt: %w", err)
	}
	return t, nil
}

func (a *authManager) VerifyJWT(token string) (string, error) {

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Auth jwt, verify jwt, token method: ")
		}
		return []byte(a.secretKey), nil
	})
	if err != nil {
		return "", fmt.Errorf("Auth jwt, verify jwt, jwt parse: %w", err)
	}

	if !parsedToken.Valid {
		return "", errors.New("Auth jwt, verify jwt, parse token validation: ")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("Auth jwt, verify jwt, token claims: ")
	}

	return claims["username"].(string), nil
}
