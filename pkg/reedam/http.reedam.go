package reedam

import (
	"chem-factory/pkg/dto"
	"chem-factory/pkg/lang"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorMessageHTTP(ctx *gin.Context, err error) {
	if err == nil {
		return
	}
	if r := As(err); r != nil {
		if r.Status == 0 {
			r.Status = http.StatusInternalServerError
			r.ErrName = lang.ErrorUnexpected
			r.Message = lang.MessageUnexpected
		}
		if r.ErrName == "" {
			r.ErrName = r.Err.Error()
		}
		type response struct {
			dto.ErrorResponse
			dto.MessageResponse
		}
		ctx.JSON(r.Status, response{
			ErrorResponse:   dto.ErrorResponse{Error: r.ErrName},
			MessageResponse: dto.MessageResponse{Message: r.Message},
		})
		return
	}
	ErrorMessageHTTP(ctx, Unexpected(err).WithLog())
}
