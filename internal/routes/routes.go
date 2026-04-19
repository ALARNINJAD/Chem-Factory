package routes

import (
	"chem-factory/internal/service"
	"fmt"

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
	fmt.Println("Server Is Started")
	route = r
	registerRoutes()
	route.server.Run(":8090")
}

func registerRoutes() {
	// route.server.POST("/", root)
	route.server.POST("/login", login)
	route.server.POST("/register", register)
	route.server.GET("/user", getUserData)
	// route.server.POST("/inventory", postInventoryItems)
	route.server.DELETE("/inventory", deleteInventoryItems)
	route.server.GET("/inventory", getInventory)
	route.server.GET("/shop", getShopItems)
	route.server.POST("/shop/buy", getBuyShopItems)
	route.server.POST("/shop", postSetForSell)
	route.server.POST("/mixer", postAddToMixer)

	// get inventory
	// get all products on sell
	// buy product
	// get mixing recipie
	// post new product
}
