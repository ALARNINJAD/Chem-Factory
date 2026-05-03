package auth

import "os"

type Manager struct {
	Hash      *hashManager
	JWT       *jwtManager
}

func New() *Manager {
	return &Manager{
		Hash: NewHashManager(),
		JWT: NewJWTManager(os.Getenv("SECRET_KEY")),
	}
}
