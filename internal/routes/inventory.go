package routes

import (
	"chem-factory/internal/dto/error"
	"chem-factory/internal/dto/inventory"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (route *Manager) getInventory(context *gin.Context) {

	request := inventory.ExportRequest{Token: context.Request.Header.Get("Authorization")}

	if request.Token == "" {
		context.JSON(http.StatusBadRequest, error.Response{Error: "Could not bind jwt."})
		return
	}

	response, err := route.service.ExportUserInventory(request)
	if err != nil {
		context.JSON(http.StatusForbidden, error.Response{Error: err.Error()})
	}

	context.JSON(http.StatusOK, response)
}
