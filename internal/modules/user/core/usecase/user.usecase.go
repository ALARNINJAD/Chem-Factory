package usecase

import (
	"chem-factory/internal/modules/user/adapter/http/dto"
	"chem-factory/internal/modules/user/core/port"
	"context"
	"log"
)

type userUsecase struct {
	userRepo port.UserRepository
}

func NewUserUsecase(userRepo port.UserRepository) *userUsecase {
	log.Println("Initializing user usecase")
	return &userUsecase{userRepo: userRepo}
}

func (service *userUsecase) GetProfile(ctx context.Context, userID uint) (dto.ProfileResponse, error) {

	user, err := service.userRepo.FindByID(ctx, userID)
	if err != nil {
		return dto.ProfileResponse{}, err
	}

	return dto.ProfileResponse{
		ID:       user.ID,
		Username: user.Username,
		Balance:  user.Balance,
		XP:       user.XP,
		Level:    user.Level,
	}, nil
}
