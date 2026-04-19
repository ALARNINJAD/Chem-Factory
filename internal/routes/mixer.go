package routes

import (
	e "chem-factory/internal/dto/error"
	m "chem-factory/internal/dto/mixer"
	"net/http"

	"github.com/gin-gonic/gin"
)

func postAddToMixer(context *gin.Context) {

	request := m.MixerAddRequest{Token: context.Request.Header.Get("Authorization")}

	if request.Token == "" {
		context.JSON(http.StatusBadRequest, e.ErrorResponse{Error: "Invalid token."})
		return
	}

	var err error
	if err = context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, e.ErrorResponse{Error: err.Error()})
		return
	}

	var id int
	if id, err = route.service.AddToMixer(request); err != nil {
		context.JSON(http.StatusBadRequest, e.ErrorResponse{Error: err.Error()})
		return
	}

	context.JSON(http.StatusOK, m.MixerAddResponse{ID: id})
}
