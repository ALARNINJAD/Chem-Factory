package routes

import (
	e "chem-factory/internal/dto/error"
	i "chem-factory/internal/dto/inventory"
	"net/http"

	"github.com/gin-gonic/gin"
)

func postInventoryItems(context *gin.Context) {

	request := i.InventoryAddRequest{Token: context.Request.Header.Get("Authorization")}

	if request.Token == "" {
		context.JSON(http.StatusBadRequest, e.ErrorResponse{Error: "Could not bind jwt."})
		return
	}

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, e.ErrorResponse{Error: err.Error()})
		return
	}

	response, err := route.service.AddToInventory(request)
	if err != nil {
		context.JSON(http.StatusForbidden, e.ErrorResponse{Error: err.Error()})
		return
	}

	context.JSON(http.StatusBadRequest, response)
}

func deleteInventoryItems(context *gin.Context) {

	request := i.InventoryAddRequest{Token: context.Request.Header.Get("Authorization")}

	if request.Token == "" {
		context.JSON(http.StatusBadRequest, e.ErrorResponse{Error: "Could not bind jwt."})
		return
	}

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, e.ErrorResponse{Error: err.Error()})
		return
	}

	response, err := route.service.RemoveFromInventory(request)
	if err != nil {
		context.JSON(http.StatusForbidden, e.ErrorResponse{Error: err.Error()})
		return
	}

	context.JSON(http.StatusOK, response)
}

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
