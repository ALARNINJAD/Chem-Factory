package handler

import (
	"chem-factory/internal/modules/inventory/core/port"
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

	userID := ctx.GetInt("user_id")

	response, err := h.service.Export(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not export inventory"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
