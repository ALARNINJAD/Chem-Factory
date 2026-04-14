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
	// data, err := os.ReadFile(filepath.Join(".", "configs", "route.json"))
	// if err != nil {
	// 	panic("Could not access route config file")
	// }
	// var rt []*routesManager
	// err = json.Unmarshal(data, &rt)
	// if err != nil {
	// 	panic("Could not access route config")
	// }
	// log.Println(rt[0].port)
	return &routesManager{server: svr, service: svc}
}

func (r *routesManager) Start() {
	fmt.Println("Server Is Started")
	route = r
	registerRoutes()
	route.server.Run(":8090")
}

func registerRoutes() {
	route.server.POST("/", root)
	route.server.POST("/login", login)
	route.server.POST("/register", register)
	route.server.GET("/user", getUserData)
	// get inventory
	// get all products on sell
	// buy product
	// get mixing recipie
	// post new product
}

func root(context *gin.Context) {

}
