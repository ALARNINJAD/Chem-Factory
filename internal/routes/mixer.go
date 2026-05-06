package routes

import (
	"chem-factory/internal/dto/error"
	"chem-factory/internal/dto/massage"
	"chem-factory/internal/dto/mixer"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (route *Manager) postAddToMixer(context *gin.Context) {

	request := mixer.AddRequest{Token: context.Request.Header.Get("Authorization")}

	if request.Token == "" {
		context.JSON(http.StatusBadRequest, error.Response{Error: "Invalid token."})
		return
	}

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, error.Response{Error: err.Error()})
		return
	}

	id, err := route.service.AddToMixer(request)
	if err != nil {
		context.JSON(http.StatusBadRequest, error.Response{Error: err.Error()})
		return
	}

	context.JSON(http.StatusOK, mixer.MixerAddResponse{ID: id})
}

func (route *Manager) getCheckMixTime(context *gin.Context) {

	request := mixer.CheckRequest{Token: context.Request.Header.Get("Authorization")}
	if request.Token == "" {
		context.JSON(http.StatusBadRequest, error.Response{Error: "Invalid token."})
		return
	}

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, error.Response{Error: err.Error()})
		return
	}

	response, err := route.service.CheckMix(request)
	if err != nil {
		context.JSON(http.StatusForbidden, error.Response{Error: err.Error()})
		return
	}

	context.JSON(http.StatusOK, response)
}

func (route *Manager) patchPickMix(context *gin.Context) {

	request := mixer.PickRequest{Token: context.Request.Header.Get("Authorization")}
	if request.Token == "" {
		context.JSON(http.StatusBadRequest, error.Response{Error: "Invalid token."})
		return
	}

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, error.Response{Error: err.Error()})
		return
	}

	if err := route.service.PickMix(request); err != nil {
		context.JSON(http.StatusForbidden, error.Response{Error: err.Error()})
		return
	}

	context.JSON(http.StatusOK, massage.Response{Massage: "Done."})
}

func (route *Manager) patchPickNewMix(context *gin.Context) {

	request := mixer.PickNewRequest{Token: context.Request.Header.Get("Authorization")}
	if request.Token == "" {
		context.JSON(http.StatusBadRequest, error.Response{Error: "Invalid token."})
		return
	}

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, error.Response{Error: err.Error()})
		return
	}

	if err := route.service.PickNewMix(request); err != nil {
		context.JSON(http.StatusForbidden, error.Response{Error: err.Error()})
		return
	}

	context.JSON(http.StatusOK, massage.Response{Massage: "Done."})
}
