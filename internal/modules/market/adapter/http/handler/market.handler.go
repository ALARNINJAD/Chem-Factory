package handler

import (
	"chem-factory/internal/modules/market/adapter/http/dto"
	"chem-factory/internal/modules/market/core/port"
	"chem-factory/pkg/reedam"
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
		reedam.ErrorMessageHTTP(ctx, err)
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

	response, err := h.service.Buy(ctx, request, ctx.GetUint("user_id"))
	if err != nil {
		reedam.ErrorMessageHTTP(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *MarketHandler) SetForSell(ctx *gin.Context) {

	request := dto.SetForSellRequest{}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.SetForSell(ctx, request, ctx.GetUint("user_id"))
	if err != nil {
		reedam.ErrorMessageHTTP(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}
