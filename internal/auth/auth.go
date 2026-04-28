package auth

import "os"

type AuthManager interface {
	// hash
	CheckPassword(password, hashedPassword string) error
	HashPassword(password string) (string, error)
	// jwt
	GenerateJWT(username string) (string, error)
	VerifyJWT(token string) (string, error)
}

type authManager struct {
	secretKey  string
}

func Init() *authManager {
	return &authManager{secretKey: os.Getenv("SECRET_KEY")}
}
