package usecase

import (
	"chem-factory/internal/domain"
	"chem-factory/internal/modules/auth/adapter/http/dto"
	"chem-factory/internal/modules/auth/core/port"
	"chem-factory/pkg/hash"
	"context"
	"errors"
)

type AuthUsecase struct {
	userRepo   port.UserRepository
	transactor port.Transactor
}

func NewAuthUsecase(userRepo port.UserRepository, transactor port.Transactor) *AuthUsecase {
	return &AuthUsecase{userRepo: userRepo, transactor: transactor}
}

func (service *AuthUsecase) Register(ctx context.Context, request dto.RegisterRequest) error {

	id, _ := service.userRepo.FindIDByUsername(ctx, request.Username)
	if id != 0 {
		return errors.New("user already exists")
	}

	var err error
	request.Password, err = hash.New(10).HashPassword(request.Password) // cost of hash must be in config
	if err != nil {
		return err
	}

	user := domain.User{}
	if err = user.New(request.Username, request.Password); err != nil {
		return err
	}
	user.IncreaseBalance(500) // for test
	return service.userRepo.Add(ctx, user)
}

func (service *AuthUsecase) Login(ctx context.Context, request dto.LoginRequest) (uint, error) {

	userID, err := service.userRepo.FindIDByUsername(ctx, request.Username)
	if err != nil {
		return 0, err
	}

	password, err := service.userRepo.FindPasswordByID(ctx, userID)
	if err != nil {
		return 0, err
	}

	if err := hash.New(10).CheckPassword(request.Password, password); err != nil { // cost of hash must be in config
		return 0, err
	}

	return userID, nil
}
