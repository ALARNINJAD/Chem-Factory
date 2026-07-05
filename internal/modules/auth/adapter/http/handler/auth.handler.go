package handler

import (
	"chem-factory/internal/modules/auth/adapter/http/dto"
	"chem-factory/internal/modules/auth/core/port"
	"chem-factory/internal/routes/http/middleware"
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

	userID, err := h.service.Login(ctx.Request.Context(), request)
	if err != nil || userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := h.jwt.Generate(request.Username, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	ctx.JSON(http.StatusOK, dto.LoginResponse{Token: token})
}

func (h *AuthHandler) Register(ctx *gin.Context) {

	var request dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Register(ctx.Request.Context(), request); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Done."})
}
