package routes

import (
	"chem-factory/internal/service"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

type Manager struct {
	server  *gin.Engine
	service *service.Manager
}

func New(service *service.Manager) *Manager {
	return &Manager{server: gin.Default(), service: service}
}

func (route *Manager) Start() {
	route.registerRoutes()
	route.server.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}

func (route *Manager) registerRoutes() {
	route.server.POST("/login", route.login)
	route.server.POST("/register", route.register)
	route.server.GET("/user", route.getUserData)
	route.server.GET("/inventory", route.getInventory)
	route.server.GET("/shop", route.getShopItems)
	route.server.POST("/shop/buy", route.getBuyShopItems)
	route.server.POST("/shop", route.postSetForSell)
	route.server.POST("/mixer", route.postAddToMixer)
	route.server.GET("/mixer", route.getCheckMixTime)
	route.server.PATCH("/mixer", route.patchPickMix)
	route.server.PATCH("/mixer/new", route.patchPickNewMix)
}
