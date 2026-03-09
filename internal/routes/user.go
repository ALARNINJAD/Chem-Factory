package routes

import (
	"chem-factory/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getUserData(context *gin.Context) {

	token := context.Request.Header.Get("Authorization")
	if token == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Could not bind jwt."})
		return
	}

	user, err := route.service.UserData(token)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Token is invalid."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"user": user})
}

func login(context *gin.Context) {

	var user model.User
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Could not bind json."})
		return
	}

	token, err := route.service.Login(user.Username, user.Password)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Could not login."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"token": token})
}

func register(context *gin.Context) {

	var user model.User
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Could not bind json."})
		return
	}

	if err := route.service.Register(user.Username, user.Password); err != nil {
		context.JSON(http.StatusForbidden, gin.H{"error": "Could not register user."})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"massage": "Done."})
}
