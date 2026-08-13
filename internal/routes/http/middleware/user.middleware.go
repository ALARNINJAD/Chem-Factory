package middleware

import (
	"chem-factory/pkg/dto"
	"chem-factory/pkg/lang"
	"net/http"

	"github.com/gin-gonic/gin"
)

type JWTManager interface {
	Generate(username string, userID uint) (string, error)
	Verify(token string) (string, uint, error)
}

func Auth(jwt JWTManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{Error: lang.ErrorMissingToken})
			return
		}
		username, userID, err := jwt.Verify(token)
		if err != nil || userID == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{Error: lang.ErrorInvalidToken})
			return
		}
		ctx.Set("username", username)
		ctx.Set("user_id", userID)
		ctx.Next()
	}
}
