package dto

import "chem-factory/pkg/dto"

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type RegisterResponse struct {
	dto.MessageResponse
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type LoginResponse struct {
	UserID uint   `json:"user_id"`
	Token  string `json:"token"`
	dto.MessageResponse
}
