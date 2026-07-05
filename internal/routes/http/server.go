package http

import (
	userBootstrap "chem-factory/internal/modules/user/bootstrap"
	authBootstrap "chem-factory/internal/modules/auth/bootstrap"
	"chem-factory/internal/routes/http/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine
}

func NewServer(jwt middleware.JWTManager, userModule userBootstrap.Module, authModule authBootstrap.Module) *Server {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())

	api := engine.Group("/api")
	{
		user := api.Group("/user")
		user.Use(middleware.Auth(jwt))
		{
			user.GET("/profile", userModule.UserHandler.GetProfile)
		}
	}
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authModule.AuthHandler.Register)
			auth.POST("/login", authModule.AuthHandler.Login)
		}
	}

	return &Server{engine: engine}
}

func (s *Server) Start(addr string) {
	log.Fatal(s.engine.Run(addr))
}
