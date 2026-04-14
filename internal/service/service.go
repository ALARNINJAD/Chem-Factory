package service

import (
	"chem-factory/internal/auth"
	u "chem-factory/internal/dto/user"
	"chem-factory/internal/repository"
)

type ServiceManager interface {
	// user
	Login(userLR u.UserLoginRequest) (string, error)
	Register(userRR u.UserRegisterRequest) error
	UserData(token string) (u.UserDataResponse, error)
}

type serviceManager struct {
	auth       auth.AuthManager
	repository repository.RepositoryManager
}

func Init(a auth.AuthManager, r repository.RepositoryManager) *serviceManager {
	return &serviceManager{auth: a, repository: r}
}
