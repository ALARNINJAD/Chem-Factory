package service

import (
	"chem-factory/internal/dto/user"
	"errors"
	"fmt"
)

func (service *Manager) UserData(token string) (user.DataResponse, error) {

	a := service.auth

	username, err := a.JWT.Verify(token)
	if err != nil {
		return user.DataResponse{}, fmt.Errorf("Service user, get user data: %w", err)
	}

	r := service.repository

	usr, err := r.User.FindByUsername(username)
	if err != nil {
		return user.DataResponse{}, fmt.Errorf("Service user, get user data: %w", err)
	}

	return user.DataResponse{
		Username: usr.Username,
		Balance:  usr.Balance,
		XP:       usr.XP,
		Level:    usr.Level,
	}, nil
}

func (service *Manager) Login(request user.LoginRequest) (string, error) {

	r := service.repository

	savedHashedPassword, err := r.User.FindPasswordByUsername(request.Username)
	if err != nil {
		return "", fmt.Errorf("Service user, loging: %w", err)
	}

	a := service.auth

	if err := a.Hash.CheckPassword(request.Password, savedHashedPassword); err != nil {
		return "", fmt.Errorf("Service user, loging: %w", err)
	}

	token, err := a.JWT.Generate(request.Username)
	if err != nil {
		return "", fmt.Errorf("Service user, loging: %w", err)
	}

	return token, nil
}

func (service *Manager) Register(request user.RegisterRequest) error {

	r := service.repository

	if _, err := r.User.FindIDbyUsername(request.Username); err == nil {
		return errors.New("Service user, register, find username: user already exists.")
	}

	a := service.auth

	hashedPassword, err := a.Hash.HashPassword(request.Password)
	if err != nil {
		return fmt.Errorf("Service user, register: %w", err)
	}

	tx, err := r.Transaction()
	if err != nil {
		return fmt.Errorf("Service user, register: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = r.User.Add(tx, request.Username, hashedPassword); err != nil {
		return fmt.Errorf("Service user, register: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("Service user, register, commit: %w", err)
	}

	// n := service.notification
	// err = n.SendSMSWithProvider(n.Kavenegar, notification.SimpleSMS{
	// 	Receptor: []string{"09120538596"},
	// 	Massage:  fmt.Sprintf("Welcome %s", request.Username),
	// })
	// if err != nil {
	// 	log.Printf("Service user, register: %s", err)
	// }

	return nil
}
