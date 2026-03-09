package service

import (
	"chem-factory/internal/model"
	"errors"
)

func (service *serviceManager) UserData(token string) (model.User, error) {

	var user model.User
	var err error

	if user.Username, err = service.auth.VerifyJWT(token); err != nil {
		return model.User{}, err
	}

	if err = service.repository.ExportUserByUsername(&user); err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (service *serviceManager) Login(username string, password string) (string, error) {

	savedHashedPassword, err := service.repository.ExportPasswordByUsername(username)
	if err != nil {
		return "", err
	}

	if !service.auth.CheckPassword(password, savedHashedPassword) {
		return "", errors.New("Password is wrong.")
	}

	token, err := service.auth.GenerateJWT(username)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (service *serviceManager) Register(username string, password string) error {

	if _, err := service.repository.ExportIDbyUsername(username); err == nil {
		return errors.New("Username already exists.")
	}

	hashedPassword, err := service.auth.HashPassword(password)
	if err != nil {
		return err
	}

	err = service.repository.SaveNewUser(model.User{Username: username, Password: hashedPassword})
	if err != nil {
		return err
	}

	return nil
}
