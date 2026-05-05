package routes

import (
	e "chem-factory/internal/dto/error"
	u "chem-factory/internal/dto/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (route *Manager) getUserData(context *gin.Context) {

	request := u.UserDataRequest{Token: context.Request.Header.Get("Authorization")}

	if request.Token == "" {
		context.JSON(http.StatusBadRequest, e.ErrorResponse{Error: "Could not bind jwt."})
		return
	}

	response, err := route.service.UserData(request.Token)
	if err != nil {
		context.JSON(http.StatusUnauthorized, e.ErrorResponse{Error: err.Error()})
		return
	}

	context.JSON(http.StatusOK, response)
}

func (route *Manager) login(context *gin.Context) {

	var request u.UserLoginRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, e.ErrorResponse{Error: err.Error()})
		return
	}

	t, err := route.service.Login(request)
	if err != nil {
		context.JSON(http.StatusUnauthorized, e.ErrorResponse{Error: err.Error()})
		return
	}

	context.JSON(http.StatusOK, u.UserLoginResponse{Token: t})
}

func (route *Manager) register(context *gin.Context) {

	var request u.UserRegisterRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, e.ErrorResponse{Error: err.Error()})
		return
	}

	if err := route.service.Register(request); err != nil {
		context.JSON(http.StatusForbidden, e.ErrorResponse{Error: err.Error()})
		return
	}

	context.JSON(http.StatusCreated, u.UserRegisterResponse{Massage: "Done."})
}
