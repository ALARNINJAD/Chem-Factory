package routes

import "github.com/gin-gonic/gin"

func Register(server *gin.Engine) {
	server.POST("/", root)
	server.POST("/login", login)
	server.POST("/register", login)
}

func root(context *gin.Context) {

}

func login(context *gin.Context) {

}

func register(context *gin.Context) {

}
