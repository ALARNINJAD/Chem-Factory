package routes

import (
	"chem-factory/internal/auth"
	"chem-factory/internal/model"
	"chem-factory/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(server *gin.Engine) {
	server.POST("/", root)
	server.POST("/login", login)
	server.POST("/register", register)
}

func root(context *gin.Context) {

}

func login(context *gin.Context) {

	var user model.User

	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if !auth.CheckPassword(&user) {
		context.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	token, err := auth.GenerateJWT(user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	context.JSON(http.StatusOK, gin.H{"token": token})
}

func register(context *gin.Context) {

	var user model.User

	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{})
	}

	if err := auth.HashPassword(&user); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{})
	}

	if err := repository.Save(user); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{})
	}

	context.JSON(http.StatusCreated, gin.H{})
}
