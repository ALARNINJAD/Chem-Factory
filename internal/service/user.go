package service

import (
	u "chem-factory/internal/dto/user"
	"chem-factory/internal/notification"
	"errors"
	"fmt"
	"log"
)

func (m *Manager) UserData(token string) (u.UserDataResponse, error) {

	a := m.auth

	username, err := a.JWT.Verify(token)
	if err != nil {
		return u.UserDataResponse{}, fmt.Errorf("Service user, get user data: %w", err)
	}

	r := m.repository

	user, err := r.User.FindByUsername(username)
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

func (m *Manager) Login(request u.UserLoginRequest) (string, error) {

	r := m.repository

	savedHashedPassword, err := r.User.FindPasswordByUsername(request.Username)
	if err != nil {
		return "", fmt.Errorf("Service user, loging: %w", err)
	}

	a := m.auth

	if err := a.Hash.CheckPassword(request.Password, savedHashedPassword); err != nil {
		return "", fmt.Errorf("Service user, loging: %w", err)
	}

	token, err := a.JWT.Generate(request.Username)
	if err != nil {
		return "", fmt.Errorf("Service user, loging: %w", err)
	}

	return token, nil
}

func (m *Manager) Register(request u.UserRegisterRequest) error {

	r := m.repository

	if _, err := r.User.FindIDbyUsername(request.Username); err == nil {
		return errors.New("Service user, register, find username: user already exists.")
	}

	a := m.auth

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

	n := m.notification

	err = n.SendSMSWithProvider(n.Kavenegar, notification.SimpleSMS{
		Receptor: []string{"09120538596"},
		Massage:  fmt.Sprintf("Welcome %s", request.Username),
	})
	if err != nil {
		log.Printf("Service user, register: %s", err)
	}

	return nil
}
