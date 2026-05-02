package service

import (
	u "chem-factory/internal/dto/user"
	"errors"
	"fmt"
)

func (s *serviceManager) UserData(token string) (u.UserDataResponse, error) {

	username, err := s.auth.VerifyJWT(token)
	if err != nil {
		return u.UserDataResponse{}, fmt.Errorf("Service user, get user data: %w ", err)
	}

	user, err := s.repository.FindUserByUsername(username)
	if err != nil {
		return u.UserDataResponse{}, fmt.Errorf("Service user, get user data: %w ", err)
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
		return "", fmt.Errorf("Service user, loging: %w ", err)
	}

	if err := s.auth.CheckPassword(request.Password, savedHashedPassword); err != nil {
		return "", fmt.Errorf("Service user, loging: %w ", err)
	}

	token, err := s.auth.GenerateJWT(request.Username)
	if err != nil {
		return "", fmt.Errorf("Service user, loging: %w ", err)
	}

	return token, nil
}

func (s *serviceManager) Register(request u.UserRegisterRequest) error {

	if _, err := s.repository.FindIDbyUsername(request.Username); err == nil {
		return errors.New("Service user, register, find username: user already exists")
	}

	hashedPassword, err := s.auth.HashPassword(request.Password)
	if err != nil {
		return fmt.Errorf("Service user, register: %w ", err)
	}

	tx, err := s.repository.Transaction()
	if err != nil {
		return fmt.Errorf("Service user, register: %w ", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = s.repository.SaveUser(tx, request.Username, hashedPassword); err != nil {
		return fmt.Errorf("Service user, register: %w ", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("Service user, register, commit: %w ", err)
	}

	// err = s.notification.SendSMSWithProvider(s.notification.Kavenegar, notification.SimpleSMS{
	// 	Receptor: []string{"09120538596"},
	// 	Massage:  fmt.Sprintf("Welcome %s.", request.Username),
	// })
	// if err != nil {
	// 	log.Printf("Service user, register: %s ", err)
	// }

	return nil
}
