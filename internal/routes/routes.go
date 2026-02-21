package routes

import "github.com/gin-gonic/gin"

func Register(server *gin.Engine) {
	server.POST("/", route)
	// server.POST("/login")
	// server.POST("/register")
}

func route(context *gin.Context) {

}
