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
	server.GET("/user", getUserData)
	// get inventory
	// get all products on sell
	// buy product
	// get mixing recipie
	// post new product
}

func root(context *gin.Context) {

}

func getUserData(context *gin.Context) {

	var user model.User

	token := context.Request.Header.Get("Authorization")
	if token == "" {
		context.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	var err error
	if user.Username, err = auth.VerifyJWT(token); err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	if err = repository.Export(&user); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	user.ID, user.Password = 0, ""
	context.JSON(http.StatusOK, gin.H{"user": user})
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
		return
	}

	if err := auth.HashPassword(&user); err != nil {
		context.JSON(http.StatusForbidden, gin.H{})
		return
	}

	if err := repository.Save(&user); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	context.JSON(http.StatusCreated, gin.H{})
}
