package handler

import (
	"chem-factory/internal/modules/inventory/core/port"
	"chem-factory/pkg/reedam"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	service port.InventoryService
}

func NewInventoryHandler(service port.InventoryService) *InventoryHandler {
	return &InventoryHandler{service: service}
}

func (h *InventoryHandler) Export(ctx *gin.Context) {

	userID := ctx.GetUint("user_id")

	response, err := h.service.Export(ctx.Request.Context(), userID)
	if err != nil {
		reedam.ErrorMessageHTTP(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}
