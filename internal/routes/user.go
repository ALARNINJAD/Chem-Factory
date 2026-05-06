package routes

import (
	"chem-factory/internal/dto/error"
	"chem-factory/internal/dto/massage"
	"chem-factory/internal/dto/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (route *Manager) getUserData(context *gin.Context) {

	request := user.DataRequest{Token: context.Request.Header.Get("Authorization")}

	if request.Token == "" {
		context.JSON(http.StatusBadRequest, error.Response{Error: "Could not bind jwt."})
		return
	}

	response, err := route.service.UserData(request.Token)
	if err != nil {
		context.JSON(http.StatusUnauthorized, error.Response{Error: err.Error()})
		return
	}

	context.JSON(http.StatusOK, response)
}

func (route *Manager) login(context *gin.Context) {

	var request user.LoginRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, error.Response{Error: err.Error()})
		return
	}

	t, err := route.service.Login(request)
	if err != nil {
		context.JSON(http.StatusUnauthorized, error.Response{Error: err.Error()})
		return
	}

	context.JSON(http.StatusOK, user.LoginResponse{Token: t})
}

func (route *Manager) register(context *gin.Context) {

	var request user.RegisterRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, error.Response{Error: err.Error()})
		return
	}

	if err := route.service.Register(request); err != nil {
		context.JSON(http.StatusForbidden, error.Response{Error: err.Error()})
		return
	}

	context.JSON(http.StatusCreated, massage.Response{Massage: "Done."})
}
