package handler

import (
	"chem-factory/internal/modules/mixer/adapter/http/dto"
	"chem-factory/internal/modules/mixer/core/port"
	"chem-factory/pkg/reedam"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MixerHandler struct {
	service port.MixerService
}

func NewMixerHandler(service port.MixerService) *MixerHandler {
	return &MixerHandler{service: service}
}

func (h *MixerHandler) Mixes(ctx *gin.Context) {

	userID := ctx.GetUint("user_id")

	response, err := h.service.Mixes(ctx.Request.Context(), userID)
	if err != nil {
		reedam.ErrorMessageHTTP(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *MixerHandler) Mix(ctx *gin.Context) {

	request := dto.MixRequest{}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.Mix(ctx, request, ctx.GetUint("user_id"))
	if err != nil {
		reedam.ErrorMessageHTTP(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *MixerHandler) Check(ctx *gin.Context) {

	request := dto.CheckRequest{}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.Check(ctx, request, ctx.GetUint("user_id"))
	if err != nil {
		reedam.ErrorMessageHTTP(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *MixerHandler) Pick(ctx *gin.Context) {

	request := dto.PickRequest{}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.Pick(ctx, request, ctx.GetUint("user_id"))
	if err != nil {
		reedam.ErrorMessageHTTP(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *MixerHandler) NewMaterial(ctx *gin.Context) {

	request := dto.NewMaterialRequest{}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.NewMaterial(ctx, request, ctx.GetUint("user_id"))
	if err != nil {
		reedam.ErrorMessageHTTP(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}
