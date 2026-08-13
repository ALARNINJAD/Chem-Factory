package reedam

import (
	"chem-factory/pkg/lang"
	"net/http"
)

func UserNotFound(err error) *Reedam {
	return New().WithError(err).WithErrName(lang.ErrorUserNotFound).WithMessage(lang.MessageUserNotFound).WithStatus(http.StatusNotFound)
}

func Unexpected(err error) *Reedam {
	return New().WithError(err).WithErrName(lang.ErrorUnexpected).WithMessage(lang.MessageUnexpected).WithStatus(http.StatusInternalServerError)
}

func InvalidToken(err error) *Reedam {
	return New().WithError(err).WithErrName(lang.ErrorInvalidToken).WithStatus(http.StatusUnauthorized)
}

func InventoryNotFound(err error) *Reedam {
	return New().WithError(err).WithErrName(lang.ErrorInventoryNotFound).WithMessage(lang.MessageInventoryNotFound).WithStatus(http.StatusNotFound)
}

func MaterialNotFound(err error) *Reedam {
	return New().WithError(err).WithErrName(lang.ErrorMaterialNotFound).WithMessage(lang.MessageMaterialNotFound).WithStatus(http.StatusNotFound)
}

func MarketNotFound(err error) *Reedam {
	return New().WithError(err).WithErrName(lang.ErrorMarketNotFound).WithMessage(lang.MessageMarketNotFound).WithStatus(http.StatusNotFound)
}

func MixNotFound(err error) *Reedam {
	return New().WithError(err).WithErrName(lang.ErrorMixNotFound).WithMessage(lang.MessageMixNotFound).WithStatus(http.StatusNotFound)
}
