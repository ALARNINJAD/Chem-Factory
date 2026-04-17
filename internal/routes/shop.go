package routes

import (
	e "chem-factory/internal/dto/error"
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
