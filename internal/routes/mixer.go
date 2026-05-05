package routes

import (
	e "chem-factory/internal/dto/error"
	"chem-factory/internal/dto/massage"
	"chem-factory/internal/dto/mixer"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (route *Manager) postAddToMixer(context *gin.Context) {

	request := mixer.MixerAddRequest{Token: context.Request.Header.Get("Authorization")}

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

	context.JSON(http.StatusOK, mixer.MixerAddResponse{ID: id})
}

func (route *Manager) getCheckMixTime(context *gin.Context) {

	request := mixer.CheckMixRequest{Token: context.Request.Header.Get("Authorization")}
	if request.Token == "" {
		context.JSON(http.StatusBadRequest, e.ErrorResponse{Error: "Invalid token."})
		return
	}

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, e.ErrorResponse{Error: err.Error()})
		return
	}

	response, err := route.service.CheckMix(request)
	if err != nil {
		context.JSON(http.StatusForbidden, e.ErrorResponse{Error: err.Error()})
		return
	}

	context.JSON(http.StatusOK, response)
}

func (route *Manager) patchPickMix(context *gin.Context) {

	request := mixer.PickMixRequest{Token: context.Request.Header.Get("Authorization")}
	if request.Token == "" {
		context.JSON(http.StatusBadRequest, e.ErrorResponse{Error: "Invalid token."})
		return
	}

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, e.ErrorResponse{Error: err.Error()})
		return
	}

	if err := route.service.PickMix(request); err != nil {
		context.JSON(http.StatusForbidden, e.ErrorResponse{Error: err.Error()})
		return
	}

	context.JSON(http.StatusOK, massage.MassageResponse{Massage: "Done."})
}

func (route *Manager) patchPickNewMix(context *gin.Context) {

	request := mixer.PickNewMixRequest{Token: context.Request.Header.Get("Authorization")}
	if request.Token == "" {
		context.JSON(http.StatusBadRequest, e.ErrorResponse{Error: "Invalid token."})
		return
	}

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, e.ErrorResponse{Error: err.Error()})
		return
	}

	if err := route.service.PickNewMix(request); err != nil {
		context.JSON(http.StatusForbidden, e.ErrorResponse{Error: err.Error()})
		return
	}

	context.JSON(http.StatusOK, massage.MassageResponse{Massage: "Done."})
}
