package service

import (
	u "chem-factory/internal/dto/user"
	"errors"
	"fmt"
)

func (s *serviceManager) UserData(token string) (u.UserDataResponse, error) {

	username, err := s.auth.VerifyJWT(token)
	if err != nil {
		return u.UserDataResponse{}, fmt.Errorf("Service user, get user data: %w", err)
	}

	user, err := s.repository.FindUserByUsername(username)
	if err != nil {
		return u.UserDataResponse{}, fmt.Errorf("Service user, get user data: %w", err)
	}

	return u.UserDataResponse{
		Username: user.Username,
		Balance:  user.Balance,
		XP:       user.XP,
		Level:    user.Level,
	}, nil
}

func (s *serviceManager) Login(request u.UserLoginRequest) (string, error) {

	savedHashedPassword, err := s.repository.FindPasswordByUsername(request.Username)
	if err != nil {
		return "", fmt.Errorf("Service user, login: %w", err)
	}

	if !s.auth.CheckPassword(request.Password, savedHashedPassword) {
		return "", fmt.Errorf("Service user, login: %w", err)
	}

	token, err := s.auth.GenerateJWT(request.Username)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *serviceManager) Register(userRR u.UserRegisterRequest) error {

	if _, err := s.repository.FindIDbyUsername(userRR.Username); err == nil {
		return errors.New("Username already exists.")
	}

	hashedPassword, err := s.auth.HashPassword(userRR.Password)
	if err != nil {
		return err
	}

	err = s.repository.SaveUser(userRR.Username, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}
