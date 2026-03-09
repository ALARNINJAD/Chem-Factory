package service

import (
	"chem-factory/internal/auth"
	"chem-factory/internal/model"
	"chem-factory/internal/repository"
)

type ServiceManager interface {
	// user
	Login(username string, password string) (string, error)
	Register(username string, password string) error
	UserData(token string) (model.User, error)
}

type serviceManager struct {
	auth       auth.AuthManager
	repository repository.RepositoryManager
}

func Init(a auth.AuthManager, r repository.RepositoryManager) *serviceManager {
	return &serviceManager{auth: a, repository: r}
}
