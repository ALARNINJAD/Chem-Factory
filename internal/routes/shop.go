package routes

import (
	e "chem-factory/internal/dto/error"
	s "chem-factory/internal/dto/shop"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getShopItems(context *gin.Context) {

	response, err := route.service.ItemsForSell()
	if err != nil {
		context.JSON(http.StatusInternalServerError, e.ErrorResponse{Error: err.Error()})
		return
	}

	context.JSON(http.StatusOK, response)
}

func getBuyShopItems(context *gin.Context) {

	request := s.ShopBuyRequest{Token: context.Request.Header.Get("Authorization")}

	if request.Token == "" {
		context.JSON(http.StatusBadRequest, e.ErrorResponse{Error: "Invalid token."})
		return
	}

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, e.ErrorResponse{Error: err.Error()})
		return
	}

	if err := route.service.Buy(request); err != nil {
		context.JSON(http.StatusInternalServerError, e.ErrorResponse{Error: err.Error()})
		return
	}

	context.JSON(http.StatusOK, s.ShopBuyResponse{Massage: "Done."})
}

func postSetForSell(context *gin.Context) {

	request := s.ShopSetForSellRequest{Token: context.Request.Header.Get("Authorization")}

	if request.Token == "" {
		context.JSON(http.StatusBadRequest, e.ErrorResponse{Error: "Invalid token."})
		return
	}

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, e.ErrorResponse{Error: err.Error()})
		return
	}

	if err := route.service.SetForSell(request); err != nil {
		context.JSON(http.StatusBadRequest, e.ErrorResponse{Error: err.Error()})
		return
	}

	context.JSON(http.StatusOK, s.ShopSetForSellResponse{Massage: "Done."})
}
