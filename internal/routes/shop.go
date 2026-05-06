package routes

import (
	"chem-factory/internal/dto/error"
	"chem-factory/internal/dto/massage"
	"chem-factory/internal/dto/shop"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (route *Manager) getShopItems(context *gin.Context) {

	response, err := route.service.ItemsForSell()
	if err != nil {
		context.JSON(http.StatusInternalServerError, error.Response{Error: err.Error()})
		return
	}

	context.JSON(http.StatusOK, response)
}

func (route *Manager) getBuyShopItems(context *gin.Context) {

	request := shop.BuyRequest{Token: context.Request.Header.Get("Authorization")}

	if request.Token == "" {
		context.JSON(http.StatusBadRequest, error.Response{Error: "Invalid token."})
		return
	}

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, error.Response{Error: err.Error()})
		return
	}

	if err := route.service.BuyMaterial(request); err != nil {
		context.JSON(http.StatusInternalServerError, error.Response{Error: err.Error()})
		return
	}

	context.JSON(http.StatusOK, massage.Response{Massage: "Done."})
}

func (route *Manager) postSetForSell(context *gin.Context) {

	request := shop.SetForSellRequest{Token: context.Request.Header.Get("Authorization")}

	if request.Token == "" {
		context.JSON(http.StatusBadRequest, error.Response{Error: "Invalid token."})
		return
	}

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, error.Response{Error: err.Error()})
		return
	}

	if err := route.service.SetForSell(request); err != nil {
		context.JSON(http.StatusBadRequest, error.Response{Error: err.Error()})
		return
	}

	context.JSON(http.StatusOK, massage.Response{Massage: "Done."})
}
