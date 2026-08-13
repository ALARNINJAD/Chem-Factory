package handler

import (
	"chem-factory/internal/modules/user/core/port"
	"chem-factory/pkg/reedam"
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

	userID := ctx.GetUint("user_id")

	response, err := h.service.GetProfile(ctx.Request.Context(), userID)
	if err != nil {
		reedam.ErrorMessageHTTP(ctx, err)
	}

	ctx.JSON(http.StatusOK, response)
}
