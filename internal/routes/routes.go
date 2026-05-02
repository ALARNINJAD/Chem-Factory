package routes

import (
	"chem-factory/internal/service"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

var route *Manager

type Manager struct {
	server  *gin.Engine
	service *service.Manager
}

func New(svr *gin.Engine, svc *service.Manager) *Manager {
	return &Manager{server: svr, service: svc}
}

func (r *Manager) Start() {
	route = r
	registerRoutes()
	route.server.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}

func registerRoutes() {
	route.server.POST("/login", login)
	route.server.POST("/register", register)
	route.server.GET("/user", getUserData)
	route.server.GET("/inventory", getInventory)
	route.server.GET("/shop", getShopItems)
	route.server.POST("/shop/buy", getBuyShopItems)
	route.server.POST("/shop", postSetForSell)
	route.server.POST("/mixer", postAddToMixer)
	route.server.GET("/mixer", getCheckMixTime)
	route.server.PATCH("/mixer", patchPickMix)
	route.server.PATCH("/mixer/new", patchPickNewMix)
}
