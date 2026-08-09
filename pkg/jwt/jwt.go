package jwt

import (
	"chem-factory/pkg/reedam"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWT struct{ secretKey string }

func NewJWT(secretKey string) *JWT { return &JWT{secretKey: secretKey} }

func (j JWT) Generate(username string, userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"user_id":  userID,
		"exp":      time.Now().Add(time.Hour).Unix(),
	})
	t, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", reedam.InternalError(err)
	}
	return t, nil
}

func (j JWT) Verify(token string) (string, uint, error) {

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Auth jwt, verify jwt, token method: func.")
		}
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return "", 0, reedam.InternalError(err)
	}

	if !parsedToken.Valid {
		return "", 0, reedam.InternalError(err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", 0, reedam.InternalError(err)
	}

	idFloat, ok := claims["user_id"].(float64)
	if !ok {
		return "", 0, reedam.InternalError(err)
	}

	username, ok := claims["username"].(string)
	if !ok {
		return "", 0, reedam.InternalError(err)
	}

	return username, uint(idFloat), nil
}
