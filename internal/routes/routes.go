package routes

import (
	"chem-factory/internal/service"
	"log"

	"github.com/gin-gonic/gin"
)

var route *routesManager

type RoutesManager interface {
	Route()
}

type routesManager struct {
	port    string
	server  *gin.Engine
	service service.ServiceManager
}

func Init(svr *gin.Engine, svc service.ServiceManager) *routesManager {
	return &routesManager{server: svr, service: svc}
}

func (r *routesManager) Start() {
	log.Println("Server is started.")
	route = r
	registerRoutes()
	route.server.Run(":8090")
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
