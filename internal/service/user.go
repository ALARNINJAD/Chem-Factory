package service

import (
	u "chem-factory/internal/dto/user"
	"errors"
)

func (service *serviceManager) UserData(token string) (u.UserDataResponse, error) {

	username, err := service.auth.VerifyJWT(token)
	if err != nil {
		return u.UserDataResponse{}, err
	}

	user, err := service.repository.FindUserByUsername(username)
	if err != nil {
		return u.UserDataResponse{}, err
	}

	userData := u.UserDataResponse{
		Username: user.Username,
		Balance:  user.Balance,
		XP:       user.XP,
		Level:    user.Level,
	}

	return userData, nil
}

func (service *serviceManager) Login(userLR u.UserLoginRequest) (string, error) {

	savedHashedPassword, err := service.repository.FindPasswordByUsername(userLR.Username)
	if err != nil {
		return "", err
	}

	if !service.auth.CheckPassword(userLR.Password, savedHashedPassword) {
		return "", errors.New("Password is wrong.")
	}

	token, err := service.auth.GenerateJWT(userLR.Username)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (service *serviceManager) Register(userRR u.UserRegisterRequest) error {

	// it's better to not use a functino like this
	if _, err := service.repository.FindIDbyUsername(userRR.Username); err == nil {
		return errors.New("Username already exists.")
	}

	hashedPassword, err := service.auth.HashPassword(userRR.Password)
	if err != nil {
		return err
	}

	err = service.repository.SaveUser(userRR.Username, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}
