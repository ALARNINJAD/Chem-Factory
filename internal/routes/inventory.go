package routes

import (
	e "chem-factory/internal/dto/error"
	i "chem-factory/internal/dto/inventory"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getInventory(context *gin.Context) {

	request := i.InventoryExportRequest{Token: context.Request.Header.Get("Authorization")}

	if request.Token == "" {
		context.JSON(http.StatusBadRequest, e.ErrorResponse{Error: "Could not bind jwt."})
		return
	}

	response, err := route.service.ExportUserInventory(request)
	if err != nil {
		context.JSON(http.StatusForbidden, e.ErrorResponse{Error: err.Error()})
	}

	context.JSON(http.StatusOK, response)
}
