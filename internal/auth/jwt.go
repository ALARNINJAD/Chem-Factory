package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtManager struct{ secretKey string }

func NewJWTManager(secretKey string) *jwtManager { return &jwtManager{secretKey: secretKey} }

func (j *jwtManager) Generate(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour).Unix(),
	})
	t, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", fmt.Errorf("Auth jwt, generate jwt: %w", err)
	}
	return t, nil
}

func (j *jwtManager) Verify(token string) (string, error) {

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Auth jwt, verify jwt, token method: func.")
		}
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return "", fmt.Errorf("Auth jwt, verify jwt, jwt parse: %w", err)
	}

	if !parsedToken.Valid {
		return "", errors.New("Auth jwt, verify jwt, parse token validation: invalid token.")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("Auth jwt, verify jwt, token claims: not ok.")
	}

	return claims["username"].(string), nil
}
