package http

import (
	authBootstrap "chem-factory/internal/modules/auth/bootstrap"
	inventoryBootstrap "chem-factory/internal/modules/inventory/bootstrap"
	userBootstrap "chem-factory/internal/modules/user/bootstrap"
	"chem-factory/internal/routes/http/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine
}

func NewServer(
	jwt middleware.JWTManager,
	userModule userBootstrap.Module,
	authModule authBootstrap.Module,
	inventoryModule inventoryBootstrap.Module,
) *Server {
	log.Println("Grouping routes")

	gin.SetMode(gin.DebugMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())

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
	{
		inventory := api.Group("/inventory")
		inventory.Use(middleware.Auth(jwt))
		{
			inventory.GET("/export", inventoryModule.InventoryHandler.Export)
		}
	}

	return &Server{engine: engine}
}

func (s *Server) Start(addr string) {
	log.Println("Starting server")
	log.Fatal(s.engine.Run(addr))
}
