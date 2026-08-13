package handler

import (
	"chem-factory/internal/modules/auth/adapter/http/dto"
	"chem-factory/internal/modules/auth/core/port"
	"chem-factory/internal/routes/http/middleware"
	"chem-factory/pkg/reedam"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service port.AuthService
	jwt     middleware.JWTManager
}

func NewAuthHandler(service port.AuthService, jwt middleware.JWTManager) *AuthHandler {
	return &AuthHandler{service: service, jwt: jwt}
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var request dto.LoginRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.Login(ctx.Request.Context(), request)
	if err != nil {
		reedam.ErrorMessageHTTP(ctx, err)
		return
	}

	response.Token, err = h.jwt.Generate(request.Username, response.UserID)
	if err != nil {
		reedam.ErrorMessageHTTP(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *AuthHandler) Register(ctx *gin.Context) {

	var request dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.Register(ctx.Request.Context(), request)
	if err != nil {
		reedam.ErrorMessageHTTP(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, response)
}
