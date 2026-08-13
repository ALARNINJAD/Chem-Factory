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
		return "", reedam.Unexpected(err).WithLog(username, userID)
	}
	return t, nil
}

func (j JWT) Verify(token string) (string, uint, error) {

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, reedam.InvalidToken(errors.New("token method is not *jwt.SigningMethodHMAC"))
		}
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return "", 0, reedam.InvalidToken(err)
	}

	if !parsedToken.Valid {
		return "", 0, reedam.InvalidToken(errors.New("parsed token is not valid"))
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", 0, reedam.InvalidToken(errors.New("parsed token claims is not jwt.MapClaims"))
	}

	idFloat, ok := claims["user_id"].(float64)
	if !ok {
		return "", 0, reedam.InvalidToken(errors.New("claims user_id is not float64"))
	}

	username, ok := claims["username"].(string)
	if !ok {
		return "", 0, reedam.InvalidToken(errors.New("claims username is not string"))
	}

	return username, uint(idFloat), nil
}
