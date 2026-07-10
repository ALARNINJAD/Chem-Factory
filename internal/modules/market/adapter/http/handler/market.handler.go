package handler

import (
	"chem-factory/internal/modules/market/adapter/http/dto"
	"chem-factory/internal/modules/market/core/port"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MarketHandler struct {
	service port.MarketService
}

func NewMarketHandler(service port.MarketService) *MarketHandler {
	return &MarketHandler{service: service}
}

func (h *MarketHandler) Export(ctx *gin.Context) {

	response, err := h.service.Export(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch market list"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *MarketHandler) Buy(ctx *gin.Context) {

	request := dto.BuyRequest{}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.Buy(ctx, request, ctx.GetInt("user_id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not complete purchase"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "purchase completed successfully"})
}

func (h *MarketHandler) SetForSell(ctx *gin.Context) {

	request := dto.SetForSellRequest{}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.SetForSell(ctx, request, ctx.GetInt("user_id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not set item for sale"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "item set for sale successfully"})
}
