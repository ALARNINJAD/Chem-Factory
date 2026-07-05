package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type JWTManager interface {
	Generate(username string, id int) (string, error)
	Verify(token string) (string, int, error)
}

func Auth(jwt JWTManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}
		username, id, err := jwt.Verify(token)
		if err != nil || id == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		ctx.Set("username", username)
		ctx.Set("id", id)
		ctx.Next()
	}
}
