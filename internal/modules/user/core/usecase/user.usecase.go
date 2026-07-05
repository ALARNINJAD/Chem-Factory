package usecase

import (
	"chem-factory/internal/modules/user/adapter/http/dto"
	"chem-factory/internal/modules/user/core/port"
	"context"
)

type userUsecase struct {
	userRepo   port.UserRepository
	transactor port.Transactor
}

func NewUserUsecase(userRepo port.UserRepository, transactor port.Transactor) *userUsecase{
	return &userUsecase{userRepo: userRepo, transactor: transactor}
}

func (service *userUsecase) GetProfile(ctx context.Context, userID int) (dto.ProfileResponse, error) {

	user, err := service.userRepo.FindByID(ctx, userID)
	if err != nil {
		return dto.ProfileResponse{}, err
	}

	return dto.ProfileResponse{
		Username: user.Username,
		Balance:  user.Balance,
		XP:       user.XP,
		Level:    user.Level,
	}, nil
}
