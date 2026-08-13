package usecase

import (
	"chem-factory/internal/domain"
	"chem-factory/internal/modules/auth/adapter/http/dto"
	"chem-factory/internal/modules/auth/core/port"
	o "chem-factory/pkg/dto"
	"chem-factory/pkg/hash"
	"chem-factory/pkg/lang"
	"chem-factory/pkg/reedam"
	"context"
	"errors"
	"net/http"
)

type AuthUsecase struct {
	userRepo port.UserRepository
}

func NewAuthUsecase(userRepo port.UserRepository) *AuthUsecase {
	return &AuthUsecase{userRepo: userRepo}
}

func (service *AuthUsecase) Register(ctx context.Context, request dto.RegisterRequest) (dto.RegisterResponse, error) {

	id, _ := service.userRepo.FindIDByUsername(ctx, request.Username)
	if id != 0 {
		return dto.RegisterResponse{}, reedam.New().WithError(errors.New("user already exists")).WithErrName(lang.ErrorUserAlreadyExists).WithMessage(lang.MessageUserAlreadyExists).WithStatus(http.StatusForbidden)
	}

	var err error
	request.Password, err = hash.New(10).HashPassword(request.Password) // cost of hash must be in config
	if err != nil {
		return dto.RegisterResponse{}, err
	}

	user := domain.User{}
	if err = user.New(request.Username, request.Password); err != nil {
		return dto.RegisterResponse{}, err
	}

	if err = service.userRepo.Add(ctx, user); err != nil {
		return dto.RegisterResponse{}, err
	}

	return dto.RegisterResponse{MessageResponse: o.MessageResponse{Message: lang.MessageRegisterSuccessfully}}, nil
}

func (service *AuthUsecase) Login(ctx context.Context, request dto.LoginRequest) (dto.LoginResponse, error) {

	userID, err := service.userRepo.FindIDByUsername(ctx, request.Username)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	password, err := service.userRepo.FindPasswordByID(ctx, userID)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	if err := hash.New(10).CheckPassword(request.Password, password); err != nil { // cost of hash must be in config
		return dto.LoginResponse{}, err
	}

	return dto.LoginResponse{UserID: userID, MessageResponse: o.MessageResponse{Message: lang.MessageLoginSuccessfully}}, nil
}
