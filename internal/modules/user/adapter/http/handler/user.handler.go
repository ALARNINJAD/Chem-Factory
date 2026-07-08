package handler

import (
	"chem-factory/internal/modules/user/core/port"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service port.UserService
}

func NewUserHandler(service port.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetProfile(ctx *gin.Context) {

	id := ctx.GetInt("id")

	response, err := h.service.GetProfile(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch user profile"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
