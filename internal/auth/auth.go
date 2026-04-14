package auth

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type AuthManager interface {
	// hash
	CheckPassword(password, hashedPassword string) bool
	HashPassword(password string) (string, error)
	// jwt
	GenerateJWT(username string) (string, error)
	VerifyJWT(token string) (string, error)
}

type authManager struct {
	secretKey  string `json:"secret_key"`
	costOfHash int    `json:"cost_of_hash"`
}

func Init() *authManager {
	data, err := os.ReadFile(filepath.Join(".", "configs", "auth.json"))
	if err != nil {
		panic("Could not access auth config file")
	}

	var auth []authManager
	err = json.Unmarshal(data, &auth)
	if err != nil {
		panic("Could not access auth config")
	}

	return &auth[0]
}
